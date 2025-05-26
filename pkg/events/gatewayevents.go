package events

import (
	"context"
	"encoding/json"
	"os"
	"strings"
	"time"

	v1alpha1 "github.com/cybercoder/tlscdn-controller/pkg/apis/v1alpha1"
	v1alpha1Types "github.com/cybercoder/tlscdn-controller/pkg/apis/v1alpha1/types"
	"github.com/cybercoder/tlscdn-controller/pkg/k8s"
	"github.com/cybercoder/tlscdn-controller/pkg/logger"
	"github.com/cybercoder/tlscdn-controller/pkg/redis"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
)

func OnAddGateway(obj interface{}) {
	u := obj.(*unstructured.Unstructured)
	var gateway v1alpha1Types.Gateway
	if err := runtime.DefaultUnstructuredConverter.FromUnstructured(u.Object, &gateway); err != nil {
		logger.Errorf("Error converting unstructured to gateway object: %v", err)
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
		logger.Errorf("Failed to marshal data: %v", err)
	}
	CDN_HOSTNAME := os.Getenv("CDN_HOSTNAME")
	if CDN_HOSTNAME == "" {
		CDN_HOSTNAME = "tlscdn.ir"
	}
	redisKey := strings.Replace(string(gateway.GetUID()), "-", "", -1) + "." + CDN_HOSTNAME
	status := redisClient.Set(context.Background(), redisKey, jsonData, 0)
	logger.Infof("Redis set status for gateway %s: %v", gateway.GetName(), status)

}

func OnUpdateGateway(prev interface{}, obj interface{}) {
	gateway, err := convertUnstructToGateway(obj)
	if err != nil {
		logger.Errorf("Error converting unstructured to gateway object: %v", err)
		return
	}

	oldgateway, err := convertUnstructToGateway(prev)
	if err != nil {
		logger.Errorf("Error converting unstructured previous to gateway object: %v", err)
		return
	}

	// Compare specs by marshaling them to JSON
	if compareSpecs(gateway.Spec, oldgateway.Spec) {
		return
	}

	// find all httproutes associated to the gateway.
	httpRoutes, err := findHttpRoutesByGatewayName(gateway.GetName(), gateway.GetNamespace())
	if err != nil {
		logger.Errorf("Error finding httproutes by gateway name: %v", err)
		return
	}
	if len(httpRoutes.Items) == 0 {
		logger.Info("No associated httproutes found for gateway")
		return
	}
	k := k8s.CreateDynamicClient()

	for _, httpRoute := range httpRoutes.Items {
		deletedUpstream := true
		for _, upstreams := range gateway.Spec.Upstreams {
			if httpRoute.Spec.UpstreamName == upstreams.Name {
				deletedUpstream = false
			}
		}
		if deletedUpstream {
			httpRoute.Spec.UpstreamName = ""
		}

		annotations := httpRoute.GetAnnotations()
		if annotations == nil {
			annotations = make(map[string]string)
		}
		annotations["cdngateway.cdn.ik8s.ir/updated"] = time.Now().Format(time.RFC3339Nano)
		// labels := httpRoute.GetLabels()
		// if labels == nil {
		// 	labels = make(map[string]string)
		// }
		uo := &unstructured.Unstructured{
			Object: map[string]interface{}{
				"apiVersion": httpRoute.APIVersion,
				"kind":       httpRoute.Kind,
				"metadata": map[string]interface{}{
					"name":      httpRoute.Name,
					"namespace": httpRoute.Namespace,
				},
				"spec": httpRoute.Spec,
			},
		}
		uo.SetAnnotations(annotations)
		// uo.SetLabels(labels)
		uo.SetResourceVersion(httpRoute.GetResourceVersion())
		_, err = k.Resource(v1alpha1.HTTPRouteGVR).Namespace(httpRoute.GetNamespace()).Update(context.Background(), uo, metav1.UpdateOptions{})
		if err != nil {
			logger.Errorf("Error on httproute apply changes: %v", err)
			return
		}
		logger.Infof("Updated httproute %s for gateway %s", httpRoute.GetName(), gateway.GetName())
	}
}

func convertUnstructToGateway(obj interface{}) (*v1alpha1Types.Gateway, error) {
	u := obj.(*unstructured.Unstructured)
	var gateway v1alpha1Types.Gateway
	if err := runtime.DefaultUnstructuredConverter.FromUnstructured(u.Object, &gateway); err != nil {
		logger.Errorf("Error converting unstructured to gateway object: %v", err)
		return nil, err
	}
	return &gateway, nil
}

// compareSpecs compares two GatewaySpec objects by marshaling them to JSON
func compareSpecs(spec1, spec2 v1alpha1Types.GatewaySpec) bool {
	spec1JSON, err1 := json.Marshal(spec1)
	spec2JSON, err2 := json.Marshal(spec2)

	if err1 != nil || err2 != nil {
		logger.Errorf("Error marshaling specs: %v, %v", err1, err2)
		return false // If we can't compare, assume they're different
	}

	return string(spec1JSON) == string(spec2JSON)
}

func findHttpRoutesByGatewayName(name string, namespace string) (*v1alpha1Types.HTTPRouteList, error) {
	k := k8s.CreateDynamicClient()
	uglist, err := k.Resource(v1alpha1.HTTPRouteGVR).Namespace(namespace).List(context.TODO(), metav1.ListOptions{
		FieldSelector: "spec.gateway.name=" + name,
	})
	if err != nil {
		return nil, err
	}
	var httpRouteList v1alpha1Types.HTTPRouteList
	if err := runtime.DefaultUnstructuredConverter.FromUnstructured(uglist.UnstructuredContent(), &httpRouteList); err != nil {
		logger.Errorf("Error converting unstructured to HTTPRoute object: %v", err)
		return nil, err
	}
	return &httpRouteList, nil
}

func OnDeleteGateway(obj interface{}) {
	u := obj.(*unstructured.Unstructured)
	var gateway v1alpha1Types.Gateway
	if err := runtime.DefaultUnstructuredConverter.FromUnstructured(u.Object, &gateway); err != nil {
		logger.Errorf("Error converting unstructured to gateway object: %v", err)
		return
	}
	k := k8s.CreateDynamicClient()
	err := k.Resource(v1alpha1.HTTPRouteGVR).Namespace(gateway.GetNamespace()).DeleteCollection(
		context.Background(),
		metav1.DeleteOptions{},
		metav1.ListOptions{
			FieldSelector: "spec.gateway.name=" + gateway.GetName(),
		})
	if err != nil {
		logger.Errorf("Error on deleting associated httroutes for gateway %s on namespace %s: %v", gateway.GetName(), gateway.GetNamespace(), err)
	}
	redisClient := redis.CreateClient()
	CDN_HOSTNAME := os.Getenv("CDN_HOSTNAME")
	if CDN_HOSTNAME == "" {
		CDN_HOSTNAME = "tlscdn.ir"
	}
	redisKey := strings.Replace(string(gateway.GetUID()), "-", "", -1) + "." + CDN_HOSTNAME
	err = redisClient.Del(context.Background(), redisKey).Err()
	if err != nil {
		logger.Errorf("Error on deleting associated redis key for gateway %s on namespace %s: %v", gateway.GetName(), gateway.GetNamespace(), err)
	}
	err = redisClient.Publish(context.Background(), "invalidate_gateway_cache", redisKey).Err()
	if err != nil {
		logger.Errorf("[Delete gateway] cache invalidation, publish message for gateway %s in namespace %s was unsuccessful: %v", gateway.GetName(), gateway.GetNamespace(), err)
		return
	}
}
