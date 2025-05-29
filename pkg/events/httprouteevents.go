package events

import (
	"context"
	"encoding/json"
	"os"
	"strings"

	v1alpha1 "github.com/cybercoder/tlscdn-controller/pkg/apis/v1alpha1"
	v1alpha1Types "github.com/cybercoder/tlscdn-controller/pkg/apis/v1alpha1/types"
	"github.com/cybercoder/tlscdn-controller/pkg/k8s"
	"github.com/cybercoder/tlscdn-controller/pkg/logger"
	"github.com/cybercoder/tlscdn-controller/pkg/redis"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
)

func OnAddHTTPRoute(obj interface{}) {
	u := obj.(*unstructured.Unstructured)
	var httproute v1alpha1Types.HTTPRoute
	if err := runtime.DefaultUnstructuredConverter.FromUnstructured(u.Object, &httproute); err != nil {
		logger.Errorf("Error converting unstructured to httproute object: %v", err)
		return
	}

	k := k8s.CreateDynamicClient()

	gName := httproute.Spec.Gateway.Name
	g, err := k.Resource(v1alpha1.GatewayGVR).Namespace(httproute.GetNamespace()).Get(context.TODO(), string(gName), metav1.GetOptions{})
	if err != nil {
		logger.Errorf("Gateway %s not found in namespace %s, orphaned route: %s, err: %v", gName, httproute.GetNamespace(), httproute.GetName(), err)
		return
	}
	var gateway v1alpha1Types.Gateway
	if err := runtime.DefaultUnstructuredConverter.FromUnstructured(g.Object, &gateway); err != nil {
		logger.Errorf("Error converting unstructured to gateway object: %v", err)
		return
	}

	var upstream v1alpha1Types.Upstream
	for _, u := range gateway.Spec.Upstreams {
		logger.Debugf("Upstream name: %s, HTTPRoute upstream: %s", u.Name, httproute.Spec.UpstreamName)
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
		"cache":     httproute.Spec.Cache,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		logger.Errorf("Failed to marshal data: %v", err)
	}
	CDN_HOSTNAME := os.Getenv("CDN_HOSTNAME")
	if CDN_HOSTNAME == "" {
		CDN_HOSTNAME = "cdntls.ir"
	}
	redisKey := ""
	if gateway.Spec.Domain == "" {
		redisKey = "httproute:" + strings.Replace(string(gateway.GetUID()), "-", "", -1) + "." + CDN_HOSTNAME + ":" + httproute.Spec.Path.Type + ":" + httproute.Spec.Path.Path
	} else {
		redisKey = "httproute:" + gateway.Spec.Domain + ":" + httproute.Spec.Path.Type + ":" + httproute.Spec.Path.Path
	}

	err = redisClient.Set(context.Background(),
		redisKey,
		jsonData, 0).Err()
	if err != nil {
		logger.Errorf("Error on storing httproute: %s to redis: %v", redisKey, err)
		return
	}

	err = unstructured.SetNestedField(u.Object, redisKey, "status", "redisKey")
	if err != nil {
		logger.Errorf("Failed to set redis key on httproute %s status: %v", httproute.GetName(), err)
		return
	}
	_, err = k.Resource(v1alpha1.HTTPRouteGVR).Namespace(httproute.GetNamespace()).UpdateStatus(context.TODO(), u, metav1.UpdateOptions{})
	if err != nil {
		logger.Errorf("Failed to update httproute %s status: %v", httproute.GetName(), err)
		return
	}
}

func OnUpdateHTTPRoute(prev interface{}, obj interface{}) {
	u := obj.(*unstructured.Unstructured)
	var httproute *v1alpha1Types.HTTPRoute
	if err := runtime.DefaultUnstructuredConverter.FromUnstructured(u.Object, &httproute); err != nil {
		logger.Errorf("Error converting unstructured to httproute object: %v", err)
		return
	}

	k := k8s.CreateDynamicClient()

	gName := httproute.Spec.Gateway.Name
	g, err := k.Resource(v1alpha1.GatewayGVR).Namespace(httproute.GetNamespace()).Get(context.TODO(), string(gName), metav1.GetOptions{})
	if err != nil {
		logger.Errorf("Gateway %s not found in namespace %s, orphaned route: %s, err: %v", gName, httproute.GetNamespace(), httproute.GetName(), err)
		return
	}
	var gateway v1alpha1Types.Gateway
	if err := runtime.DefaultUnstructuredConverter.FromUnstructured(g.Object, &gateway); err != nil {
		logger.Errorf("Error converting unstructured to gateway object: %v", err)
		return
	}

	var upstream v1alpha1Types.Upstream
	for _, u := range gateway.Spec.Upstreams {
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
		"cache":     httproute.Spec.Cache,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		logger.Errorf("Failed to marshal data: %v", err)
	}
	CDN_HOSTNAME := os.Getenv("CDN_HOSTNAME")
	if CDN_HOSTNAME == "" {
		CDN_HOSTNAME = "cdntls.ir"
	}
	redisKey := "httproute:" + strings.Replace(string(gateway.GetUID()), "-", "", -1) + "." + CDN_HOSTNAME + ":" + httproute.Spec.Path.Type + ":" + httproute.Spec.Path.Path
	err = redisClient.Set(context.Background(), redisKey, jsonData, 0).Err()
	if err != nil {
		logger.Errorf("Redis update set for httproute %s in namespace %s was unsuccessful: %v", httproute.GetName(), httproute.GetNamespace(), err)
		return
	}
	err = redisClient.Publish(context.Background(), "invalidate_httproute_cache", redisKey).Err()
	if err != nil {
		logger.Errorf("[Update httproute] cache invalidation, publish message for httproute %s in namespace %s was unsuccessful: %v", httproute.GetName(), httproute.GetNamespace(), err)
		return
	}

}

func OnDeleteHTTPRoute(obj interface{}) {
	u := obj.(*unstructured.Unstructured)
	var httproute *v1alpha1Types.HTTPRoute
	if err := runtime.DefaultUnstructuredConverter.FromUnstructured(u.Object, &httproute); err != nil {
		logger.Errorf("Error converting unstructured to httproute object: %v", err)
		return
	}
	redisClient := redis.CreateClient()
	err := redisClient.Del(context.Background(), httproute.Status.RedisKey).Err()
	if err != nil {
		logger.Errorf(
			"Error on deleting httproute %s in namespace %s key from redis: %v",
			u.GetName(), u.GetNamespace(), err,
		)
	}
	err = redisClient.Publish(context.Background(), "invalidate_httproute_cache", httproute.Status.RedisKey).Err()
	if err != nil {
		logger.Errorf("[Delete httproute] cache invalidation, publish message for httproute %s in namespace %s was unsuccessful: %v", httproute.GetName(), httproute.GetNamespace(), err)
		return
	}
}
