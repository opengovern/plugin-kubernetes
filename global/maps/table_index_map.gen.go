package maps

import (
	"github.com/opengovern/og-describer-kubernetes/discovery/pkg/es"
)

var ResourceTypesToTables = map[string]string{
  "Kubernetes/Node": "kubernetes_node",
  "Kubernetes/PersistentVolume": "kubernetes_persistent_volume",
  "Kubernetes/PersistentVolumeClaim": "kubernetes_persistent_volume_claim",
  "Kubernetes/Pod": "kubernetes_pod",
  "Kubernetes/Secret": "kubernetes_secret",
  "Kubernetes/Service": "kubernetes_service",
}

var ResourceTypeToDescription = map[string]interface{}{
  "Kubernetes/Node": opengovernance.KubernetesNode{},
  "Kubernetes/PersistentVolume": opengovernance.KubernetesPersistentVolume{},
  "Kubernetes/PersistentVolumeClaim": opengovernance.KubernetesPersistentVolumeClaim{},
  "Kubernetes/Pod": opengovernance.KubernetesPod{},
  "Kubernetes/Secret": opengovernance.KubernetesSecret{},
  "Kubernetes/Service": opengovernance.KubernetesService{},
}

var TablesToResourceTypes = map[string]string{
  "kubernetes_node": "Kubernetes/Node",
  "kubernetes_persistent_volume": "Kubernetes/PersistentVolume",
  "kubernetes_persistent_volume_claim": "Kubernetes/PersistentVolumeClaim",
  "kubernetes_pod": "Kubernetes/Pod",
  "kubernetes_secret": "Kubernetes/Secret",
  "kubernetes_service": "Kubernetes/Service",
}
