package main

import (
	"github.com/opengovern/og-describer-kubernetes/global/constants"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// Config represents the JSON input configuration
type Config struct {
}

func IntegrationHealthcheck(creds constants.IntegrationCredentials, cfg Config) (bool, error) {
	config, err := clientcmd.RESTConfigFromKubeConfig([]byte(creds.KubeConfig))
	if err != nil {
		return false, err
	}
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		return false, err
	}

	_, err = clientSet.ServerVersion()
	if err != nil {
		return false, err
	}

	return true, nil
}
