package telegram

import (
	"go.uber.org/zap" // структурированное логирование
	"regexp"
	"strings"
	"unicode"
)

// cleanText очищает входной текст перед обработкой ботом.
// Параметры:
// - input: исходный текст
// - cleanFilter: регулярное выражение для удаления нежелательных символов
// - removeDup: флаг удаления повторяющихся букв и дубликатов слов
// - logger: логгер для отладки
//
// Функция возвращает "очищенный" текст.
func cleanText(input, cleanFilter string, removeDup bool, logger *zap.SugaredLogger) string {
	logger.Debugw("cleanText Input", "input", input)
	for i, r := range input {
		logger.Debugw("Character", "index", i, "unicode", string(r), "code", r)
	}

	// Удаляем символы, не подходящие под cleanFilter
	re := regexp.MustCompile(cleanFilter)
	cleaned := re.ReplaceAllString(input, "")
	logger.Debugw("After clean_filter", "result", cleaned)

	// Приведение текста к нижнему регистру
	cleaned = strings.ToLower(cleaned)
	logger.Debugw("After ToLower", "result", cleaned)

	// Если нужно удалять дубликаты букв и слов
	if removeDup {
		words := strings.Fields(cleaned) // разбиваем на слова
		for i, word := range words {
			var result strings.Builder
			var prev rune
			for _, r := range word {
				// оставляем цифры и символы, которые не повторяются подряд
				if unicode.IsDigit(r) || r != prev {
					result.WriteRune(r)
				}
				prev = r
			}
			words[i] = result.String()
		}
		cleaned = strings.Join(words, " ")

		// Убираем повторяющиеся слова
		words = strings.Fields(cleaned)
		uniqueWords := make([]string, 0, len(words))
		seen := make(map[string]bool)
		for _, word := range words {
			if !seen[word] {
				uniqueWords = append(uniqueWords, word)
				seen[word] = true
			}
		}
		cleaned = strings.Join(uniqueWords, " ")
	}

	// Нормализация пробелов: заменяем несколько пробелов на один
	reSpace := regexp.MustCompile(`\s+`)
	cleaned = reSpace.ReplaceAllString(cleaned, " ")
	logger.Debugw("After normalize spaces", "result", cleaned)

	// Убираем пробелы в начале и конце
	cleaned = strings.TrimSpace(cleaned)
	logger.Debugw("Final", "result", cleaned)

	return cleaned
}
