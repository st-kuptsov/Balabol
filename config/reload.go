package config

import (
	"crypto/sha256"
	"fmt"
	"os"
	"time"

	"github.com/st-kuptsov/balabol/pkg/metrics" // кастомные Prometheus метрики
)

// ReloadIfChanged проверяет, изменился ли конфиг или секреты.
// Если изменился — перечитывает их и обновляет CachedConfig.
// Возвращает:
//   - changed = true, если конфиг или секреты обновились
//   - error при проблемах с чтением или компиляцией правил
func (c *CachedConfig) ReloadIfChanged(path string) (bool, error) {
	// Чтение основного конфига
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			// Если файла нет, считаем что изменений нет
			return false, nil
		}
		return false, fmt.Errorf("cannot read config file: %w", err)
	}

	// Вычисляем SHA256 хеш файла
	newHash := fmt.Sprintf("%x", sha256.Sum256(data))

	// Загружаем конфиг
	cfg, err := GetConfig(path)
	if err != nil {
		return false, err
	}

	changed := false

	// Проверяем, изменился ли основной конфиг
	if newHash != c.ConfigHash {
		c.Config = cfg
		c.ConfigHash = newHash
		changed = true
	}

	// Проверяем файл секретов
	if cfg.SecretsPath != "" {
		sData, err := os.ReadFile(cfg.SecretsPath)
		if err != nil {
			return changed, fmt.Errorf("cannot read secrets file: %w", err)
		}
		newSecretsHash := fmt.Sprintf("%x", sha256.Sum256(sData))
		if newSecretsHash != c.SecretsHash {
			// Перечитываем секреты
			if err := c.Config.LoadSecrets(); err != nil {
				return changed, err
			}
			c.SecretsHash = newSecretsHash
			changed = true
		}
	}

	// Перекомпилируем все правила (например, регулярные выражения)
	for i := range c.Config.Rules {
		if err := c.Config.Rules[i].Compile(); err != nil {
			return changed, fmt.Errorf("invalid regexp %q: %w", c.Config.Rules[i].Pattern, err)
		}
	}

	return changed, nil
}

// ReloadWithMetrics проверяет изменения конфигурации и обновляет Prometheus метрики
func (c *CachedConfig) ReloadWithMetrics(path string) (bool, error) {
	start := time.Now()

	// Проверка изменений
	changed, err := c.ReloadIfChanged(path)

	// Записываем длительность операции в метрики
	duration := time.Since(start).Seconds()
	metrics.ConfigReloadDuration.Set(duration)

	// В случае ошибки увеличиваем счетчик ошибок
	if err != nil {
		metrics.ConfigReloadErrorsTotal.Inc()
		return changed, err
	}

	// Если были изменения — увеличиваем счетчик успешных reload
	if changed {
		metrics.ConfigReloadTotal.Inc()
	}

	return changed, nil
}
