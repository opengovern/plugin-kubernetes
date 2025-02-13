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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Metadata struct {
}

type KubernetesClusterRoleDescription struct {
	MetaObject  metav1.ObjectMeta
	ClusterRole rbacv1.ClusterRole
}

type KubernetesClusterRoleBindingDescription struct {
	MetaObject         metav1.ObjectMeta
	ClusterRoleBinding rbacv1.ClusterRoleBinding
}

type KubernetesConfigMapDescription struct {
	MetaObject metav1.ObjectMeta
	ConfigMap  corev1.ConfigMap
}

type KubernetesCronJobDescription struct {
	MetaObject metav1.ObjectMeta
	CronJob    batchv1.CronJob
}

type KubernetesCustomResourceDescription struct {
	MetaObject         metav1.ObjectMeta
	FullyQualifiedName string
	CustomResource     any
}

type KubernetesCustomResourceDefinitionDescription struct {
	MetaObject               metav1.ObjectMeta
	CustomResourceDefinition apiextensionsv1.CustomResourceDefinition
}

type KubernetesDaemonSetDescription struct {
	MetaObject metav1.ObjectMeta
	DaemonSet  appsv1.DaemonSet
}

type KubernetesDeploymentDescription struct {
	MetaObject metav1.ObjectMeta
	Deployment appsv1.Deployment
}

type KubernetesEndpointSliceDescription struct {
	MetaObject    metav1.ObjectMeta
	EndpointSlice discoveryv1.EndpointSlice
}

type KubernetesEndpointDescription struct {
	MetaObject metav1.ObjectMeta
	Endpoint   corev1.Endpoints
}

type KubernetesEventDescription struct {
	MetaObject metav1.ObjectMeta
	Event      corev1.Event
}

type KubernetesHorizontalPodAutoscalerDescription struct {
	MetaObject              metav1.ObjectMeta
	HorizontalPodAutoscaler autoscalingv1.HorizontalPodAutoscaler
}

type KubernetesIngressDescription struct {
	MetaObject metav1.ObjectMeta
	Ingress    networkingv1.Ingress
}

type KubernetesJobDescription struct {
	MetaObject metav1.ObjectMeta
	Job        batchv1.Job
}

type KubernetesLimitRangeDescription struct {
	MetaObject metav1.ObjectMeta
	LimitRange corev1.LimitRange
}

type KubernetesNamespaceDescription struct {
	MetaObject metav1.ObjectMeta
	Namespace  corev1.Namespace
}

type KubernetesNetworkPolicyDescription struct {
	MetaObject    metav1.ObjectMeta
	NetworkPolicy networkingv1.NetworkPolicy
}

type KubernetesNodeDescription struct {
	MetaObject metav1.ObjectMeta
	Node       corev1.Node
}

type KubernetesPersistentVolumeDescription struct {
	MetaObject metav1.ObjectMeta
	PV         corev1.PersistentVolume
}

type KubernetesPersistentVolumeClaimDescription struct {
	MetaObject metav1.ObjectMeta
	PVC        corev1.PersistentVolumeClaim
}

type KubernetesPodDescription struct {
	MetaObject metav1.ObjectMeta
	Pod        corev1.Pod
}

type KubernetesPodDisruptionBudgetDescription struct {
	MetaObject          metav1.ObjectMeta
	PodDisruptionBudget policyv1.PodDisruptionBudget
}

type KubernetesPodTemplateDescription struct {
	MetaObject  metav1.ObjectMeta
	PodTemplate corev1.PodTemplate
}

type KubernetesReplicaSetDescription struct {
	MetaObject metav1.ObjectMeta
	ReplicaSet appsv1.ReplicaSet
}

type KubernetesReplicationControllerDescription struct {
	MetaObject            metav1.ObjectMeta
	ReplicationController corev1.ReplicationController
}

type KubernetesResourceQuotaDescription struct {
	MetaObject    metav1.ObjectMeta
	ResourceQuota corev1.ResourceQuota
}

type KubernetesRoleDescription struct {
	MetaObject metav1.ObjectMeta
	Role       rbacv1.Role
}

type KubernetesRoleBindingDescription struct {
	MetaObject  metav1.ObjectMeta
	RoleBinding rbacv1.RoleBinding
}

type KubernetesSecretDescription struct {
	MetaObject metav1.ObjectMeta
	Secret     corev1.Secret
}
type KubernetesServiceDescription struct {
	MetaObject metav1.ObjectMeta
	Service    corev1.Service
}

type KubernetesServiceAccountDescription struct {
	MetaObject     metav1.ObjectMeta
	ServiceAccount corev1.ServiceAccount
}

type KubernetesStatefulSetDescription struct {
	MetaObject  metav1.ObjectMeta
	StatefulSet appsv1.StatefulSet
}

type KubernetesStorageClassDescription struct {
	MetaObject   metav1.ObjectMeta
	StorageClass storagev1.StorageClass
}
