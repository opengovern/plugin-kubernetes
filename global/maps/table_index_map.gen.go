package maps

import (
	"github.com/opengovern/og-describer-kubernetes/discovery/pkg/es"
)

var ResourceTypesToTables = map[string]string{
  "Kubernetes/Node": "kubernetes_node",
}

var ResourceTypeToDescription = map[string]interface{}{
  "Kubernetes/Node": opengovernance.KubernetesNode{},
}

var TablesToResourceTypes = map[string]string{
  "kubernetes_node": "Kubernetes/Node",
}
