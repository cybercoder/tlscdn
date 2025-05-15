package events

import "log"

func OnAddHTTPRoute(obj interface{}) {
	log.Printf("httproute: %v", obj)
}
