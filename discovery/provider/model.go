// Implement types for each resource

package provider

import (
	appsv1 "k8s.io/api/apps/v1"
	autoscalingv1 "k8s.io/api/autoscaling/v1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	discoveryv1 "k8s.io/api/discovery/v1"
	networkingv1 "k8s.io/api/networking/v1"
	policyv1 "k8s.io/api/policy/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	storagev1 "k8s.io/api/storage/v1"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
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
	MetaObject  ObjectMeta
	ClusterRole ClusterRole
}

type KubernetesClusterRoleBindingDescription struct {
	MetaObject         ObjectMeta
	ClusterRoleBinding ClusterRoleBinding
}

type KubernetesConfigMapDescription struct {
	MetaObject ObjectMeta
	ConfigMap  corev1.ConfigMap
}

type KubernetesCronJobDescription struct {
	MetaObject ObjectMeta
	CronJob    batchv1.CronJob
}

type KubernetesCustomResourceDescription struct {
	MetaObject         ObjectMeta
	FullyQualifiedName string
	CustomResource     any
}

type KubernetesCustomResourceDefinitionDescription struct {
	MetaObject               ObjectMeta
	CustomResourceDefinition apiextensionsv1.CustomResourceDefinition
}

type KubernetesDaemonSetDescription struct {
	MetaObject ObjectMeta
	DaemonSet  appsv1.DaemonSet
}

type KubernetesDeploymentDescription struct {
	MetaObject ObjectMeta
	Deployment appsv1.Deployment
}

type KubernetesEndpointSliceDescription struct {
	MetaObject    ObjectMeta
	EndpointSlice discoveryv1.EndpointSlice
}

type KubernetesEndpointDescription struct {
	MetaObject ObjectMeta
	Endpoint   corev1.Endpoints
}

type KubernetesEventDescription struct {
	MetaObject ObjectMeta
	Event      corev1.Event
}

type KubernetesHorizontalPodAutoscalerDescription struct {
	MetaObject              ObjectMeta
	HorizontalPodAutoscaler autoscalingv1.HorizontalPodAutoscaler
}

type KubernetesIngressDescription struct {
	MetaObject ObjectMeta
	Ingress    networkingv1.Ingress
}

type KubernetesJobDescription struct {
	MetaObject ObjectMeta
	Job        batchv1.Job
}

type KubernetesLimitRangeDescription struct {
	MetaObject ObjectMeta
	LimitRange corev1.LimitRange
}

type KubernetesNamespaceDescription struct {
	MetaObject ObjectMeta
	Namespace  corev1.Namespace
}

type KubernetesNetworkPolicyDescription struct {
	MetaObject    ObjectMeta
	NetworkPolicy networkingv1.NetworkPolicy
}

type KubernetesNodeDescription struct {
	MetaObject ObjectMeta
	Node       corev1.Node
}

type KubernetesPersistentVolumeDescription struct {
	MetaObject ObjectMeta
	PV         corev1.PersistentVolume
}

type KubernetesPersistentVolumeClaimDescription struct {
	MetaObject ObjectMeta
	PVC        corev1.PersistentVolumeClaim
}

type KubernetesPodDescription struct {
	MetaObject ObjectMeta
	Pod        corev1.Pod
}

type KubernetesPodDisruptionBudgetDescription struct {
	MetaObject          ObjectMeta
	PodDisruptionBudget policyv1.PodDisruptionBudget
}

type KubernetesPodTemplateDescription struct {
	MetaObject  ObjectMeta
	PodTemplate corev1.PodTemplate
}

type KubernetesReplicaSetDescription struct {
	MetaObject ObjectMeta
	ReplicaSet appsv1.ReplicaSet
}

type KubernetesReplicationControllerDescription struct {
	MetaObject            ObjectMeta
	ReplicationController corev1.ReplicationController
}

type KubernetesResourceQuotaDescription struct {
	MetaObject    ObjectMeta
	ResourceQuota corev1.ResourceQuota
}

type KubernetesRoleDescription struct {
	MetaObject ObjectMeta
	Role       rbacv1.Role
}

type KubernetesRoleBindingDescription struct {
	MetaObject  ObjectMeta
	RoleBinding rbacv1.RoleBinding
}

type KubernetesSecretDescription struct {
	MetaObject ObjectMeta
	Secret     corev1.Secret
}
type KubernetesServiceDescription struct {
	MetaObject ObjectMeta
	Service    corev1.Service
}

type KubernetesServiceAccountDescription struct {
	MetaObject     ObjectMeta
	ServiceAccount corev1.ServiceAccount
}

type KubernetesStatefulSetDescription struct {
	MetaObject  ObjectMeta
	StatefulSet appsv1.StatefulSet
}

type KubernetesStorageClassDescription struct {
	MetaObject   ObjectMeta
	StorageClass storagev1.StorageClass
}
