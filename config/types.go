package config

import "regexp"

// Config представляет основную конфигурацию приложения.
type Config struct {
	Rules       []Rule         `yaml:"rules"`                             // Список правил фильтрации/ответов
	Telegram    TelegramConfig `yaml:"telegram"`                          // Настройки Telegram-бота
	Logging     LogConfig      `yaml:"log_settings"`                      // Настройки логирования
	CleanFilter string         `yaml:"clean_filter"`                      // Фильтр для очистки текста перед обработкой
	RemoveDup   bool           `yaml:"remove_duplicate_letters"`          // Удалять ли повторяющиеся буквы
	BotMode     string         `yaml:"bot_mode" env-default:"first_last"` // Режим работы бота
	SecretsPath string         `yaml:"secrets"`                           // Путь к файлу секретов (например, токен Telegram)
	ServicePort int            `yaml:"service_port" env-default:"9090"`   // Порт сервиса для Prometheus метрик
}

// Rule представляет одно правило для бота:
//   - Pattern — регулярное выражение для сопоставления текста
//   - Response — текст ответа, если правило сработало
//   - Text — дополнительное описание правила
//   - re — скомпилированное регулярное выражение (не сохраняется в YAML)
type Rule struct {
	Text     string         `yaml:"text"`     // Описание правила
	Pattern  string         `yaml:"pattern"`  // Регулярное выражение в виде строки
	Response string         `yaml:"response"` // Ответ бота при совпадении
	re       *regexp.Regexp `yaml:"-"`        // Скомпилированное регулярное выражение
}

// TelegramConfig хранит настройки Telegram-бота
type TelegramConfig struct {
	Token string // Токен бота
}

// LogConfig хранит настройки логирования приложения
type LogConfig struct {
	Directory  string `yaml:"directory" env-default:"logs"`        // Директория для логов
	Filename   string `yaml:"filename" env-default:"app.log"`      // Имя файла логов
	MaxSize    int    `yaml:"max_size" env-default:"10"`           // Максимальный размер файла в МБ
	MaxBackups int    `yaml:"max_backups" env-default:"1"`         // Количество резервных копий
	MaxAge     int    `yaml:"max_age" env-default:"1"`             // Срок хранения логов в днях
	Compress   bool   `yaml:"compress" env-default:"true"`         // Сжимать ли старые файлы
	Level      string `yaml:"level" env-default:"info"`            // Уровень логирования (debug/info/warn/error)
	Console    bool   `yaml:"console_enabled" env-default:"false"` // Вывод логов в консоль
}

// CachedConfig хранит загруженный конфиг и хеши файлов
// для отслеживания изменений и безопасного reload
type CachedConfig struct {
	Config      *Config // Основная конфигурация
	ConfigHash  string  // SHA256 хеш основного конфига
	SecretsHash string  // SHA256 хеш файла секретов
}
