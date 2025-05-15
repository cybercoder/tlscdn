package k8s

import (
	"log"
	"os"
	"path/filepath"

	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"

	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

var client *kubernetes.Clientset
var dynamicClient *dynamic.DynamicClient

func CreateClient() *kubernetes.Clientset {
	// singleton
	if client != nil {
		return client
	}

	config, err := createConfig()
	if err != nil {
		log.Fatalf("Error creating config: %v", err)
	}

	client, err = kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("Error creating dynamicClient: %v", err)
	}
	return client
}

func CreateDynamicClient() *dynamic.DynamicClient {
	// singleton
	if dynamicClient != nil {
		return dynamicClient
	}

	config, err := createConfig()
	if err != nil {
		log.Fatalf("Error creating config: %v", err)
	}

	dynamicClient, err = dynamic.NewForConfig(config)
	if err != nil {
		log.Fatalf("Error creating dynamicClient: %v", err)
	}
	return dynamicClient
}

func createConfig() (*rest.Config, error) {
	configFile := filepath.Join(homedir.HomeDir(), ".kube", "config")
	_, err := os.Stat(configFile)
	if err != nil {
		return rest.InClusterConfig()
	}
	return clientcmd.BuildConfigFromFlags("", configFile)
}
