package telegram

import (
	"go.uber.org/zap"
	"strings"

	"github.com/st-kuptsov/balabol/config"
)

// hit представляет совпадение текста с правилом
type hit struct {
	pos      int    // позиция совпадения в тексте
	ruleIdx  int    // индекс правила в списке rules
	resp     string // ответ, связанный с правилом
	ruleName string // название правила (Pattern)
	ruleText string // текстовое описание правила
}

// MatchRules проверяет текст на соответствие правилам.
// Параметры:
// - text: текст для проверки
// - rules: список правил (config.Rule)
// - mode: режим обработки ("first_last" или "all")
// - logger: логгер для отладки
//
// Возвращает список hit — все совпадения с правилами.
func MatchRules(text string, rules []config.Rule, mode string, logger *zap.SugaredLogger) []hit {
	var hits []hit

	switch mode {
	case "first_last":
		// В режиме "first_last" проверяем только начало и конец текста
		for i, rule := range rules {
			re := rule.Re()
			if re == nil {
				logger.Debugw("Rule has nil regexp", "ruleText", rule.Text)
				continue
			}

			locs := re.FindAllStringIndex(text, -1)      // позиции совпадений
			matchedStrings := re.FindAllString(text, -1) // сами совпадения
			logger.Debugw("Pattern matching",
				"text", text,
				"pattern", rule.Pattern,
				"matches", locs,
				"matchedStrings", matchedStrings)

			if len(locs) == 0 {
				continue
			}

			// Если совпадение в начале текста
			if locs[0][0] == 0 {
				hits = append(hits, hit{
					pos:      locs[0][0],
					ruleIdx:  i,
					resp:     rule.Response,
					ruleName: rule.Pattern,
					ruleText: rule.Text,
				})
			}

			// Если совпадение в конце текста
			lastLoc := locs[len(locs)-1]
			if lastLoc[1] == len(text) || strings.TrimSpace(text[lastLoc[1]:]) == "" {
				if lastLoc[0] != 0 || len(locs) > 1 {
					hits = append(hits, hit{
						pos:      lastLoc[0],
						ruleIdx:  i,
						resp:     rule.Response,
						ruleName: rule.Pattern,
						ruleText: rule.Text,
					})
				}
			}
		}

	case "all":
		// В режиме "all" добавляем все совпадения правил
		for i, rule := range rules {
			re := rule.Re()
			if re == nil {
				logger.Debugw("Rule has nil regexp", "ruleText", rule.Text)
				continue
			}

			locs := re.FindAllStringIndex(text, -1)
			matchedStrings := re.FindAllString(text, -1)
			logger.Debugw("Pattern matching",
				"text", text,
				"pattern", rule.Pattern,
				"matches", locs,
				"matchedStrings", matchedStrings)

			for _, loc := range locs {
				hits = append(hits, hit{
					pos:      loc[0],
					ruleIdx:  i,
					resp:     rule.Response,
					ruleName: rule.Pattern,
					ruleText: rule.Text,
				})
			}
		}
	}

	return hits
}
