package config

import (
	"crypto/sha256" // для вычисления хеша конфигурационных файлов
	"fmt"
	"github.com/ilyakaznacheev/cleanenv" // библиотека для чтения YAML/ENV конфигов
	"log"
	"os"
)

// GetConfig загружает конфигурацию из файла по указанному пути.
// Выполняет:
// 1. Чтение основного конфига через cleanenv
// 2. Загрузку секретов (например, токена Telegram)
// 3. Компиляцию правил (Rule.Compile)
func GetConfig(path string) (*Config, error) {
	var cfg Config

	// Чтение конфигурации из файла
	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		// Если произошла ошибка, выводим справку по конфигу
		if help, err2 := cleanenv.GetDescription(&cfg, nil); err2 == nil {
			log.Println(help)
		}
		return nil, fmt.Errorf("ошибка загрузки конфигурации: %w", err)
	}

	// Загрузка секретов из отдельного файла (если указан)
	if err := cfg.LoadSecrets(); err != nil {
		return nil, err
	}

	// Компиляция всех правил из конфига (например, регулярные выражения)
	for i := range cfg.Rules {
		if err := cfg.Rules[i].Compile(); err != nil {
			return nil, err
		}
	}

	return &cfg, nil
}

// LoadSecrets загружает секреты из отдельного файла (SecretsPath).
// Поддерживаются, например, токены для Telegram.
func (c *Config) LoadSecrets() error {
	if c.SecretsPath == "" {
		// Если путь к секретам не указан, пропускаем
		return nil
	}

	// Структура для парсинга секретов
	type secrets struct {
		Telegram struct {
			Token string `yaml:"token"`
		} `yaml:"telegram"`
	}

	var sec secrets
	// Чтение файла секретов
	if err := cleanenv.ReadConfig(c.SecretsPath, &sec); err != nil {
		return err
	}

	// Присвоение токена из секрета в основную конфигурацию
	c.Telegram.Token = sec.Telegram.Token
	return nil
}

// LoadConfigWithHash загружает конфигурацию и вычисляет SHA256 хеши
// 1. Основного конфига
// 2. Файла секретов (если указан)
// Возвращает CachedConfig, который хранит конфиг и его хеши
func LoadConfigWithHash(path string) (*CachedConfig, error) {
	// Чтение основного конфиг файла
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("cannot read config file: %w", err)
	}
	// Вычисление SHA256 хеша основного конфига
	hash := fmt.Sprintf("%x", sha256.Sum256(data))

	// Загрузка конфигурации
	cfg, err := GetConfig(path)
	if err != nil {
		return nil, err
	}

	// Вычисление хеша для файла секретов, если он есть
	secretsHash := ""
	if cfg.SecretsPath != "" {
		if sData, err := os.ReadFile(cfg.SecretsPath); err == nil {
			secretsHash = fmt.Sprintf("%x", sha256.Sum256(sData))
		}
	}

	return &CachedConfig{
		Config:      cfg,
		ConfigHash:  hash,
		SecretsHash: secretsHash,
	}, nil
}
