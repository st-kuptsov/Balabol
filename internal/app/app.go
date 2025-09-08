package app

import (
	"errors"
	"fmt"
	"github.com/prometheus/client_golang/prometheus/promhttp" // обработчик метрик Prometheus
	"github.com/st-kuptsov/balabol/config"                    // работа с конфигурацией
	"github.com/st-kuptsov/balabol/internal/telegram"         // Telegram-бот
	logs "github.com/st-kuptsov/balabol/pkg/logs"             // кастомный логгер
	"github.com/st-kuptsov/balabol/pkg/metrics"               // инициализация метрик
	"go.uber.org/zap"                                         // структурированное логирование
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Run запускает основную логику приложения.
// version — версия приложения, передается для логирования.
func Run(version string) error {
	// Путь к конфигурационному файлу
	configPath := "config/config.yaml"

	// Загружаем конфиг с контролем хеша (чтобы отслеживать изменения)
	conf, err := config.LoadConfigWithHash(configPath)
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	// Создаём логгер согласно конфигурации
	logger := logs.DefaultLogger(conf.Config.Logging)
	logger.Infow("starting balabol",
		"config", configPath,
		"logLevel", conf.Config.Logging.Level,
		"pid", os.Getpid(),
		"version", version,
	)

	// Инициализация метрик Prometheus
	logger.Debug("initializing metrics server")
	metrics.InitMetrics()                               // инициализация метрик приложения
	startMetricsServer(conf.Config.ServicePort, logger) // запуск HTTP-сервера для метрик

	// Инициализация Telegram-бота
	logger.Debug("initializing telegram bot")
	bot, err := telegram.NewBot(
		conf.Config.Telegram.Token,
		conf.Config.BotMode,
		conf.Config.CleanFilter,
		conf.Config.RemoveDup,
		func() []config.Rule { return conf.Config.Rules }, // ленивый доступ к правилам
		logger,
	)
	if err != nil {
		return fmt.Errorf("telegram bot init: %w", err)
	}
	logger.Info("telegram bot initialized")

	// Запуск бота в отдельной горутине
	go bot.Start()

	// Канал для сигналов остановки приложения
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// Запуск горутины, которая периодически перезагружает конфиг при изменении
	go reloadConfigLoop(configPath, conf, stop, logger)

	// Основная блокировка: ждем сигнала остановки
	<-stop
	logger.Info("shutting down gracefully...")
	bot.Stop() // остановка бота
	return nil
}

// startMetricsServer запускает HTTP-сервер для Prometheus метрик
func startMetricsServer(port int, logger *zap.SugaredLogger) {
	go func() {
		servicePort := fmt.Sprintf(":%d", port)
		http.Handle("/metrics", promhttp.Handler()) // обработчик метрик
		logger.Infow("metrics server started", "port", servicePort)
		// Запуск HTTP сервера
		if err := http.ListenAndServe(servicePort, nil); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Errorw("metrics server failed", "error", err)
		}
	}()
}

// reloadConfigLoop периодически (каждые 5 секунд) проверяет изменения конфигурации.
// Если конфиг изменился, логгирует обновление и применяет новые правила.
func reloadConfigLoop(path string, conf *config.CachedConfig, stop chan os.Signal, logger *zap.SugaredLogger) {
	ticker := time.NewTicker(5 * time.Second) // тикер с интервалом 5 секунд
	defer ticker.Stop()

	for {
		select {
		case <-stop: // сигнал на завершение
			return
		case <-ticker.C: // тикер срабатывает
			changed, err := conf.ReloadWithMetrics(path) // проверка и перезагрузка конфига
			if err != nil {
				logger.Errorw("config reload failed", "error", err)
				continue
			}
			if changed {
				logger.Infow("config reloaded", "mode", conf.Config.BotMode, "rules_count", len(conf.Config.Rules))
			}
		}
	}
}
