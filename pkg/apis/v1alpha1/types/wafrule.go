package v1alpha1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// --- Rule Phases ---
type RulePhase int

const (
	PhaseRequestHeaders  RulePhase = 1
	PhaseRequestBody     RulePhase = 2
	PhaseResponseHeaders RulePhase = 3
	PhaseResponseBody    RulePhase = 4
	PhaseLogging         RulePhase = 5
)

// --- Rule Actions ---
type RuleAction string

const (
	ActionAllow      RuleAction = "allow"
	ActionDeny       RuleAction = "deny"
	ActionDrop       RuleAction = "drop"
	ActionPass       RuleAction = "pass"
	ActionProxy      RuleAction = "proxy"
	ActionRedirect   RuleAction = "redirect"
	ActionChain      RuleAction = "chain"
	ActionLog        RuleAction = "log"
	ActionAuditLog   RuleAction = "auditlog"
	ActionStatus     RuleAction = "status"
	ActionCapture    RuleAction = "capture"
	ActionTx         RuleAction = "t"
	ActionMsg        RuleAction = "msg"
	ActionTag        RuleAction = "tag"
	ActionLogData    RuleAction = "logdata"
	ActionSeverity   RuleAction = "severity"
	ActionMultiMatch RuleAction = "multiMatch"
	ActionVer        RuleAction = "ver"
	ActionRev        RuleAction = "rev"
	ActionID         RuleAction = "id"
	ActionSkip       RuleAction = "skip"
	ActionSkipAfter  RuleAction = "skipAfter"
	ActionCtl        RuleAction = "ctl"
	ActionInitCol    RuleAction = "initcol"
	ActionSetEnv     RuleAction = "setenv"
	ActionSetVar     RuleAction = "setvar"
)

// --- Transformations ---
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

// --- Operators ---
type RuleOperator string

// --- WAFRule ---
type WAFRule struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:",inline"`

	Spec   WAFRuleSpec    `json:"spec"`
	Status *WAFRuleStatus `json:"status,omitempty"`
}

type WAFRuleSpec struct {
	Gateway         *GatewayReference `json:"gateway,omitempty"`
	RuleID          int               `json:"ruleId"`
	Phase           RulePhase         `json:"phase"`
	Action          RuleAction        `json:"action"`
	Operator        RuleOperator      `json:"operator,omitempty"`
	Match           string            `json:"match,omitempty"`
	Variables       []RuleVariable    `json:"variables,omitempty"`
	Transformations []Transformation  `json:"transformations,omitempty"`
	Metadata        *RuleMetadata     `json:"metadata,omitempty"`
	Enabled         bool              `json:"enabled,omitempty"`
	HTTPStatus      int               `json:"httpStatus,omitempty"`
	ActionConfig    *ActionConfig     `json:"actionConfig,omitempty"`
	ChainRules      []WAFRuleSpec     `json:"chainRules,omitempty"` // Inline recursion
}

// Gateway reference
type GatewayReference struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace,omitempty"`
}

// Rule variable
type RuleVariable struct {
	Type      string `json:"type"`
	Selector  string `json:"selector,omitempty"`
	Count     bool   `json:"count,omitempty"`
	Exclusion bool   `json:"exclusion,omitempty"`
	Key       string `json:"key,omitempty"`
}

// Metadata for severity/message/tags
type RuleMetadata struct {
	Severity string   `json:"severity,omitempty"`
	Message  string   `json:"message,omitempty"`
	Tags     []string `json:"tags,omitempty"`
}

// Action config
type ActionConfig struct {
	Redirect  *RedirectAction `json:"redirect,omitempty"`
	Status    *StatusAction   `json:"status,omitempty"`
	SetVar    *SetVarAction   `json:"setvar,omitempty"`
	InitCol   *InitColAction  `json:"initcol,omitempty"`
	Skip      int             `json:"skip,omitempty"`
	SkipAfter string          `json:"skipAfter,omitempty"`
	LogData   string          `json:"logdata,omitempty"`
}

// Nested actions
type RedirectAction struct {
	URL string `json:"url"`
}
type StatusAction struct {
	Code int `json:"code"`
}
type SetVarAction struct {
	Variable string `json:"variable"`
	Value    string `json:"value"`
	Scope    string `json:"scope,omitempty"`
}
type InitColAction struct {
	Collection string `json:"collection"`
	Variable   string `json:"variable,omitempty"`
}

// Status
type WAFRuleStatus struct {
	Conditions      []metav1.Condition `json:"conditions,omitempty"`
	LastAppliedHash string             `json:"lastAppliedHash,omitempty"`
	Synced          bool               `json:"synced,omitempty"`
}

// List
type WAFRuleList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"inline"`
	Items           []WAFRule `json:"items"`
}
