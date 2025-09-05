package app

import (
	"errors"
	"fmt"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/st-kuptsov/balabol/config"
	"github.com/st-kuptsov/balabol/internal/telegram"
	logs "github.com/st-kuptsov/balabol/pkg/logs"
	"github.com/st-kuptsov/balabol/pkg/metrics"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Run запускает приложение
func Run(version string) error {
	configPath := "config/config.yaml"

	conf, err := config.LoadConfigWithHash(configPath)
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	logger := logs.DefaultLogger(conf.Config.Logging)
	logger.Infow("starting balabol",
		"config", configPath,
		"logLevel", conf.Config.Logging.Level,
		"pid", os.Getpid(),
		"version", version,
	)

	// Инициализация метрик
	logger.Debug("initializing metrics server")
	metrics.InitMetrics()
	startMetricsServer(conf.Config.ServicePort, logger)

	// Инициализация Telegram-бота
	logger.Debug("initializing telegram bot")
	bot, err := telegram.NewBot(
		conf.Config.Telegram.Token,
		conf.Config.BotMode,
		func() []config.Rule { return conf.Config.Rules },
	)
	if err != nil {
		return fmt.Errorf("telegram bot init: %w", err)
	}
	logger.Info("telegram bot initialized")

	// Запуск бота
	go bot.Start()

	// канал для остановки
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// reload конфигурации
	go reloadConfigLoop(configPath, conf, stop, logger)

	// ожидание сигнала остановки
	<-stop
	logger.Info("shutting down gracefully...")
	bot.Stop()
	return nil
}

// startMetricsServer запускает HTTP сервер для Prometheus
func startMetricsServer(port int, logger *zap.SugaredLogger) {
	go func() {
		servicePort := fmt.Sprintf(":%d", port)
		http.Handle("/metrics", promhttp.Handler())
		logger.Infow("metrics server started", "port", servicePort)
		if err := http.ListenAndServe(servicePort, nil); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Errorw("metrics server failed", "error", err)
		}
	}()
}

// reloadConfigLoop проверяет изменения конфига каждые 5 секунд
func reloadConfigLoop(path string, conf *config.CachedConfig, stop chan os.Signal, logger *zap.SugaredLogger) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-stop:
			return
		case <-ticker.C:
			changed, err := conf.ReloadWithMetrics(path)
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
