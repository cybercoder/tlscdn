package events

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func OnAddGateway(obj interface{}) {
	u := obj.(*unstructured.Unstructured)

}
