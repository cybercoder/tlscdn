package waf

import (
	"fmt"
	"strings"

	v1alpha1types "github.com/cybercoder/tlscdn-controller/pkg/apis/v1alpha1/types"
)

// Convert a WAFRule to SecLang string
func RuleToSecLang(rule *v1alpha1types.WAFRule) string {
	var variablesPart string
	if len(rule.Spec.Variables) > 0 {
		var parts []string
		for _, v := range rule.Spec.Variables {
			if v.Selector != "" {
				parts = append(parts, fmt.Sprintf("%s:%s", v.Type, v.Selector))
			} else {
				parts = append(parts, v.Type)
			}
		}
		variablesPart = strings.Join(parts, "|")
	} else {
		variablesPart = "REQUEST_URI"
	}

	operator := string(rule.Spec.Operator)
	if operator != "" && operator[0] == '@' {
		operator = operator[1:]
	}
	operatorWithPattern := operator
	if rule.Spec.Match != "" {
		operatorWithPattern = fmt.Sprintf("@%s %s", operator, escape(rule.Spec.Match))
	} else if operator != "" {
		operatorWithPattern = fmt.Sprintf("@%s", operator)
	}

	var actions []string
	actions = append(actions, fmt.Sprintf("id:%d", rule.Spec.RuleID))
	actions = append(actions, fmt.Sprintf("phase:%d", rule.Spec.Phase))
	actions = append(actions, string(rule.Spec.Action))

	for _, t := range rule.Spec.Transformations {
		actions = append(actions, fmt.Sprintf("t:%s", t))
	}

	if rule.Spec.Metadata.Severity != "" {
		actions = append(actions, fmt.Sprintf("severity:%s", rule.Spec.Metadata.Severity))
	}
	if rule.Spec.Metadata.Message != "" {
		actions = append(actions, fmt.Sprintf("msg:%s", escape(rule.Spec.Metadata.Message)))
	}
	for _, tag := range rule.Spec.Metadata.Tags {
		actions = append(actions, fmt.Sprintf("tag:%s", tag))
	}

	if !rule.Spec.Enabled {
		actions = append(actions, "ctl:ruleEngine=Off")
	}
	if rule.Spec.HTTPStatus != 0 {
		actions = append(actions, fmt.Sprintf("status:%d", rule.Spec.HTTPStatus))
	}

	// ActionConfig
	if rule.Spec.ActionConfig != nil {
		ac := rule.Spec.ActionConfig
		if ac.Redirect != nil && ac.Redirect.URL != "" {
			actions = append(actions, fmt.Sprintf("redirect:%s", escape(ac.Redirect.URL)))
		}
		if ac.Status != nil && ac.Status.Code != 0 {
			actions = append(actions, fmt.Sprintf("status:%d", ac.Status.Code))
		}
		if ac.SetVar != nil && ac.SetVar.Variable != "" {
			scope := ""
			if ac.SetVar.Scope != "" {
				scope = fmt.Sprintf("@%s", ac.SetVar.Scope)
			}
			actions = append(actions, fmt.Sprintf("setvar:%s=%s%s", ac.SetVar.Variable, ac.SetVar.Value, scope))
		}
		if ac.InitCol != nil && ac.InitCol.Collection != "" {
			actions = append(actions, fmt.Sprintf("initcol:%s=%s", ac.InitCol.Collection, ac.InitCol.Variable))
		}
		if ac.Skip > 0 {
			actions = append(actions, fmt.Sprintf("skip:%d", ac.Skip))
		}
		if ac.SkipAfter != "" {
			actions = append(actions, fmt.Sprintf("skipAfter:%s", ac.SkipAfter))
		}
		if ac.LogData != "" {
			actions = append(actions, fmt.Sprintf("logdata:%s", escape(ac.LogData)))
		}
	}

	actionsPart := strings.Join(actions, ",")
	ruleLine := fmt.Sprintf("SecRule %s \"%s\" \"%s\"", variablesPart, operatorWithPattern, actionsPart)

	// Chain rules
	if len(rule.Spec.ChainRules) > 0 {
		var lines []string
		lines = append(lines, ruleLine)
		for _, cr := range rule.Spec.ChainRules {
			lines = append(lines, chainRuleToSecLang(cr))
		}
		return strings.Join(lines, "\n")
	}

	return ruleLine
}

func chainRuleToSecLang(cr v1alpha1types.WAFRuleSpec) string {
	var variablesPart string
	if len(cr.Variables) > 0 {
		var parts []string
		for _, v := range cr.Variables {
			if v.Selector != "" {
				parts = append(parts, fmt.Sprintf("%s:%s", v.Type, v.Selector))
			} else {
				parts = append(parts, v.Type)
			}
		}
		variablesPart = strings.Join(parts, "|")
	} else {
		variablesPart = "REQUEST_URI"
	}

	operator := string(cr.Operator)
	if operator != "" && operator[0] == '@' {
		operator = operator[1:]
	}
	operatorWithPattern := operator
	if cr.Match != "" {
		operatorWithPattern = fmt.Sprintf("@%s %s", operator, escape(cr.Match))
	} else if operator != "" {
		operatorWithPattern = fmt.Sprintf("@%s", operator)
	}

	var actions []string
	actions = append(actions, "chain")
	for _, t := range cr.Transformations {
		actions = append(actions, fmt.Sprintf("t:%s", t))
	}

	actionsPart := strings.Join(actions, ",")
	return fmt.Sprintf("SecRule %s \"%s\" \"%s\"", variablesPart, operatorWithPattern, actionsPart)
}

func escape(s string) string {
	s = strings.ReplaceAll(s, "\\", "\\\\")
	s = strings.ReplaceAll(s, "\"", "\\\"")
	s = strings.ReplaceAll(s, "!", "\\!")
	s = strings.ReplaceAll(s, "/", "\\/")
	return s
}
