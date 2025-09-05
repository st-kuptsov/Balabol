package telegram

import (
	"strconv"
	"strings"
	"time"

	"github.com/st-kuptsov/balabol/config"
	"github.com/st-kuptsov/balabol/pkg/metrics"
	tb "gopkg.in/telebot.v3"
)

// NewBot создает и возвращает Telegram-бота
func NewBot(token string, mode string, rulesFn func() []config.Rule) (*tb.Bot, error) {
	pref := tb.Settings{
		Token:  token,
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	}

	bot, err := tb.NewBot(pref)
	if err != nil {
		return nil, err
	}

	bot.Handle(tb.OnText, func(c tb.Context) error {
		start := time.Now()

		text := cleanText(strings.TrimSpace(c.Message().Text))
		chatID := strconv.FormatInt(c.Chat().ID, 10)
		metrics.MessagesTotal.WithLabelValues(chatID).Inc()

		if text == "" {
			metrics.NoMatchTotal.Inc()
			metrics.ObserveProcessing(start)
			return nil
		}

		rules := rulesFn()
		hits := MatchRules(text, rules, mode)

		if len(hits) == 0 {
			metrics.NoMatchTotal.Inc()
			metrics.ObserveProcessing(start)
			return nil
		}

		// Формируем ответ и обновляем метрики
		replies := make([]string, 0, len(hits))
		for _, h := range hits {
			replies = append(replies, h.resp)
			metrics.RuleHitsTotal.WithLabelValues(h.ruleText).Inc()
		}

		reply := strings.Join(replies, ". ")
		metrics.RepliesTotal.Inc()
		metrics.ObserveProcessing(start)

		return c.Reply(reply, &tb.SendOptions{ReplyTo: c.Message()})
	})

	return bot, nil
}
