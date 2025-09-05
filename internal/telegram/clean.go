package telegram

import (
	"regexp"
	"strings"
)

// cleanText очищает текст: убирает все символы кроме букв/цифр/пробелов
// и удаляет повторяющиеся символы подряд
func cleanText(input string) string {
	re := regexp.MustCompile(`[^a-zA-Zа-яА-ЯёЁ0-9 ]+`)
	cleaned := re.ReplaceAllString(input, "")

	var result strings.Builder
	var prev rune
	for _, r := range cleaned {
		if r != prev || r == ' ' {
			result.WriteRune(r)
		}
		prev = r
	}
	return strings.TrimSpace(result.String())
}
