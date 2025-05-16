package v1alpha1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

type Gateway struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              GatewaySpec   `json:"spec"`
	Status            GatewayStatus `json:"status,omitempty"`
}

type GatewaySpec struct {
	Upstreams []Upstream `json:"upstreams,omitempty"`
}

type GatewayStatus struct {
	// ...
}

type GatewayList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Gateway `json:"items"`
}
