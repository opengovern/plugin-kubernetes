package main

import (
	"github.com/opengovern/og-describer-kubernetes/global/constants"
)

// Config represents the JSON input configuration
type Config struct {
}

func IntegrationHealthcheck(creds constants.IntegrationCredentials, cfg Config) (bool, error) {
	return DoHealthcheck(creds.KubeConfig)
}
