package config

import (
	"fmt"
	"regexp"
)

// Compile компилирует строковое регулярное выражение Rule.Pattern
// и сохраняет его в поле re для последующего использования.
// Возвращает ошибку, если регулярное выражение некорректное.
func (r *Rule) Compile() error {
	re, err := regexp.Compile(r.Pattern)
	if err != nil {
		return fmt.Errorf("invalid regexp %q: %w", r.Pattern, err)
	}
	r.re = re
	return nil
}

// Re возвращает скомпилированное регулярное выражение.
// Используется для поиска или проверки текста согласно правилу.
func (r *Rule) Re() *regexp.Regexp {
	return r.re
}
