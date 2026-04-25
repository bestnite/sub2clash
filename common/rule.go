package common

import (
	"fmt"
	"strings"
)

func PrependRuleProvider(
	sub *generatedConfig, providerName string, group string,
) {
	PrependRules(
		sub,
		fmt.Sprintf("RULE-SET,%s,%s", providerName, group),
	)
}

func AppenddRuleProvider(
	sub *generatedConfig, providerName string, group string,
) {
	AppendRules(sub, fmt.Sprintf("RULE-SET,%s,%s", providerName, group))
}

// PrependRules 用于在规则头部插入新规则。
// 这通常对应用户显式要求 prepend 的场景。
func PrependRules(sub *generatedConfig, rules ...string) {
	if sub.Rule == nil {
		sub.Rule = make([]string, 0)
	}
	sub.Rule = append(rules, sub.Rule...)
}

// AppendRules 在规则尾部追加，但如果尾部已有 MATCH，则保持 MATCH 仍然是最后一条。
func AppendRules(sub *generatedConfig, rules ...string) {
	if sub.Rule == nil {
		sub.Rule = make([]string, 0)
	}
	if len(sub.Rule) == 0 {
		sub.Rule = append(sub.Rule, rules...)
		return
	}
	matchRule := sub.Rule[len(sub.Rule)-1]
	if strings.Contains(matchRule, "MATCH") {
		sub.Rule = append(sub.Rule[:len(sub.Rule)-1], rules...)
		sub.Rule = append(sub.Rule, matchRule)
		return
	}
	sub.Rule = append(sub.Rule, rules...)
}
