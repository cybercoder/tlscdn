package v1alpha1

import "k8s.io/apimachinery/pkg/runtime/schema"

var (
	HTTPRouteGVR = schema.GroupVersionResource{
		Group:    "cdn.ik8s.ir",
		Version:  "v1alpha1",
		Resource: "cdnhttproutes",
	}

	GatewayGVR = schema.GroupVersionResource{
		Group:    "cdn.ik8s.ir",
		Version:  "v1alpha1",
		Resource: "cdngateways",
	}

	WafRuleGVR = schema.GroupVersionResource{
		Group:    "cdn.ik8s.ir",
		Version:  "v1alpha1",
		Resource: "wafrules",
	}
)
