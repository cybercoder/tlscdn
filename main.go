package main

import (
	"github.com/cybercoder/tlscdn-controller/pkg/events"
	"github.com/cybercoder/tlscdn-controller/pkg/k8s"
	"k8s.io/client-go/tools/cache"
)

func main() {
	gatewayInformer := k8s.CreateGatewayInformer()
	gatewayInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: events.OnAddGateway,
	})

	httpRouteInformer := k8s.CreateHTTPRouteInformer()
	httpRouteInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: events.OnAddHTTPRoute,
	})

	stopCh := make(chan struct{})

	defer close(stopCh)

	go gatewayInformer.Run(stopCh)
	go httpRouteInformer.Run(stopCh)

	select {}
}
