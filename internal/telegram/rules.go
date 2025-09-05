package telegram

import (
	"sort"
	"strings"

	"github.com/st-kuptsov/balabol/config"
)

type hit struct {
	pos      int
	ruleIdx  int
	resp     string
	ruleName string
	ruleText string
}

// MatchRules проверяет текст на совпадения с правилами в заданном режиме
func MatchRules(text string, rules []config.Rule, mode string) []hit {
	var hits []hit
	switch mode {
	case "first_last":
		words := strings.Fields(text)
		if len(words) == 0 {
			return nil
		}
		for i, rule := range rules {
			re := rule.Re()
			if re == nil {
				continue
			}
			if re.MatchString(words[0]) {
				hits = append(hits, hit{pos: 0, ruleIdx: i, resp: rule.Response, ruleName: rule.Pattern, ruleText: rule.Text})
			}
			if len(words) > 1 && re.MatchString(words[len(words)-1]) {
				hits = append(hits, hit{pos: len(text), ruleIdx: i, resp: rule.Response, ruleName: rule.Pattern, ruleText: rule.Text})
			}
		}

	case "all":
		for i, rule := range rules {
			re := rule.Re()
			if re == nil {
				continue
			}
			locs := re.FindAllStringIndex(text, -1)
			for _, loc := range locs {
				hits = append(hits, hit{pos: loc[0], ruleIdx: i, resp: rule.Response, ruleName: rule.Pattern, ruleText: rule.Text})
			}
		}

	default: // fallback
		for i, rule := range rules {
			re := rule.Re()
			if re == nil {
				continue
			}
			locs := re.FindAllStringIndex(text, -1)
			for _, loc := range locs {
				hits = append(hits, hit{pos: loc[0], ruleIdx: i, resp: rule.Response, ruleName: rule.Pattern, ruleText: rule.Text})
			}
		}
	}

	// сортировка по позиции в тексте
	sort.Slice(hits, func(i, j int) bool {
		if hits[i].pos == hits[j].pos {
			return hits[i].ruleIdx < hits[j].ruleIdx
		}
		return hits[i].pos < hits[j].pos
	})

	return hits
}
