package v1alpha1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

type HTTPRoute struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              HTTPRouteSpec   `json:"spec"`
	Status            HTTPRouteStatus `json:"status,omitempty"`
}

type HTTPRouteSpec struct {
	Gateway      *LocalObjectReference `json:"gateway,omitempty"`
	UpstreamName string                `json:"upstreamName"`
	HostHeader   string                `json:"hostHeader,omitempty"`
	LbMethod     string                `json:"lbMethod"`
	Path         Path                  `json:"path,omitempty"`
	Cache        HTTPRouteSpecCache    `json:"cache,omitempty"`
}

type HTTPRouteSpecCache struct {
	Level                       string `json:"level"`
	IgnoreUpstreamCacheSettings bool   `json:"ignoreUpstreamCacheSettings"`
	BrowserTTL                  int    `json:"browserTTL"`
	EdgeTTL                     int    `json:"edgeTTL"`
	NonSuccessTTL               int    `json:"nonSuccessTTL"`
	StaleTTL                    int    `json:"staleTTL"`
	Immutable                   bool   `json:"immutable"`
}

type HTTPRouteStatus struct {
	// Conditions []metav1.Condition `json:"conditions,omitempty"`
	RedisKey string `json:"redisKey,omitempty"`
}

type HTTPRouteList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []HTTPRoute `json:"items"`
}
