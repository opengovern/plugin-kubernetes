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

	"Kubernetes/ReplicaSet": {
		IntegrationType:      constants.IntegrationName,
		ResourceName:         "Kubernetes/ReplicaSet",
		Tags:                 map[string][]string{
        },
		Labels:               map[string]string{
        },
		Annotations:          map[string]string{
        },
		ListDescriber:        provider.DescribeByIntegration(describers.KubernetesReplicaSet),
		GetDescriber:         nil,
	},

	"Kubernetes/DaemonSet": {
		IntegrationType:      constants.IntegrationName,
		ResourceName:         "Kubernetes/DaemonSet",
		Tags:                 map[string][]string{
        },
		Labels:               map[string]string{
        },
		Annotations:          map[string]string{
        },
		ListDescriber:        provider.DescribeByIntegration(describers.KubernetesDaemonSet),
		GetDescriber:         nil,
	},

	"Kubernetes/Endpoint": {
		IntegrationType:      constants.IntegrationName,
		ResourceName:         "Kubernetes/Endpoint",
		Tags:                 map[string][]string{
        },
		Labels:               map[string]string{
        },
		Annotations:          map[string]string{
        },
		ListDescriber:        provider.DescribeByIntegration(describers.KubernetesEndpoint),
		GetDescriber:         nil,
	},

	"Kubernetes/EndpointSlice": {
		IntegrationType:      constants.IntegrationName,
		ResourceName:         "Kubernetes/EndpointSlice",
		Tags:                 map[string][]string{
        },
		Labels:               map[string]string{
        },
		Annotations:          map[string]string{
        },
		ListDescriber:        provider.DescribeByIntegration(describers.KubernetesEndpointSlice),
		GetDescriber:         nil,
	},

	"Kubernetes/Event": {
		IntegrationType:      constants.IntegrationName,
		ResourceName:         "Kubernetes/Event",
		Tags:                 map[string][]string{
        },
		Labels:               map[string]string{
        },
		Annotations:          map[string]string{
        },
		ListDescriber:        provider.DescribeByIntegration(describers.KubernetesEvent),
		GetDescriber:         nil,
	},

	"Kubernetes/Job": {
		IntegrationType:      constants.IntegrationName,
		ResourceName:         "Kubernetes/Job",
		Tags:                 map[string][]string{
        },
		Labels:               map[string]string{
        },
		Annotations:          map[string]string{
        },
		ListDescriber:        provider.DescribeByIntegration(describers.KubernetesJob),
		GetDescriber:         nil,
	},

	"Kubernetes/CronJob": {
		IntegrationType:      constants.IntegrationName,
		ResourceName:         "Kubernetes/CronJob",
		Tags:                 map[string][]string{
        },
		Labels:               map[string]string{
        },
		Annotations:          map[string]string{
        },
		ListDescriber:        provider.DescribeByIntegration(describers.KubernetesCronJob),
		GetDescriber:         nil,
	},

	"Kubernetes/Ingress": {
		IntegrationType:      constants.IntegrationName,
		ResourceName:         "Kubernetes/Ingress",
		Tags:                 map[string][]string{
        },
		Labels:               map[string]string{
        },
		Annotations:          map[string]string{
        },
		ListDescriber:        provider.DescribeByIntegration(describers.KubernetesIngress),
		GetDescriber:         nil,
	},

	"Kubernetes/NetworkPolicy": {
		IntegrationType:      constants.IntegrationName,
		ResourceName:         "Kubernetes/NetworkPolicy",
		Tags:                 map[string][]string{
        },
		Labels:               map[string]string{
        },
		Annotations:          map[string]string{
        },
		ListDescriber:        provider.DescribeByIntegration(describers.KubernetesNetworkPolicy),
		GetDescriber:         nil,
	},

	"Kubernetes/Role": {
		IntegrationType:      constants.IntegrationName,
		ResourceName:         "Kubernetes/Role",
		Tags:                 map[string][]string{
        },
		Labels:               map[string]string{
        },
		Annotations:          map[string]string{
        },
		ListDescriber:        provider.DescribeByIntegration(describers.KubernetesRole),
		GetDescriber:         nil,
	},

	"Kubernetes/RoleBinding": {
		IntegrationType:      constants.IntegrationName,
		ResourceName:         "Kubernetes/RoleBinding",
		Tags:                 map[string][]string{
        },
		Labels:               map[string]string{
        },
		Annotations:          map[string]string{
        },
		ListDescriber:        provider.DescribeByIntegration(describers.KubernetesRoleBinding),
		GetDescriber:         nil,
	},

	"Kubernetes/Cluster": {
		IntegrationType:      constants.IntegrationName,
		ResourceName:         "Kubernetes/Cluster",
		Tags:                 map[string][]string{
        },
		Labels:               map[string]string{
        },
		Annotations:          map[string]string{
        },
		ListDescriber:        provider.DescribeByIntegration(describers.KubernetesCluster),
		GetDescriber:         nil,
	},

	"Kubernetes/ClusterRole": {
		IntegrationType:      constants.IntegrationName,
		ResourceName:         "Kubernetes/ClusterRole",
		Tags:                 map[string][]string{
        },
		Labels:               map[string]string{
        },
		Annotations:          map[string]string{
        },
		ListDescriber:        provider.DescribeByIntegration(describers.KubernetesClusterRole),
		GetDescriber:         nil,
	},

	"Kubernetes/ClusterRoleBinding": {
		IntegrationType:      constants.IntegrationName,
		ResourceName:         "Kubernetes/ClusterRoleBinding",
		Tags:                 map[string][]string{
        },
		Labels:               map[string]string{
        },
		Annotations:          map[string]string{
        },
		ListDescriber:        provider.DescribeByIntegration(describers.KubernetesClusterRoleBinding),
		GetDescriber:         nil,
	},

	"Kubernetes/PodDisruptionBudget": {
		IntegrationType:      constants.IntegrationName,
		ResourceName:         "Kubernetes/PodDisruptionBudget",
		Tags:                 map[string][]string{
        },
		Labels:               map[string]string{
        },
		Annotations:          map[string]string{
        },
		ListDescriber:        provider.DescribeByIntegration(describers.KubernetesPodDisruptionBudget),
		GetDescriber:         nil,
	},

	"Kubernetes/PodTemplate": {
		IntegrationType:      constants.IntegrationName,
		ResourceName:         "Kubernetes/PodTemplate",
		Tags:                 map[string][]string{
        },
		Labels:               map[string]string{
        },
		Annotations:          map[string]string{
        },
		ListDescriber:        provider.DescribeByIntegration(describers.KubernetesPodTemplate),
		GetDescriber:         nil,
	},

	"Kubernetes/HorizontalPodAutoscaler": {
		IntegrationType:      constants.IntegrationName,
		ResourceName:         "Kubernetes/HorizontalPodAutoscaler",
		Tags:                 map[string][]string{
        },
		Labels:               map[string]string{
        },
		Annotations:          map[string]string{
        },
		ListDescriber:        provider.DescribeByIntegration(describers.KubernetesHorizontalPodAutoscaler),
		GetDescriber:         nil,
	},

	"Kubernetes/CustomResourceDefinition": {
		IntegrationType:      constants.IntegrationName,
		ResourceName:         "Kubernetes/CustomResourceDefinition",
		Tags:                 map[string][]string{
        },
		Labels:               map[string]string{
        },
		Annotations:          map[string]string{
        },
		ListDescriber:        provider.DescribeByIntegration(describers.KubernetesCustomResourceDefinition),
		GetDescriber:         nil,
	},

	"Kubernetes/CustomResource": {
		IntegrationType:      constants.IntegrationName,
		ResourceName:         "Kubernetes/CustomResource",
		Tags:                 map[string][]string{
        },
		Labels:               map[string]string{
        },
		Annotations:          map[string]string{
        },
		ListDescriber:        provider.DescribeByIntegration(describers.KubernetesCustomResource),
		GetDescriber:         nil,
	},

	"Kubernetes/StorageClass": {
		IntegrationType:      constants.IntegrationName,
		ResourceName:         "Kubernetes/StorageClass",
		Tags:                 map[string][]string{
        },
		Labels:               map[string]string{
        },
		Annotations:          map[string]string{
        },
		ListDescriber:        provider.DescribeByIntegration(describers.KubernetesStorageClass),
		GetDescriber:         nil,
	},

	"Kubernetes/LimitRange": {
		IntegrationType:      constants.IntegrationName,
		ResourceName:         "Kubernetes/LimitRange",
		Tags:                 map[string][]string{
        },
		Labels:               map[string]string{
        },
		Annotations:          map[string]string{
        },
		ListDescriber:        provider.DescribeByIntegration(describers.KubernetesLimitRange),
		GetDescriber:         nil,
	},

	"Kubernetes/Namespace": {
		IntegrationType:      constants.IntegrationName,
		ResourceName:         "Kubernetes/Namespace",
		Tags:                 map[string][]string{
        },
		Labels:               map[string]string{
        },
		Annotations:          map[string]string{
        },
		ListDescriber:        provider.DescribeByIntegration(describers.KubernetesNamespace),
		GetDescriber:         nil,
	},

	"Kubernetes/ReplicationController": {
		IntegrationType:      constants.IntegrationName,
		ResourceName:         "Kubernetes/ReplicationController",
		Tags:                 map[string][]string{
        },
		Labels:               map[string]string{
        },
		Annotations:          map[string]string{
        },
		ListDescriber:        provider.DescribeByIntegration(describers.KubernetesReplicationController),
		GetDescriber:         nil,
	},

	"Kubernetes/RessourceQuota": {
		IntegrationType:      constants.IntegrationName,
		ResourceName:         "Kubernetes/RessourceQuota",
		Tags:                 map[string][]string{
        },
		Labels:               map[string]string{
        },
		Annotations:          map[string]string{
        },
		ListDescriber:        provider.DescribeByIntegration(describers.KubernetesResourceQuota),
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

	"Kubernetes/ReplicaSet": {
		Name:         "Kubernetes/ReplicaSet",
		IntegrationType:      constants.IntegrationName,
		Description:                 "",
		
	},

	"Kubernetes/DaemonSet": {
		Name:         "Kubernetes/DaemonSet",
		IntegrationType:      constants.IntegrationName,
		Description:                 "",
		
	},

	"Kubernetes/Endpoint": {
		Name:         "Kubernetes/Endpoint",
		IntegrationType:      constants.IntegrationName,
		Description:                 "",
		
	},

	"Kubernetes/EndpointSlice": {
		Name:         "Kubernetes/EndpointSlice",
		IntegrationType:      constants.IntegrationName,
		Description:                 "",
		
	},

	"Kubernetes/Event": {
		Name:         "Kubernetes/Event",
		IntegrationType:      constants.IntegrationName,
		Description:                 "",
		
	},

	"Kubernetes/Job": {
		Name:         "Kubernetes/Job",
		IntegrationType:      constants.IntegrationName,
		Description:                 "",
		
	},

	"Kubernetes/CronJob": {
		Name:         "Kubernetes/CronJob",
		IntegrationType:      constants.IntegrationName,
		Description:                 "",
		
	},

	"Kubernetes/Ingress": {
		Name:         "Kubernetes/Ingress",
		IntegrationType:      constants.IntegrationName,
		Description:                 "",
		
	},

	"Kubernetes/NetworkPolicy": {
		Name:         "Kubernetes/NetworkPolicy",
		IntegrationType:      constants.IntegrationName,
		Description:                 "",
		
	},

	"Kubernetes/Role": {
		Name:         "Kubernetes/Role",
		IntegrationType:      constants.IntegrationName,
		Description:                 "",
		
	},

	"Kubernetes/RoleBinding": {
		Name:         "Kubernetes/RoleBinding",
		IntegrationType:      constants.IntegrationName,
		Description:                 "",
		
	},

	"Kubernetes/Cluster": {
		Name:         "Kubernetes/Cluster",
		IntegrationType:      constants.IntegrationName,
		Description:                 "",
		
	},

	"Kubernetes/ClusterRole": {
		Name:         "Kubernetes/ClusterRole",
		IntegrationType:      constants.IntegrationName,
		Description:                 "",
		
	},

	"Kubernetes/ClusterRoleBinding": {
		Name:         "Kubernetes/ClusterRoleBinding",
		IntegrationType:      constants.IntegrationName,
		Description:                 "",
		
	},

	"Kubernetes/PodDisruptionBudget": {
		Name:         "Kubernetes/PodDisruptionBudget",
		IntegrationType:      constants.IntegrationName,
		Description:                 "",
		
	},

	"Kubernetes/PodTemplate": {
		Name:         "Kubernetes/PodTemplate",
		IntegrationType:      constants.IntegrationName,
		Description:                 "",
		
	},

	"Kubernetes/HorizontalPodAutoscaler": {
		Name:         "Kubernetes/HorizontalPodAutoscaler",
		IntegrationType:      constants.IntegrationName,
		Description:                 "",
		
	},

	"Kubernetes/CustomResourceDefinition": {
		Name:         "Kubernetes/CustomResourceDefinition",
		IntegrationType:      constants.IntegrationName,
		Description:                 "",
		
	},

	"Kubernetes/CustomResource": {
		Name:         "Kubernetes/CustomResource",
		IntegrationType:      constants.IntegrationName,
		Description:                 "",
		
	},

	"Kubernetes/StorageClass": {
		Name:         "Kubernetes/StorageClass",
		IntegrationType:      constants.IntegrationName,
		Description:                 "",
		
	},

	"Kubernetes/LimitRange": {
		Name:         "Kubernetes/LimitRange",
		IntegrationType:      constants.IntegrationName,
		Description:                 "",
		
	},

	"Kubernetes/Namespace": {
		Name:         "Kubernetes/Namespace",
		IntegrationType:      constants.IntegrationName,
		Description:                 "",
		
	},

	"Kubernetes/ReplicationController": {
		Name:         "Kubernetes/ReplicationController",
		IntegrationType:      constants.IntegrationName,
		Description:                 "",
		
	},

	"Kubernetes/RessourceQuota": {
		Name:         "Kubernetes/RessourceQuota",
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
  "Kubernetes/ReplicaSet",
  "Kubernetes/DaemonSet",
  "Kubernetes/Endpoint",
  "Kubernetes/EndpointSlice",
  "Kubernetes/Event",
  "Kubernetes/Job",
  "Kubernetes/CronJob",
  "Kubernetes/Ingress",
  "Kubernetes/NetworkPolicy",
  "Kubernetes/Role",
  "Kubernetes/RoleBinding",
  "Kubernetes/Cluster",
  "Kubernetes/ClusterRole",
  "Kubernetes/ClusterRoleBinding",
  "Kubernetes/PodDisruptionBudget",
  "Kubernetes/PodTemplate",
  "Kubernetes/HorizontalPodAutoscaler",
  "Kubernetes/CustomResourceDefinition",
  "Kubernetes/CustomResource",
  "Kubernetes/StorageClass",
  "Kubernetes/LimitRange",
  "Kubernetes/Namespace",
  "Kubernetes/ReplicationController",
  "Kubernetes/RessourceQuota",
}