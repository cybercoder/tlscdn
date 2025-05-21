package main

import (
	"log"

	"github.com/cybercoder/tlscdn-controller/pkg/events"
	"github.com/cybercoder/tlscdn-controller/pkg/k8s"
	"github.com/joho/godotenv"
	"k8s.io/client-go/tools/cache"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Error loading .env file")
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

	stopCh := make(chan struct{})

	defer close(stopCh)

	go gatewayInformer.Run(stopCh)
	go httpRouteInformer.Run(stopCh)

	select {}
}
