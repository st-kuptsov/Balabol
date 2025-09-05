package config

import (
	"crypto/sha256"
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
)

func GetConfig(path string) (*Config, error) {
	var cfg Config
	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		if help, err2 := cleanenv.GetDescription(&cfg, nil); err2 == nil {
			log.Println(help)
		}
		return nil, fmt.Errorf("ошибка загрузки конфигурации: %w", err)
	}

	if err := cfg.LoadSecrets(); err != nil {
		return nil, err
	}

	for i := range cfg.Rules {
		if err := cfg.Rules[i].Compile(); err != nil {
			return nil, err
		}
	}

	return &cfg, nil
}

func (c *Config) LoadSecrets() error {
	if c.SecretsPath == "" {
		return nil
	}

	type secrets struct {
		Telegram struct {
			Token string `yaml:"token"`
		} `yaml:"telegram"`
	}

	var sec secrets
	if err := cleanenv.ReadConfig(c.SecretsPath, &sec); err != nil {
		return err
	}

	c.Telegram.Token = sec.Telegram.Token
	return nil
}

func LoadConfigWithHash(path string) (*CachedConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("cannot read config file: %w", err)
	}
	hash := fmt.Sprintf("%x", sha256.Sum256(data))

	cfg, err := GetConfig(path)
	if err != nil {
		return nil, err
	}

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
