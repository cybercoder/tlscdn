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
	ActionAllow    RuleAction = "allow"
	ActionDeny     RuleAction = "deny"
	ActionDrop     RuleAction = "drop"
	ActionPass     RuleAction = "pass"
	ActionProxy    RuleAction = "proxy"
	ActionRedirect RuleAction = "redirect"
	ActionChain    RuleAction = "chain"
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
	OperatorContains           RuleOperator = "contains"
	OperatorBeginsWith         RuleOperator = "beginsWith"
	OperatorEndsWith           RuleOperator = "endsWith"
	OperatorEquals             RuleOperator = "equals"
	OperatorRegex              RuleOperator = "regex"
	OperatorGeoLookup          RuleOperator = "geoLookup"
	OperatorIPMatch            RuleOperator = "ipMatch"
	OperatorUnconditionalMatch RuleOperator = "unconditionalMatch"
	OperatorWithin             RuleOperator = "within"
	OperatorDetectSQLi         RuleOperator = "detectSQLi"
	OperatorDetectXSS          RuleOperator = "detectXSS"
)

// WAFRule represents a single Coraza rule with support for chaining
type WAFRule struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              WAFRuleSpec   `json:"spec"`
	Status            WAFRuleStatus `json:"status,omitempty"`
}

type WAFRuleSpec struct {
	Gateway  *GatewayReference `json:"gateway,omitempty"` // Changed from LocalObjectReference
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
	Enabled    bool         `json:"enabled,omitempty"`    // Changed from Status to Enabled
	HTTPStatus int          `json:"httpStatus,omitempty"` // Added HTTPStatus field

	// Chaining support
	Chain      bool               `json:"chain,omitempty"`
	ChainRules []WAFRuleChainLink `json:"chainRules,omitempty"` // Changed to ChainLink type
}

// GatewayReference represents a reference to a Gateway resource
type GatewayReference struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace,omitempty"`
}

// WAFRuleChainLink represents a chained rule (simplified structure)
type WAFRuleChainLink struct {
	RuleID          int              `json:"ruleId,omitempty"`
	Phase           RulePhase        `json:"phase,omitempty"`
	Operator        RuleOperator     `json:"operator,omitempty"`
	Variables       []RuleVariable   `json:"variables,omitempty"`
	Match           string           `json:"match,omitempty"`
	Transformations []Transformation `json:"transformations,omitempty"`
}

type RuleVariable struct {
	Type     string `json:"type"` // Required field
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
