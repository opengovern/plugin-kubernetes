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
