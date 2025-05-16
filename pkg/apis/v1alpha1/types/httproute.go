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
	HostHeader   string                `json:"host_header,omitempty"`
	LbMethod     string                `json:"lb_method"`
	Path         Path                  `json:"path,omitempty"`
}

type HTTPRouteStatus struct {
	// Conditions []metav1.Condition `json:"conditions,omitempty"`
}
