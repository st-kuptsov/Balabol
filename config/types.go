package config

import "regexp"

type Config struct {
	Rules       []Rule         `yaml:"rules"`
	Telegram    TelegramConfig `yaml:"telegram"`
	Logging     LogConfig      `yaml:"log_settings"`
	BotMode     string         `yaml:"bot_mode" env-default:"first_last"`
	SecretsPath string         `yaml:"secrets"`
	ServicePort int            `yaml:"service_port" env-default:"9090"`
}

type Rule struct {
	Text     string         `yaml:"text"`
	Pattern  string         `yaml:"pattern"`
	Response string         `yaml:"response"`
	re       *regexp.Regexp `yaml:"-"`
}

type TelegramConfig struct {
	Token string
}

type LogConfig struct {
	Directory  string `yaml:"directory" env-default:"logs"`
	Filename   string `yaml:"filename" env-default:"app.log"`
	MaxSize    int    `yaml:"max_size" env-default:"10"`
	MaxBackups int    `yaml:"max_backups" env-default:"1"`
	MaxAge     int    `yaml:"max_age" env-default:"1"`
	Compress   bool   `yaml:"compress" env-default:"true"`
	Level      string `yaml:"level" env-default:"info"`
	Console    bool   `yaml:"console_enabled" env-default:"false"`
}

type CachedConfig struct {
	Config      *Config
	ConfigHash  string
	SecretsHash string
}
