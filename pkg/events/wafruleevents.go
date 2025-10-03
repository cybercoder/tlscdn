package events

import (
	"context"
	"encoding/json"

	v1alpha1 "github.com/cybercoder/tlscdn-controller/pkg/apis/v1alpha1"
	v1alpha1Types "github.com/cybercoder/tlscdn-controller/pkg/apis/v1alpha1/types"
	"github.com/cybercoder/tlscdn-controller/pkg/k8s"
	"github.com/cybercoder/tlscdn-controller/pkg/logger"
	"github.com/cybercoder/tlscdn-controller/pkg/redis"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
)

func OnAddWafRule(obj any) {
	u := obj.(*unstructured.Unstructured)
	var wafrule v1alpha1Types.WAFRule
	if err := runtime.DefaultUnstructuredConverter.FromUnstructured(u.Object, &wafrule); err != nil {
		logger.Errorf("Error converting unstructured to wafrule object: %v", err)
		return
	}

	k := k8s.CreateDynamicClient()

	gName := wafrule.Spec.CdnGateway
	_, err := k.Resource(v1alpha1.GatewayGVR).Namespace(wafrule.GetNamespace()).Get(context.TODO(), string(gName), metav1.GetOptions{})
	if err != nil {
		logger.Errorf("Gateway %s not found in namespace %s, orphaned wafrule: %s, err: %v", gName, wafrule.GetNamespace(), wafrule.GetName(), err)
		return
	}

	redisClient := redis.CreateClient()
	redisKey := "waf:" + wafrule.GetNamespace() + ":" + string(gName) + ":rules"

	redisRules := map[string]any{
		"rules": wafrule.Spec.Rules,
	}

	jsonData, err := json.Marshal(redisRules)
	if err != nil {
		logger.Errorf("failed to marshal rules: %v", err)
		return
	}
	err = redisClient.Set(context.Background(),
		redisKey, jsonData, 0).Err()
	if err != nil {
		logger.Errorf("Error on storing waf rule: %s to redis: %v", redisKey, err)
		return
	}

}

func OnUpdateWafRule(_, obj any) {
	u := obj.(*unstructured.Unstructured)
	var wafrule v1alpha1Types.WAFRule
	if err := runtime.DefaultUnstructuredConverter.FromUnstructured(u.Object, &wafrule); err != nil {
		logger.Errorf("Error converting unstructured to wafrule object: %v", err)
		return
	}

	redisClient := redis.CreateClient()
	redisKey := "waf:" + wafrule.GetNamespace() + ":" + string(wafrule.Spec.CdnGateway) + ":rules"
	luaDictCacheKey := wafrule.GetNamespace() + ":" + wafrule.Spec.CdnGateway

	redisRules := map[string]any{
		"rules": wafrule.Spec.Rules,
	}

	jsonData, err := json.Marshal(redisRules)
	if err != nil {
		logger.Errorf("failed to marshal rules: %v", err)
		return
	}

	err = redisClient.Publish(context.Background(), "invalidate_waf_cache", luaDictCacheKey).Err()
	if err != nil {
		logger.Errorf("Error on invalidating waf cache: %v", err)
		return
	}

	err = redisClient.Set(context.Background(),
		redisKey, jsonData, 0).Err()
	if err != nil {
		logger.Errorf("Error on storing waf rule: %s to redis: %v", redisKey, err)
		return
	}

}

func OnDeleteWafRule(obj any) {
	u := obj.(*unstructured.Unstructured)
	var wafrule v1alpha1Types.WAFRule
	if err := runtime.DefaultUnstructuredConverter.FromUnstructured(u.Object, &wafrule); err != nil {
		logger.Errorf("Error converting unstructured to wafrule object: %v", err)
		return
	}

	redisClient := redis.CreateClient()
	redisKey := "waf:" + wafrule.GetNamespace() + ":" + wafrule.Spec.CdnGateway + ":rules"
	luaDictCacheKey := wafrule.GetNamespace() + ":" + wafrule.Spec.CdnGateway

	err := redisClient.Publish(context.Background(), "invalidate_waf_cache", luaDictCacheKey).Err()
	if err != nil {
		logger.Errorf("Error on invalidating waf cache: %v", err)
		return
	}

	err = redisClient.Del(context.Background(), redisKey).Err()
	if err != nil {
		logger.Errorf("Error on deleting waf rule: %s from redis: %v", redisKey, err)
		return
	}

}
