package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"time"
)

var (
	// Количество полученных сообщений
	MessagesTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "bot_messages_total",
			Help: "Total messages received by the bot",
		},
		[]string{"chat_id"},
	)

	// Количество отправленных ответов
	RepliesTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "bot_replies_total",
			Help: "Total replies sent by the bot",
		},
	)

	// Сообщения, на которые не найдено совпадений
	NoMatchTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "bot_messages_no_match_total",
			Help: "Messages received with no matching rules",
		},
	)

	// Ошибки обработки сообщений или отправки
	ErrorsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "bot_errors_total",
			Help: "Total errors occurred in bot processing",
		},
		[]string{"stage"},
	)

	// Количество срабатываний правил
	RuleHitsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "bot_rule_hits_total",
			Help: "Number of times each rule was triggered",
		},
		[]string{"rule"},
	)

	// Время обработки одного сообщения
	MessageProcessingDuration = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "bot_message_processing_duration_seconds",
			Help:    "Duration to process one message",
			Buckets: prometheus.DefBuckets,
		},
	)

	// Время reload конфигурации
	ConfigReloadDuration = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "bot_config_reload_duration_seconds",
			Help: "Time spent to reload configuration",
		},
	)

	// Количество reload конфигов
	ConfigReloadTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "bot_config_reload_total",
			Help: "Number of times config was reloaded",
		},
	)

	// Ошибки при reload
	ConfigReloadErrorsTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "bot_config_reload_errors_total",
			Help: "Number of errors during config reload",
		},
	)
)

// InitMetrics регистрирует все метрики
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

// Helper: измерение времени обработки
func ObserveProcessing(start time.Time) {
	duration := time.Since(start).Seconds()
	MessageProcessingDuration.Observe(duration)
}
