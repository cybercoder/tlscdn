package v1alpha1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// WAFRule represents the custom resource
type WAFRule struct {
	metav1.TypeMeta    `json:",inline"`
	*metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec               WAFRuleSpec `json:"spec"`
}

type WAFRuleSpec struct {
	CdnGateway string `json:"gateway"`
	Enabled    bool   `json:"enabled,omitempty"`
	Rules      []Rule `json:"rules"`
}

type Rule struct {
	ID          string        `json:"id"`
	Name        string        `json:"name"`
	Enabled     bool          `json:"enabled,omitempty"`
	Description string        `json:"description,omitempty"`
	Groups      [][]Condition `json:"groups"`
	Action      Action        `json:"action"`
}

type Condition struct {
	Param     string `json:"param"`
	Operator  string `json:"operator"`
	Value     any    `json:"value"`
	ParamName string `json:"param_name,omitempty"`
}

type Action struct {
	Type    string `json:"type"`
	Message string `json:"message,omitempty"`
	Code    int    `json:"code,omitempty"`
	Rate    int    `json:"rate,omitempty"`
	Burst   int    `json:"burst,omitempty"`
	Name    string `json:"name,omitempty"`
	Value   string `json:"value,omitempty"`
	Expires int    `json:"expires,omitempty"`
}
