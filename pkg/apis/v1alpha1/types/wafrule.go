package v1alpha1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// Use enums for fixed sets of values to prevent invalid inputs
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

type Transformation string

const (
	TransformationNone               Transformation = "none"
	TransformationLowercase          Transformation = "lowercase"
	TransformationUrlDecode          Transformation = "urlDecode"
	TransformationHtmlEntityDecode   Transformation = "htmlEntityDecode"
	TransformationBase64Decode       Transformation = "base64Decode"
	TransformationRemoveWhitespace   Transformation = "removeWhitespace"
	TransformationRemoveNulls        Transformation = "removeNulls"
	TransformationCmdLine            Transformation = "cmdLine"
	TransformationCompressWhitespace Transformation = "compressWhitespace"
	TransformationBase64Encode       Transformation = "base64Encode"
)

type RuleOperator string

const (
	OperatorContains             RuleOperator = "@contains"
	OperatorBeginsWith           RuleOperator = "@beginsWith"
	OperatorEndsWith             RuleOperator = "@endsWith"
	OperatorEquals               RuleOperator = "@equals"
	OperatorRegex                RuleOperator = "@rx"
	OperatorGeoLookup            RuleOperator = "@geoLookup"
	OperatorIPMatch              RuleOperator = "@ipMatch"
	OperatorUnconditionalMatch   RuleOperator = "@unconditionalMatch"
	OperatorWithin               RuleOperator = "@within"
	OperatorDetectSQLi           RuleOperator = "@detectSQLi"
	OperatorDetectXSS            RuleOperator = "@detectXSS"
	OperatorPm                   RuleOperator = "@pm"
	OperatorPmFromFile           RuleOperator = "@pmFromFile"
	OperatorStreq                RuleOperator = "@streq"
	OperatorValidateByteRange    RuleOperator = "@validateByteRange"
	OperatorValidateUrlEncoding  RuleOperator = "@validateUrlEncoding"
	OperatorValidateUtf8Encoding RuleOperator = "@validateUtf8Encoding"
	OperatorVerifyCC             RuleOperator = "@verifyCC"
	OperatorVerifyCPF            RuleOperator = "@verifyCPF"
	OperatorVerifyCNPJ           RuleOperator = "@verifyCNPJ"
	OperatorVerifySSN            RuleOperator = "@verifySSN"
)

// WAFRule represents a single Coraza rule with support for chaining
type WAFRule struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              WAFRuleSpec   `json:"spec"`
	Status            WAFRuleStatus `json:"status,omitempty"`
}

type WAFRuleSpec struct {
	Gateway  *GatewayReference `json:"gateway,omitempty"`
	RuleID   int               `json:"ruleId"`
	Phase    RulePhase         `json:"phase"`
	Action   RuleAction        `json:"action"`
	Operator RuleOperator      `json:"operator,omitempty"`
	Match    string            `json:"match,omitempty"`

	// Variables to inspect
	Variables []RuleVariable `json:"variables,omitempty"`

	// Transformations to apply to variables before matching
	Transformations []Transformation `json:"transformations,omitempty"`

	Metadata   RuleMetadata `json:"metadata,omitempty"`
	Enabled    bool         `json:"enabled,omitempty"`
	HTTPStatus int          `json:"httpStatus,omitempty"`

	// Action-specific fields
	RedirectURL string         `json:"redirect,omitempty"`
	StatusCode  int            `json:"status,omitempty"`
	SetVar      *SetVarAction  `json:"setvar,omitempty"`
	InitCol     *InitColAction `json:"initcol,omitempty"`
	SkipCount   int            `json:"skip,omitempty"`
	SkipAfterID string         `json:"skipAfter,omitempty"`
	LogData     string         `json:"logdata,omitempty"`
	Severity    string         `json:"severity,omitempty"` // Overrides metadata.severity for severity action
	Tags        []string       `json:"tags,omitempty"`     // For tag action
	Message     string         `json:"message,omitempty"`  // For msg action

	// Chaining support
	Chain      bool               `json:"chain,omitempty"`
	ChainRules []WAFRuleChainLink `json:"chainRules,omitempty"`
}

// GatewayReference represents a reference to a Gateway resource
type GatewayReference struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace,omitempty"`
}

// WAFRuleChainLink represents a chained rule (simplified structure)
type WAFRuleChainLink struct {
	Operator        RuleOperator     `json:"operator"`
	Variables       []RuleVariable   `json:"variables,omitempty"`
	Match           string           `json:"match"`
	Transformations []Transformation `json:"transformations,omitempty"`
}

type RuleVariable struct {
	Type      string `json:"type"`                // REQUEST_HEADERS, ARGS, TX, IP, etc.
	Selector  string `json:"selector,omitempty"`  // Specific field name
	Count     bool   `json:"count,omitempty"`     // Use _COUNT version
	Exclusion bool   `json:"exclusion,omitempty"` // Use ! prefix for exclusion
	Key       string `json:"key,omitempty"`       // For collections like TX, IP
}

type RuleMetadata struct {
	Severity string   `json:"severity,omitempty"` // Default severity
	Message  string   `json:"message,omitempty"`  // Default message
	Tags     []string `json:"tags,omitempty"`     // Default tags
}

// SetVarAction represents setvar action parameters
type SetVarAction struct {
	Variable string `json:"variable"`
	Value    string `json:"value"`
	Scope    string `json:"scope,omitempty"` // TX, IP, GLOBAL, RESOURCE, SESSION
}

// InitColAction represents initcol action parameters
type InitColAction struct {
	Collection string `json:"collection"`
	Variable   string `json:"variable,omitempty"`
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
