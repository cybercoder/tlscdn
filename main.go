package main

import (
	"github.com/cybercoder/tlscdn-controller/pkg/events"
	"github.com/cybercoder/tlscdn-controller/pkg/k8s"
	"github.com/cybercoder/tlscdn-controller/pkg/logger"
	"github.com/joho/godotenv"
	"k8s.io/client-go/tools/cache"
)

func main() {
	// Initialize the logger
	logger.Init()

	err := godotenv.Load()
	if err != nil {
		logger.Warn("Error loading .env file")
	}
	gatewayInformer := k8s.CreateGatewayInformer()
	gatewayInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    events.OnAddGateway,
		UpdateFunc: events.OnUpdateGateway,
		DeleteFunc: events.OnDeleteGateway,
	})

	httpRouteInformer := k8s.CreateHTTPRouteInformer()
	httpRouteInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    events.OnAddHTTPRoute,
		UpdateFunc: events.OnUpdateHTTPRoute,
		DeleteFunc: events.OnDeleteHTTPRoute,
	})

	wafRulesInformer := k8s.CreateWafRuleInformer()
	wafRulesInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    events.OnAddWafRule,
		UpdateFunc: events.OnUpdateWafRule,
		DeleteFunc: events.OnDeleteWafRule,
	})

	secretInformer := k8s.CreateSecretInformer()
	secretInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    events.OnAddSecret,
		UpdateFunc: events.OnUpdateSecret,
		DeleteFunc: events.OnDeleteSecret,
	})

	stopCh := make(chan struct{})

	defer close(stopCh)

	go gatewayInformer.Run(stopCh)
	go httpRouteInformer.Run(stopCh)
	go wafRulesInformer.Run(stopCh)
	go secretInformer.Run(stopCh)
	select {}
}
