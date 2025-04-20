package helpers

import (
	corev1 "k8s.io/api/core/v1"
)

// --- ConfigMap ---
type ConfigMap struct {
	TypeMeta   // Assumes TypeMeta in model_helpers.go
	ObjectMeta // Assumes ObjectMeta in model_helpers.go
	Immutable  *bool
	Data       map[string]string
}

// ConvertConfigMap creates a helper ConfigMap from a corev1 ConfigMap
func ConvertConfigMap(cm *corev1.ConfigMap) ConfigMap {
	return ConfigMap{
		TypeMeta:   ConvertTypeMeta(cm.TypeMeta),      // Assumes ConvertTypeMeta in model_helpers.go
		ObjectMeta: ConvertObjectMeta(&cm.ObjectMeta), // Assumes ConvertObjectMeta in model_helpers.go
		Immutable:  cm.Immutable,
		Data:       cm.Data,
	}
}
