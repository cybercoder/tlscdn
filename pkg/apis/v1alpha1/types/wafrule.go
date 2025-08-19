package v1alpha1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

type WAFRule struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              WAFRuleSpec   `json:"spec"`
	Status            WAFRuleStatus `json:"status,omitempty"`
}

type WAFRuleSpec struct {
	Gateway         *LocalObjectReference `json:"gateway,omitempty"`
	RuleID          int                   `json:"ruleId"`
	Phase           int                   `json:"phase"`
	Action          string                `json:"action"`
	Operator        string                `json:"operator,omitempty"`
	Variables       []RuleVariable        `json:"variables,omitempty"`
	Transformations []string              `json:"transformations,omitempty"`
	Match           string                `json:"match,omitempty"`
	Metadata        RuleMetadata          `json:"metadata,omitempty"`
	Status          int                   `json:"status,omitempty"`
}

type RuleVariable struct {
	Name     string `json:"name"`
	Selector string `json:"selector,omitempty"`
}

type RuleMetadata struct {
	Severity string   `json:"severity,omitempty"`
	Message  string   `json:"message,omitempty"`
	Tags     []string `json:"tags,omitempty"`
}

type WAFRuleStatus struct {
	Conditions      []metav1.Condition `json:"conditions,omitempty"`
	LastAppliedHash string             `json:"lastAppliedHash,omitempty"`
	Synced          bool               `json:"synced,omitempty"`
}

type WAFRuleList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []WAFRule `json:"items"`
}
