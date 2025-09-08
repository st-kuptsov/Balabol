package telegram

import (
	"go.uber.org/zap" // структурированное логирование
	"strconv"
	"strings"
	"time"

	"github.com/st-kuptsov/balabol/config"      // конфигурация приложения и правила
	"github.com/st-kuptsov/balabol/pkg/metrics" // метрики Prometheus
	tb "gopkg.in/telebot.v3"                    // библиотека для Telegram-бота
)

// NewBot создаёт и настраивает Telegram-бота.
//
// Параметры:
// - token: токен Telegram-бота
// - mode: режим работы бота (например, first_last)
// - cleanFilter: фильтр символов для очистки текста
// - removeDup: удалять ли повторяющиеся буквы
// - rulesFn: функция, возвращающая текущий список правил
// - logger: экземпляр структурированного логгера
//
// Возвращает:
// - указатель на tb.Bot
// - ошибку, если инициализация не удалась
func NewBot(token, mode, cleanFilter string, removeDup bool, rulesFn func() []config.Rule, logger *zap.SugaredLogger) (*tb.Bot, error) {
	// Настройки Telegram-бота: токен и длинный polling
	pref := tb.Settings{
		Token:  token,
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	}

	// Создаём бота
	bot, err := tb.NewBot(pref)
	if err != nil {
		return nil, err
	}

	// Обработчик входящих текстовых сообщений
	bot.Handle(tb.OnText, func(c tb.Context) error {
		start := time.Now() // для метрик времени обработки

		// Очистка текста: убираем лишние символы и дубликаты
		text := cleanText(strings.TrimSpace(c.Message().Text), cleanFilter, removeDup, logger)
		chatID := strconv.FormatInt(c.Chat().ID, 10)

		// Увеличиваем общий счетчик сообщений
		metrics.MessagesTotal.WithLabelValues(chatID).Inc()

		// Если текст пустой после очистки — учитываем как "no match"
		if text == "" {
			metrics.NoMatchTotal.Inc()
			metrics.ObserveProcessing(start)
			return nil
		}

		// Получаем текущие правила через функцию rulesFn
		rules := rulesFn()
		hits := MatchRules(text, rules, mode, logger)

		// Если нет совпадений — учитываем как "no match"
		if len(hits) == 0 {
			metrics.NoMatchTotal.Inc()
			metrics.ObserveProcessing(start)
			return nil
		}

		// Формируем ответ бота и обновляем метрики
		replies := make([]string, 0, len(hits))
		for _, h := range hits {
			replies = append(replies, h.resp)
			metrics.RuleHitsTotal.WithLabelValues(h.ruleText).Inc()
		}

		reply := strings.Join(replies, ". ") // объединяем все ответы в один текст
		metrics.RepliesTotal.Inc()
		metrics.ObserveProcessing(start) // фиксируем длительность обработки

		// Отправляем ответ пользователю
		return c.Reply(reply, &tb.SendOptions{ReplyTo: c.Message()})
	})

	return bot, nil
}
