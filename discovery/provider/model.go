// Implement types for each resource

package provider

import (
	"github.com/opengovern/og-describer-kubernetes/discovery/provider/helpers"
	batchv1 "k8s.io/api/batch/v1"
)

type Metadata struct {
}

type KubernetesClusterDescription struct {
	AuthMethod            string `json:"auth_method"`
	ContextName           string `json:"context_name"`
	Endpoint              string `json:"endpoint"`
	TLSServerVerification bool   `json:"tls_server_verification"`
	ServerVersion         string `json:"server_version"`
}

type KubernetesClusterRoleDescription struct {
	MetaObject  helpers.ObjectMeta
	ClusterRole helpers.ClusterRole
}

type KubernetesClusterRoleBindingDescription struct {
	MetaObject         helpers.ObjectMeta
	ClusterRoleBinding helpers.ClusterRoleBinding
}

type KubernetesConfigMapDescription struct {
	MetaObject helpers.ObjectMeta
	ConfigMap  helpers.ConfigMap
}

type KubernetesCronJobDescription struct {
	MetaObject helpers.ObjectMeta
	CronJob    helpers.CronJob
}

type KubernetesCustomResourceDescription struct {
	MetaObject         helpers.ObjectMeta
	FullyQualifiedName string
	CustomResource     any
}

type KubernetesCustomResourceDefinitionDescription struct {
	MetaObject               helpers.ObjectMeta
	CustomResourceDefinition helpers.CustomResourceDefinition
}

type KubernetesDaemonSetDescription struct {
	MetaObject helpers.ObjectMeta
	DaemonSet  helpers.DaemonSet
}

type KubernetesDeploymentDescription struct {
	MetaObject helpers.ObjectMeta
	Deployment helpers.Deployment
}

type KubernetesEndpointSliceDescription struct {
	MetaObject    helpers.ObjectMeta
	EndpointSlice helpers.EndpointSlice
}

type KubernetesEndpointDescription struct {
	MetaObject helpers.ObjectMeta
	Endpoint   helpers.Endpoints
}

type KubernetesEventDescription struct {
	MetaObject helpers.ObjectMeta
	Event      helpers.Event
}

type KubernetesHorizontalPodAutoscalerDescription struct {
	MetaObject              helpers.ObjectMeta
	HorizontalPodAutoscaler helpers.HorizontalPodAutoscaler
}

type KubernetesIngressDescription struct {
	MetaObject helpers.ObjectMeta
	Ingress    helpers.Ingress
}

type KubernetesJobDescription struct {
	MetaObject helpers.ObjectMeta
	Job        batchv1.Job
}

type KubernetesLimitRangeDescription struct {
	MetaObject helpers.ObjectMeta
	LimitRange helpers.LimitRange
}

type KubernetesNamespaceDescription struct {
	MetaObject helpers.ObjectMeta
	Namespace  helpers.Namespace
}

type KubernetesNetworkPolicyDescription struct {
	MetaObject    helpers.ObjectMeta
	NetworkPolicy helpers.NetworkPolicy
}

type KubernetesNodeDescription struct {
	MetaObject helpers.ObjectMeta
	Node       helpers.Node
}

type KubernetesPersistentVolumeDescription struct {
	MetaObject helpers.ObjectMeta
	PV         helpers.PersistentVolume
}

type KubernetesPersistentVolumeClaimDescription struct {
	MetaObject helpers.ObjectMeta
	PVC        helpers.PersistentVolumeClaim
}

type KubernetesPodDescription struct {
	MetaObject helpers.ObjectMeta
	Pod        helpers.Pod
}

type KubernetesPodDisruptionBudgetDescription struct {
	MetaObject          helpers.ObjectMeta
	PodDisruptionBudget helpers.PodDisruptionBudget
}

type KubernetesPodTemplateDescription struct {
	MetaObject  helpers.ObjectMeta
	PodTemplate helpers.PodTemplate
}

type KubernetesReplicaSetDescription struct {
	MetaObject helpers.ObjectMeta
	ReplicaSet helpers.ReplicaSet
}

type KubernetesReplicationControllerDescription struct {
	MetaObject            helpers.ObjectMeta
	ReplicationController helpers.ReplicationController
}

type KubernetesResourceQuotaDescription struct {
	MetaObject    helpers.ObjectMeta
	ResourceQuota helpers.ResourceQuota
}

type KubernetesRoleDescription struct {
	MetaObject helpers.ObjectMeta
	Role       helpers.Role
}

type KubernetesRoleBindingDescription struct {
	MetaObject  helpers.ObjectMeta
	RoleBinding helpers.RoleBinding
}

type KubernetesSecretDescription struct {
	MetaObject helpers.ObjectMeta
	Secret     helpers.Secret
}
type KubernetesServiceDescription struct {
	MetaObject helpers.ObjectMeta
	Service    helpers.Service
}

type KubernetesServiceAccountDescription struct {
	MetaObject     helpers.ObjectMeta
	ServiceAccount helpers.ServiceAccount
}

type KubernetesStatefulSetDescription struct {
	MetaObject  helpers.ObjectMeta
	StatefulSet helpers.StatefulSet
}

type KubernetesStorageClassDescription struct {
	MetaObject   helpers.ObjectMeta
	StorageClass helpers.StorageClass
}
