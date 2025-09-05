package config

import (
	"fmt"
	"regexp"
)

// Compile компилирует регулярное выражение
func (r *Rule) Compile() error {
	re, err := regexp.Compile(r.Pattern)
	if err != nil {
		return fmt.Errorf("invalid regexp %q: %w", r.Pattern, err)
	}
	r.re = re
	return nil
}

// Re возвращает скомпилированное регулярное выражение
func (r *Rule) Re() *regexp.Regexp {
	return r.re
}
