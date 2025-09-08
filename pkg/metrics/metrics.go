package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"time"
)

// Метрики Prometheus для Telegram-бота

var (
	// MessagesTotal — общее количество сообщений, полученных ботом
	// Лейбл "chat_id" позволяет различать сообщения по чатам
	MessagesTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "bot_messages_total",
			Help: "Total messages received by the bot",
		},
		[]string{"chat_id"},
	)

	// RepliesTotal — общее количество отправленных ботом ответов
	RepliesTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "bot_replies_total",
			Help: "Total replies sent by the bot",
		},
	)

	// NoMatchTotal — количество сообщений, на которые не найдено совпадений с правилами
	NoMatchTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "bot_messages_no_match_total",
			Help: "Messages received with no matching rules",
		},
	)

	// ErrorsTotal — количество ошибок на разных стадиях обработки сообщений
	// Лейбл "stage" указывает этап, на котором произошла ошибка
	ErrorsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "bot_errors_total",
			Help: "Total errors occurred in bot processing",
		},
		[]string{"stage"},
	)

	// RuleHitsTotal — количество срабатываний каждого правила
	// Лейбл "rule" хранит текст правила
	RuleHitsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "bot_rule_hits_total",
			Help: "Number of times each rule was triggered",
		},
		[]string{"rule"},
	)

	// MessageProcessingDuration — гистограмма времени обработки одного сообщения
	MessageProcessingDuration = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "bot_message_processing_duration_seconds",
			Help:    "Duration to process one message",
			Buckets: prometheus.DefBuckets, // стандартные интервалы Prometheus
		},
	)

	// ConfigReloadDuration — время выполнения reload конфигурации
	ConfigReloadDuration = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "bot_config_reload_duration_seconds",
			Help: "Time spent to reload configuration",
		},
	)

	// ConfigReloadTotal — количество успешных reload конфигурации
	ConfigReloadTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "bot_config_reload_total",
			Help: "Number of times config was reloaded",
		},
	)

	// ConfigReloadErrorsTotal — количество ошибок при reload конфигурации
	ConfigReloadErrorsTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "bot_config_reload_errors_total",
			Help: "Number of errors during config reload",
		},
	)
)

// InitMetrics регистрирует все метрики в Prometheus
func InitMetrics() {
	prometheus.MustRegister(
		MessagesTotal,
		RepliesTotal,
		NoMatchTotal,
		ErrorsTotal,
		RuleHitsTotal,
		MessageProcessingDuration,
		ConfigReloadDuration,
		ConfigReloadTotal,
		ConfigReloadErrorsTotal,
	)
}

// ObserveProcessing измеряет длительность обработки сообщения
// и обновляет гистограмму MessageProcessingDuration
func ObserveProcessing(start time.Time) {
	duration := time.Since(start).Seconds()
	MessageProcessingDuration.Observe(duration)
}
