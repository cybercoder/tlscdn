package k8s

import (
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic/dynamicinformer"
	"k8s.io/client-go/tools/cache"
)

var gatewayInformer cache.SharedIndexInformer
var upstreamInformer cache.SharedIndexInformer
var httprouteInformer cache.SharedIndexInformer

var gatewayGVR = schema.GroupVersionResource{
	Group:    "cdn.ik8s.ir",
	Version:  "v1alpha1",
	Resource: "cdngateways",
}

var httprouteGVR = schema.GroupVersionResource{
	Group:    "cdn.ik8s.ir",
	Version:  "v1alpha1",
	Resource: "cdnhttproutes",
}

var upstreamGVR = schema.GroupVersionResource{
	Group:    "cdn.ik8s.ir",
	Version:  "v1alpha1",
	Resource: "upstreams",
}

func CreateGatewayInformer() cache.SharedIndexInformer {
	if gatewayInformer != nil {
		return gatewayInformer
	}
	informerFactory := dynamicinformer.NewFilteredDynamicSharedInformerFactory(CreateDynamicClient(), 0, "", nil)
	gatewayInformer = informerFactory.ForResource(gatewayGVR).Informer()
	return gatewayInformer
}

func CreateHTTPRouteInformer() cache.SharedIndexInformer {
	if httprouteInformer != nil {
		return httprouteInformer
	}
	informerFactory := dynamicinformer.NewFilteredDynamicSharedInformerFactory(CreateDynamicClient(), 0, "", nil)
	httprouteInformer = informerFactory.ForResource(httprouteGVR).Informer()
	return httprouteInformer
}
