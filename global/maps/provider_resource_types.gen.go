package maps

import (
	"github.com/opengovern/og-describer-kubernetes/discovery/describers"
	model "github.com/opengovern/og-describer-kubernetes/discovery/pkg/models"
	"github.com/opengovern/og-describer-kubernetes/discovery/provider"
	"github.com/opengovern/og-describer-kubernetes/platform/constants"
	"github.com/opengovern/og-util/pkg/integration/interfaces"
)

var ResourceTypes = map[string]model.ResourceType{

	"Kubernetes/Node": {
		IntegrationType: constants.IntegrationName,
		ResourceName:    "Kubernetes/Node",
		Tags:            map[string][]string{},
		Labels:          map[string]string{},
		Annotations:     map[string]string{},
		ListDescriber:   provider.DescribeByIntegration(describers.KubernetesNode),
		GetDescriber:    nil,
	},
}

var ResourceTypeConfigs = map[string]*interfaces.ResourceTypeConfiguration{

	"Kubernetes/Node": {
		Name:            "Kubernetes/Node",
		IntegrationType: constants.IntegrationName,
		Description:     "",
	},
}

var ResourceTypesList = []string{
	"Kubernetes/Node",
}
