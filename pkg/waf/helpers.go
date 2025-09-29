package waf

import (
	"fmt"
	"strings"

	v1alpha1Types "github.com/cybercoder/tlscdn-controller/pkg/apis/v1alpha1/types"
)

// SecLang entry point
func RuleToSecLang(rule *v1alpha1Types.WAFRule) string {
	return ruleLogicToSecLang(
		rule.Spec.RuleLogic,
		rule.Spec.RuleID,
		int(rule.Spec.Phase),
		rule.Spec.Action,
		rule.Spec.Enabled,
		rule.Spec.HTTPStatus,
		true,
	)
}

// Recursive converter
func ruleLogicToSecLang(rl v1alpha1Types.RuleLogic, ruleID int, phase int, action v1alpha1Types.RuleAction, enabled bool, httpStatus int, isRoot bool) string {
	var lines []string

	// Leaf MATCH
	if rl.Type == "MATCH" {
		var variables []string
		for _, v := range rl.Variables {
			if v.Selector != "" {
				variables = append(variables, fmt.Sprintf("%s:%s", v.Type, v.Selector))
			} else {
				variables = append(variables, v.Type)
			}
		}
		variablesPart := strings.Join(variables, "|")
		if variablesPart == "" {
			variablesPart = "REQUEST_URI"
		}

		operator := rl.Operator
		if strings.HasPrefix(operator, "@") {
			operator = operator[1:]
		}

		opPattern := operator
		if rl.Match != "" {
			opPattern = fmt.Sprintf("@%s %s", operator, escapeSecLangString(rl.Match))
		} else {
			opPattern = fmt.Sprintf("@%s", operator)
		}

		var actions []string
		if isRoot {
			actions = append(actions, fmt.Sprintf("id:%d", ruleID))
			actions = append(actions, fmt.Sprintf("phase:%d", phase))
			actions = append(actions, string(action))
			if httpStatus != 0 {
				actions = append(actions, fmt.Sprintf("status:%d", httpStatus))
			}
			if !enabled {
				actions = append(actions, "ctl:ruleEngine=Off")
			}
		} else {
			actions = append(actions, "chain")
		}

		for _, t := range rl.Transformations {
			actions = append(actions, fmt.Sprintf("t:%s", string(t)))
		}

		line := fmt.Sprintf("SecRule %s \"%s\" \"%s\"", variablesPart, opPattern, strings.Join(actions, ","))
		lines = append(lines, line)
		return strings.Join(lines, "\n")
	}

	// Logic node (AND/OR)
	if rl.Type == "AND" || rl.Type == "OR" {
		for i, child := range rl.Children {
			childIsRoot := isRoot
			if i > 0 {
				childIsRoot = false // remaining children use chain
			}
			lines = append(lines, ruleLogicToSecLang(child, ruleID, int(phase), action, enabled, httpStatus, childIsRoot))
		}
	}

	return strings.Join(lines, "\n")
}

// Escaping helper
func escapeSecLangString(s string) string {
	escaped := strings.ReplaceAll(s, "\\", "\\\\")
	escaped = strings.ReplaceAll(escaped, "\"", "\\\"")
	return escaped
}
