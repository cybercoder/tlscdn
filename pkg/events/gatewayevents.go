package events

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"strings"

	v1alpha1 "github.com/cybercoder/tlscdn-controller/pkg/apis/v1alpha1/types"
	"github.com/cybercoder/tlscdn-controller/pkg/redis"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
)

func OnAddGateway(obj interface{}) {
	u := obj.(*unstructured.Unstructured)
	var gateway v1alpha1.Gateway
	if err := runtime.DefaultUnstructuredConverter.FromUnstructured(u.Object, &gateway); err != nil {
		log.Printf("Error converting unstructured to gateway object: %v", err)
		return
	}

	redisClient := redis.CreateClient()
	data := map[string]interface{}{
		"name":      gateway.GetName(),
		"namespace": gateway.GetNamespace(),
		"UID":       gateway.GetUID(),
		"upstreams": gateway.Spec.Upstreams,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Printf("failed to marshal data: %v", err)
	}
	CDN_HOSTNAME := os.Getenv("CDN_HOSTNAME")
	if CDN_HOSTNAME == "" {
		CDN_HOSTNAME = "cdntls.ir"
	}
	status := redisClient.Set(context.Background(), strings.Replace(string(gateway.GetUID()), "-", "", -1)+".cdntls.ir", jsonData, 0)
	log.Printf("status %v", status)
}
