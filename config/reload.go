package config

import (
	"crypto/sha256"
	"fmt"
	"os"
	"time"

	"github.com/st-kuptsov/balabol/pkg/metrics"
)

// ReloadIfChanged перечитывает конфиг и секреты, если хеш изменился
func (c *CachedConfig) ReloadIfChanged(path string) (bool, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, fmt.Errorf("cannot read config file: %w", err)
	}

	newHash := fmt.Sprintf("%x", sha256.Sum256(data))
	cfg, err := GetConfig(path)
	if err != nil {
		return false, err
	}

	changed := false

	if newHash != c.ConfigHash {
		c.Config = cfg
		c.ConfigHash = newHash
		changed = true
	}

	// Проверяем secrets
	if cfg.SecretsPath != "" {
		sData, err := os.ReadFile(cfg.SecretsPath)
		if err != nil {
			return changed, fmt.Errorf("cannot read secrets file: %w", err)
		}
		newSecretsHash := fmt.Sprintf("%x", sha256.Sum256(sData))
		if newSecretsHash != c.SecretsHash {
			if err := c.Config.LoadSecrets(); err != nil {
				return changed, err
			}
			c.SecretsHash = newSecretsHash
			changed = true
		}
	}

	// Перекомпилируем regexp
	for i := range c.Config.Rules {
		if err := c.Config.Rules[i].Compile(); err != nil {
			return changed, fmt.Errorf("invalid regexp %q: %w", c.Config.Rules[i].Pattern, err)
		}
	}

	return changed, nil
}

// ReloadWithMetrics проверяет изменения конфигурации и обновляет метрики
func (c *CachedConfig) ReloadWithMetrics(path string) (bool, error) {
	start := time.Now()
	changed, err := c.ReloadIfChanged(path)
	duration := time.Since(start).Seconds()
	metrics.ConfigReloadDuration.Set(duration)

	if err != nil {
		metrics.ConfigReloadErrorsTotal.Inc()
		return changed, err
	}

	if changed {
		metrics.ConfigReloadTotal.Inc()
	}

	return changed, nil
}
