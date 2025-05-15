package k8s

import (
	"time"

	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic/dynamicinformer"
	"k8s.io/client-go/tools/cache"
)

var gatewayInformer cache.SharedIndexInformer
var httprouteInformer cache.SharedIndexInformer

var gatewayGVR = schema.GroupVersionResource{
	Group:    "gateway.cdn.ik8s.ir",
	Version:  "v1alpha1",
	Resource: "gateways",
}

var httprouteGVR = schema.GroupVersionResource{
	Group:    "gateway.cdn.ik8s.ir",
	Version:  "v1",
	Resource: "httproutes",
}

func CreateGatewayInformer() cache.SharedIndexInformer {
	if gatewayInformer != nil {
		return gatewayInformer
	}
	informerFactory := dynamicinformer.NewFilteredDynamicSharedInformerFactory(CreateDynamicClient(), time.Second*30, "", nil)
	gatewayInformer = informerFactory.ForResource(gatewayGVR).Informer()
	return gatewayInformer
}

func CreateHTTPRouteInformer() cache.SharedIndexInformer {
	if httprouteInformer != nil {
		return httprouteInformer
	}
	informerFactory := dynamicinformer.NewFilteredDynamicSharedInformerFactory(CreateDynamicClient(), time.Second*30, "", nil)
	httprouteInformer = informerFactory.ForResource(httprouteGVR).Informer()
	return httprouteInformer
}
