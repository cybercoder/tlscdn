package k8s

import (
	"github.com/cybercoder/tlscdn-controller/pkg/apis/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic/dynamicinformer"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/tools/cache"
)

var gatewayInformer cache.SharedIndexInformer
var upstreamInformer cache.SharedIndexInformer
var httprouteInformer cache.SharedIndexInformer
var wafRuleInformer cache.SharedIndexInformer
var secretInformer cache.SharedIndexInformer

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

func CreateSecretInformer() cache.SharedIndexInformer {
	if secretInformer != nil {
		return secretInformer
	}
	informerfactory := informers.NewSharedInformerFactory(CreateClient(), 0)
	secretInformer = informerfactory.Core().V1().Secrets().Informer()
	return secretInformer
}

func CreateWafRuleInformer() cache.SharedIndexInformer {
	if wafRuleInformer != nil {
		return wafRuleInformer
	}
	informerFactory := dynamicinformer.NewFilteredDynamicSharedInformerFactory(CreateDynamicClient(), 0, "", nil)
	wafRuleInformer = informerFactory.ForResource(v1alpha1.WafRuleGVR).Informer()
	return wafRuleInformer
}
