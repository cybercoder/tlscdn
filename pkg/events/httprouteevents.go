package events

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"strings"

	v1alpha1 "github.com/cybercoder/tlscdn-controller/pkg/apis/v1alpha1/types"
	"github.com/cybercoder/tlscdn-controller/pkg/k8s"
	"github.com/cybercoder/tlscdn-controller/pkg/redis"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func OnAddHTTPRoute(obj interface{}) {
	u := obj.(*unstructured.Unstructured)
	var httproute v1alpha1.HTTPRoute
	if err := runtime.DefaultUnstructuredConverter.FromUnstructured(u.Object, &httproute); err != nil {
		log.Printf("Error converting unstructured to httproute object: %v", err)
		return
	}

	k := k8s.CreateDynamicClient()
	gvr := schema.GroupVersionResource{
		Group:    "cdn.ik8s.ir",
		Version:  "v1alpha1",
		Resource: "cdngateways",
	}

	gName := httproute.Spec.Gateway.Name
	g, err := k.Resource(gvr).Namespace(httproute.GetNamespace()).Get(context.TODO(), string(gName), metav1.GetOptions{})
	if err != nil {
		log.Printf("Gateway %s not found in namespace %s, orphaned route: %s, err: %v", gName, httproute.GetNamespace(), httproute.GetName(), err)
		return
	}
	var gateway v1alpha1.Gateway
	if err := runtime.DefaultUnstructuredConverter.FromUnstructured(g.Object, &gateway); err != nil {
		log.Printf("Error converting unstructured to gateway object: %v", err)
		return
	}

	var upstream v1alpha1.Upstream
	for _, u := range gateway.Spec.Upstreams {
		log.Printf("u %s U %s", u.Name, httproute.Spec.UpstreamName)
		if u.Name == httproute.Spec.UpstreamName {
			upstream = u
			break
		}
	}
	redisClient := redis.CreateClient()

	if upstream.HostHeader != "" {
		for i := range upstream.Servers {
			upstream.Servers[i].HostHeader = upstream.HostHeader
		}
	}

	data := map[string]interface{}{
		"name":      httproute.GetName(),
		"namespace": httproute.GetNamespace(),
		"UID":       httproute.GetUID(),
		"lbMetod":   httproute.Spec.LbMethod,
		"gateway":   httproute.Spec.Gateway.Name,
		"upstreams": upstream.Servers,
	}

	log.Printf("upstreamServers: %v", upstream.Servers)

	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Printf("failed to marshal data: %v", err)
	}
	CDN_HOSTNAME := os.Getenv("CDN_HOSTNAME")
	if CDN_HOSTNAME == "" {
		CDN_HOSTNAME = "cdntls.ir"
	}
	status := redisClient.Set(context.Background(),
		"httproute:"+strings.Replace(string(gateway.GetUID()), "-", "", -1)+"."+CDN_HOSTNAME+":"+httproute.Spec.Path.Type+":"+httproute.Spec.Path.Path,
		jsonData, 0)
	log.Printf("status %v", status)

}
