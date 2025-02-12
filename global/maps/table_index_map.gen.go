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
  "Kubernetes/ConfigMap": "kubernetes_config_map",
  "Kubernetes/ServiceAccount": "kubernetes_service_account",
  "Kubernetes/StatefulSet": "kubernetes_stateful_set",
  "Kubernetes/Deployment": "kubernetes_deployment",
}

var ResourceTypeToDescription = map[string]interface{}{
  "Kubernetes/Node": opengovernance.KubernetesNode{},
  "Kubernetes/PersistentVolume": opengovernance.KubernetesPersistentVolume{},
  "Kubernetes/PersistentVolumeClaim": opengovernance.KubernetesPersistentVolumeClaim{},
  "Kubernetes/Pod": opengovernance.KubernetesPod{},
  "Kubernetes/Secret": opengovernance.KubernetesSecret{},
  "Kubernetes/Service": opengovernance.KubernetesService{},
  "Kubernetes/ConfigMap": opengovernance.KubernetesConfigMap{},
  "Kubernetes/ServiceAccount": opengovernance.KubernetesServiceAccount{},
  "Kubernetes/StatefulSet": opengovernance.KubernetesStatefulSet{},
  "Kubernetes/Deployment": opengovernance.KubernetesDeployment{},
}

var TablesToResourceTypes = map[string]string{
  "kubernetes_node": "Kubernetes/Node",
  "kubernetes_persistent_volume": "Kubernetes/PersistentVolume",
  "kubernetes_persistent_volume_claim": "Kubernetes/PersistentVolumeClaim",
  "kubernetes_pod": "Kubernetes/Pod",
  "kubernetes_secret": "Kubernetes/Secret",
  "kubernetes_service": "Kubernetes/Service",
  "kubernetes_config_map": "Kubernetes/ConfigMap",
  "kubernetes_service_account": "Kubernetes/ServiceAccount",
  "kubernetes_stateful_set": "Kubernetes/StatefulSet",
  "kubernetes_deployment": "Kubernetes/Deployment",
}
