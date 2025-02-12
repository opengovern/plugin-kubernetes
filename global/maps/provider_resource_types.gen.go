package maps
import (
	"github.com/opengovern/og-describer-kubernetes/discovery/describers"
	"github.com/opengovern/og-describer-kubernetes/discovery/provider"
	"github.com/opengovern/og-describer-kubernetes/platform/constants"
	"github.com/opengovern/og-util/pkg/integration/interfaces"
	model "github.com/opengovern/og-describer-kubernetes/discovery/pkg/models"
)
var ResourceTypes = map[string]model.ResourceType{

	"Kubernetes/Node": {
		IntegrationType:      constants.IntegrationName,
		ResourceName:         "Kubernetes/Node",
		Tags:                 map[string][]string{
        },
		Labels:               map[string]string{
        },
		Annotations:          map[string]string{
        },
		ListDescriber:        provider.DescribeByIntegration(describers.KubernetesNode),
		GetDescriber:         nil,
	},

	"Kubernetes/PersistentVolume": {
		IntegrationType:      constants.IntegrationName,
		ResourceName:         "Kubernetes/PersistentVolume",
		Tags:                 map[string][]string{
        },
		Labels:               map[string]string{
        },
		Annotations:          map[string]string{
        },
		ListDescriber:        provider.DescribeByIntegration(describers.KubernetesPersistentVolume),
		GetDescriber:         nil,
	},

	"Kubernetes/PersistentVolumeClaim": {
		IntegrationType:      constants.IntegrationName,
		ResourceName:         "Kubernetes/PersistentVolumeClaim",
		Tags:                 map[string][]string{
        },
		Labels:               map[string]string{
        },
		Annotations:          map[string]string{
        },
		ListDescriber:        provider.DescribeByIntegration(describers.KubernetesPersistentVolumeClaim),
		GetDescriber:         nil,
	},

	"Kubernetes/Pod": {
		IntegrationType:      constants.IntegrationName,
		ResourceName:         "Kubernetes/Pod",
		Tags:                 map[string][]string{
        },
		Labels:               map[string]string{
        },
		Annotations:          map[string]string{
        },
		ListDescriber:        provider.DescribeByIntegration(describers.KubernetesPod),
		GetDescriber:         nil,
	},

	"Kubernetes/Secret": {
		IntegrationType:      constants.IntegrationName,
		ResourceName:         "Kubernetes/Secret",
		Tags:                 map[string][]string{
        },
		Labels:               map[string]string{
        },
		Annotations:          map[string]string{
        },
		ListDescriber:        provider.DescribeByIntegration(describers.KubernetesSecret),
		GetDescriber:         nil,
	},

	"Kubernetes/Service": {
		IntegrationType:      constants.IntegrationName,
		ResourceName:         "Kubernetes/Service",
		Tags:                 map[string][]string{
        },
		Labels:               map[string]string{
        },
		Annotations:          map[string]string{
        },
		ListDescriber:        provider.DescribeByIntegration(describers.KubernetesService),
		GetDescriber:         nil,
	},

	"Kubernetes/ConfigMap": {
		IntegrationType:      constants.IntegrationName,
		ResourceName:         "Kubernetes/ConfigMap",
		Tags:                 map[string][]string{
        },
		Labels:               map[string]string{
        },
		Annotations:          map[string]string{
        },
		ListDescriber:        provider.DescribeByIntegration(describers.KubernetesConfigMap),
		GetDescriber:         nil,
	},

	"Kubernetes/ServiceAccount": {
		IntegrationType:      constants.IntegrationName,
		ResourceName:         "Kubernetes/ServiceAccount",
		Tags:                 map[string][]string{
        },
		Labels:               map[string]string{
        },
		Annotations:          map[string]string{
        },
		ListDescriber:        provider.DescribeByIntegration(describers.KubernetesServiceAccount),
		GetDescriber:         nil,
	},

	"Kubernetes/StatefulSet": {
		IntegrationType:      constants.IntegrationName,
		ResourceName:         "Kubernetes/StatefulSet",
		Tags:                 map[string][]string{
        },
		Labels:               map[string]string{
        },
		Annotations:          map[string]string{
        },
		ListDescriber:        provider.DescribeByIntegration(describers.KubernetesStatefulSet),
		GetDescriber:         nil,
	},

	"Kubernetes/Deployment": {
		IntegrationType:      constants.IntegrationName,
		ResourceName:         "Kubernetes/Deployment",
		Tags:                 map[string][]string{
        },
		Labels:               map[string]string{
        },
		Annotations:          map[string]string{
        },
		ListDescriber:        provider.DescribeByIntegration(describers.KubernetesDeployment),
		GetDescriber:         nil,
	},
}


var ResourceTypeConfigs = map[string]*interfaces.ResourceTypeConfiguration{

	"Kubernetes/Node": {
		Name:         "Kubernetes/Node",
		IntegrationType:      constants.IntegrationName,
		Description:                 "",
		
	},

	"Kubernetes/PersistentVolume": {
		Name:         "Kubernetes/PersistentVolume",
		IntegrationType:      constants.IntegrationName,
		Description:                 "",
		
	},

	"Kubernetes/PersistentVolumeClaim": {
		Name:         "Kubernetes/PersistentVolumeClaim",
		IntegrationType:      constants.IntegrationName,
		Description:                 "",
		
	},

	"Kubernetes/Pod": {
		Name:         "Kubernetes/Pod",
		IntegrationType:      constants.IntegrationName,
		Description:                 "",
		
	},

	"Kubernetes/Secret": {
		Name:         "Kubernetes/Secret",
		IntegrationType:      constants.IntegrationName,
		Description:                 "",
		
	},

	"Kubernetes/Service": {
		Name:         "Kubernetes/Service",
		IntegrationType:      constants.IntegrationName,
		Description:                 "",
		
	},

	"Kubernetes/ConfigMap": {
		Name:         "Kubernetes/ConfigMap",
		IntegrationType:      constants.IntegrationName,
		Description:                 "",
		
	},

	"Kubernetes/ServiceAccount": {
		Name:         "Kubernetes/ServiceAccount",
		IntegrationType:      constants.IntegrationName,
		Description:                 "",
		
	},

	"Kubernetes/StatefulSet": {
		Name:         "Kubernetes/StatefulSet",
		IntegrationType:      constants.IntegrationName,
		Description:                 "",
		
	},

	"Kubernetes/Deployment": {
		Name:         "Kubernetes/Deployment",
		IntegrationType:      constants.IntegrationName,
		Description:                 "",
		
	},
}


var ResourceTypesList = []string{
  "Kubernetes/Node",
  "Kubernetes/PersistentVolume",
  "Kubernetes/PersistentVolumeClaim",
  "Kubernetes/Pod",
  "Kubernetes/Secret",
  "Kubernetes/Service",
  "Kubernetes/ConfigMap",
  "Kubernetes/ServiceAccount",
  "Kubernetes/StatefulSet",
  "Kubernetes/Deployment",
}