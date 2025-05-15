package main

import (
	"github.com/cybercoder/tlscdn/pkg/events"
	"github.com/cybercoder/tlscdn/pkg/k8s"
	"k8s.io/client-go/tools/cache"
)

func main() {
	gatewayInformer := k8s.CreateGatewayInformer()
	gatewayInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: events.OnAddGateway,
	})

	stopCh := make(chan struct{})

	defer close(stopCh)

	go gatewayInformer.Run(stopCh)

	select {}
}
