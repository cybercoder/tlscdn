package k8s

import (
	"os"
	"path/filepath"

	"github.com/cybercoder/tlscdn-controller/pkg/logger"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"

	cmclient "github.com/cert-manager/cert-manager/pkg/client/clientset/versioned"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

var client *kubernetes.Clientset
var dynamicClient *dynamic.DynamicClient
var cmClient *cmclient.Clientset

func CreateClient() *kubernetes.Clientset {
	// singleton
	if client != nil {
		return client
	}

	config, err := createConfig()
	if err != nil {
		logger.Fatalf("Error creating config: %v", err)
	}

	client, err = kubernetes.NewForConfig(config)
	if err != nil {
		logger.Fatalf("Error creating dynamicClient: %v", err)
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
		logger.Fatalf("Error creating config: %v", err)
	}

	dynamicClient, err = dynamic.NewForConfig(config)
	if err != nil {
		logger.Fatalf("Error creating dynamicClient: %v", err)
	}
	return dynamicClient
}

func CreateCertManagerClient() *cmclient.Clientset {
	if cmClient != nil {
		return cmClient
	}
	config, err := createConfig()
	if err != nil {
		logger.Fatalf("Error creating config: %v", err)
	}

	cmClient, err = cmclient.NewForConfig(config)
	if err != nil {
		logger.Fatalf("Error creating cert-manager client: %v", err)
	}
	return cmClient
}

func createConfig() (*rest.Config, error) {
	configFile := filepath.Join(homedir.HomeDir(), ".kube", "config")
	_, err := os.Stat(configFile)
	if err != nil {
		return rest.InClusterConfig()
	}
	return clientcmd.BuildConfigFromFlags("", configFile)
}
