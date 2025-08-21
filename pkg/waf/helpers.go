package waf

import (
	"fmt"
	"strings"

	v1alpha1Types "github.com/cybercoder/tlscdn-controller/pkg/apis/v1alpha1/types"
	"github.com/samber/lo"
)

func RuleToSecLang(rule *v1alpha1Types.WAFRule) string {
	// Build variables part
	var variableStrings []string
	if len(rule.Spec.Variables) > 0 {
		variableStrings = lo.Map(rule.Spec.Variables, func(v v1alpha1Types.RuleVariable, _ int) string {
			if v.Selector != "" {
				return fmt.Sprintf("%s:%s", v.Type, v.Selector)
			}
			return v.Type
		})
	}
	variablesPart := strings.Join(variableStrings, "|")
	if variablesPart == "" {
		variablesPart = "REQUEST_URI"
	}

	// Build operator with pattern
	operatorWithPattern := string(rule.Spec.Operator)
	if rule.Spec.Match != "" {
		operatorWithPattern = fmt.Sprintf("@%s %s", rule.Spec.Operator, escapeSecLangString(rule.Spec.Match))
	} else {
		operatorWithPattern = fmt.Sprintf("@%s", rule.Spec.Operator)
	}

	// Build actions
	var actions []string
	actions = append(actions, fmt.Sprintf("id:%d", rule.Spec.RuleID))
	actions = append(actions, fmt.Sprintf("phase:%d", rule.Spec.Phase))

	if rule.Spec.Chain {
		actions = append(actions, "chain")
	} else {
		actions = append(actions, string(rule.Spec.Action))
	}

	// Add transformations
	if len(rule.Spec.Transformations) > 0 {
		transforms := lo.Map(rule.Spec.Transformations, func(t v1alpha1Types.Transformation, _ int) string {
			return fmt.Sprintf("t:%s", string(t))
		})
		actions = append(actions, strings.Join(transforms, ","))
	}

	// Add metadata
	if rule.Spec.Metadata.Severity != "" {
		actions = append(actions, fmt.Sprintf("severity:%s", rule.Spec.Metadata.Severity))
	}
	if rule.Spec.Metadata.Message != "" {
		actions = append(actions, fmt.Sprintf("msg:%s", escapeSecLangString(rule.Spec.Metadata.Message)))
	}

	if len(rule.Spec.Metadata.Tags) > 0 {
		for _, tag := range rule.Spec.Metadata.Tags {
			actions = append(actions, fmt.Sprintf("tag:%s", tag))
		}
	}

	// Handle enabled status
	if !rule.Spec.Enabled {
		actions = append(actions, "ctl:ruleEngine=Off")
	}

	// Handle HTTP status
	if rule.Spec.HTTPStatus != 0 {
		actions = append(actions, fmt.Sprintf("status:%d", rule.Spec.HTTPStatus))
	}

	// Build the final rule
	actionsPart := strings.Join(actions, ",")
	ruleLine := fmt.Sprintf("SecRule %s \"%s\" \"%s\"", variablesPart, operatorWithPattern, actionsPart)

	// Handle chain rules using the helper function
	if len(rule.Spec.ChainRules) > 0 {
		chainRules := lo.Map(rule.Spec.ChainRules, func(cr v1alpha1Types.WAFRuleChainLink, _ int) string {
			return chainRuleToSecLang(cr)
		})
		return strings.Join(append([]string{ruleLine}, chainRules...), "\n")
	}

	return ruleLine
}

// Helper function for chain rules (simplified version without ID/phase/metadata)
func chainRuleToSecLang(chainLink v1alpha1Types.WAFRuleChainLink) string {
	// Build variables part
	var variableStrings []string
	if len(chainLink.Variables) > 0 {
		variableStrings = lo.Map(chainLink.Variables, func(v v1alpha1Types.RuleVariable, _ int) string {
			if v.Selector != "" {
				return fmt.Sprintf("%s:%s", v.Type, v.Selector)
			}
			return v.Type
		})
	}
	variablesPart := strings.Join(variableStrings, "|")
	if variablesPart == "" {
		variablesPart = "REQUEST_URI"
	}

	// Build operator with pattern
	operatorWithPattern := string(chainLink.Operator)
	if chainLink.Match != "" {
		operatorWithPattern = fmt.Sprintf("@%s %s", chainLink.Operator, escapeSecLangString(chainLink.Match))
	} else {
		operatorWithPattern = fmt.Sprintf("@%s", chainLink.Operator)
	}

	// Build actions - chain rules only need chain action and transformations
	var actions []string
	actions = append(actions, "chain") // All chain rules must have this

	// Add transformations
	if len(chainLink.Transformations) > 0 {
		transforms := lo.Map(chainLink.Transformations, func(t v1alpha1Types.Transformation, _ int) string {
			return fmt.Sprintf("t:%s", string(t))
		})
		actions = append(actions, strings.Join(transforms, ","))
	}

	actionsPart := strings.Join(actions, ",")

	return fmt.Sprintf("SecRule %s \"%s\" \"%s\"", variablesPart, operatorWithPattern, actionsPart)
}

func escapeSecLangString(s string) string {
	escaped := strings.ReplaceAll(s, "\\", "\\\\")
	escaped = strings.ReplaceAll(escaped, "\"", "\\\"")
	return escaped
}
