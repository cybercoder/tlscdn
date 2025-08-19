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

func OnAddWafRule(obj interface{}) {
	u := obj.(*unstructured.Unstructured)
	var wafrule v1alpha1Types.WAFRule
	if err := runtime.DefaultUnstructuredConverter.FromUnstructured(u.Object, &wafrule); err != nil {
		logger.Errorf("Error converting unstructured to wafrule object: %v", err)
		return
	}

	k := k8s.CreateDynamicClient()

	gName := wafrule.Spec.Gateway.Name
	_, err := k.Resource(v1alpha1.GatewayGVR).Namespace(wafrule.GetNamespace()).Get(context.TODO(), string(gName), metav1.GetOptions{})
	if err != nil {
		logger.Errorf("Gateway %s not found in namespace %s, orphaned wafrule: %s, err: %v", gName, wafrule.GetNamespace(), wafrule.GetName(), err)
		// return
	}

	redisClient := redis.CreateClient()
	redisKey := "WAF_RULES:" + string(gName)

	jsonData, err := json.Marshal(wafrule.Spec)
	if err != nil {
		logger.Errorf("can't convert rules %s/%s to json: %v", wafrule.GetNamespace(), wafrule.GetName(), err)
	}
	err = redisClient.Set(context.Background(),
		redisKey, jsonData, 0).Err()
	if err != nil {
		logger.Errorf("Error on storing waf rule: %s to redis: %v", redisKey, err)
		return
	}

}
