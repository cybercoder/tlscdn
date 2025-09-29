package v1alpha1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

type RulePhase int

const (
	PhaseRequestHeaders  RulePhase = 1
	PhaseRequestBody     RulePhase = 2
	PhaseResponseHeaders RulePhase = 3
	PhaseResponseBody    RulePhase = 4
	PhaseLogging         RulePhase = 5
)

type RuleAction string

const (
	ActionAllow    RuleAction = "allow"
	ActionDeny     RuleAction = "deny"
	ActionDrop     RuleAction = "drop"
	ActionPass     RuleAction = "pass"
	ActionProxy    RuleAction = "proxy"
	ActionRedirect RuleAction = "redirect"
	ActionLog      RuleAction = "log"
	ActionAuditLog RuleAction = "auditlog"
	ActionStatus   RuleAction = "status"
	ActionCapture  RuleAction = "capture"
)

type Transformation string

const (
	TransNone               Transformation = "none"
	TransLowercase          Transformation = "lowercase"
	TransUrlDecode          Transformation = "urlDecode"
	TransHtmlEntityDecode   Transformation = "htmlEntityDecode"
	TransBase64Decode       Transformation = "base64Decode"
	TransRemoveWhitespace   Transformation = "removeWhitespace"
	TransRemoveNulls        Transformation = "removeNulls"
	TransCmdLine            Transformation = "cmdLine"
	TransCompressWhitespace Transformation = "compressWhitespace"
	TransBase64Encode       Transformation = "base64Encode"
)

// Leaf or logic node
type RuleLogic struct {
	Type            string           `json:"type"` // AND/OR/MATCH
	Operator        string           `json:"operator,omitempty"`
	Match           string           `json:"match,omitempty"`
	Variables       []RuleVariable   `json:"variables,omitempty"`
	Transformations []Transformation `json:"transformations,omitempty"`
	Children        []RuleLogic      `json:"children,omitempty"` // recursion
}

type RuleVariable struct {
	Type      string `json:"type"`
	Selector  string `json:"selector,omitempty"`
	Count     bool   `json:"count,omitempty"`
	Exclusion bool   `json:"exclusion,omitempty"`
	Key       string `json:"key,omitempty"`
}

type CdnGateway struct {
	Name string `json:"name"`
}

type WAFRuleSpec struct {
	CdnGateway CdnGateway `json:"cdnGateway"`
	RuleID     int        `json:"ruleId"`
	Phase      RulePhase  `json:"phase"`
	Action     RuleAction `json:"action"`
	Enabled    bool       `json:"enabled,omitempty"`
	HTTPStatus int        `json:"httpStatus,omitempty"`
	RuleLogic  RuleLogic  `json:"ruleLogic"`
}

type WAFRuleStatus struct {
	Conditions      []metav1.Condition `json:"conditions,omitempty"`
	LastAppliedHash string             `json:"lastAppliedHash,omitempty"`
	Synced          bool               `json:"synced,omitempty"`
}

type WAFRule struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   WAFRuleSpec   `json:"spec"`
	Status WAFRuleStatus `json:"status,omitempty"`
}

type WAFRuleList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []WAFRule `json:"items"`
}
