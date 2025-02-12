// Implement types for each resource

package provider

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Metadata struct {
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

type KubernetesServiceDescription struct {
	MetaObject metav1.ObjectMeta
	Service    corev1.Service
}

type KubernetesSecretDescription struct {
	MetaObject metav1.ObjectMeta
	Secret     corev1.Secret
}
