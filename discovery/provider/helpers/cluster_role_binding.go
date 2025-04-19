package helpers

import (
	rbacv1 "k8s.io/api/rbac/v1"
)

// --- ClusterRoleBinding ---
type ClusterRoleBinding struct {
	TypeMeta
	ObjectMeta
	Subjects []Subject // Uses Subject from this file
	RoleRef  RoleRef   // Uses RoleRef from this file
}

// ConvertClusterRoleBinding creates a helper ClusterRoleBinding from a rbacv1 ClusterRoleBinding
func ConvertClusterRoleBinding(crb *rbacv1.ClusterRoleBinding) ClusterRoleBinding {
	return ClusterRoleBinding{
		TypeMeta:   ConvertTypeMeta(crb.TypeMeta),      // Assumes ConvertTypeMeta in model_helpers.go
		ObjectMeta: ConvertObjectMeta(&crb.ObjectMeta), // Assumes ConvertObjectMeta in model_helpers.go
		Subjects:   ConvertSubjects(crb.Subjects),      // Uses ConvertSubjects from this file
		RoleRef:    ConvertRoleRef(crb.RoleRef),        // Uses ConvertRoleRef from this file
	}
}

// --- Subject ---
type Subject struct {
	Kind      string
	APIGroup  string
	Name      string
	Namespace string
}

// ConvertSubject converts a single rbacv1.Subject to a helper Subject
func ConvertSubject(srcSubject rbacv1.Subject) Subject {
	return Subject{
		Kind:      srcSubject.Kind,
		APIGroup:  srcSubject.APIGroup,
		Name:      srcSubject.Name,
		Namespace: srcSubject.Namespace,
	}
}

// ConvertSubjects converts a slice of rbacv1.Subject to a slice of helper Subject
func ConvertSubjects(srcSubjects []rbacv1.Subject) []Subject {
	if srcSubjects == nil {
		return nil
	}
	subjects := make([]Subject, len(srcSubjects))
	for i, srcSubject := range srcSubjects {
		subjects[i] = ConvertSubject(srcSubject)
	}
	return subjects
}

// --- RoleRef ---
type RoleRef struct {
	APIGroup string
	Kind     string
	Name     string
}

// ConvertRoleRef converts a rbacv1.RoleRef to a helper RoleRef
func ConvertRoleRef(roleRef rbacv1.RoleRef) RoleRef {
	return RoleRef{
		APIGroup: roleRef.APIGroup,
		Kind:     roleRef.Kind,
		Name:     roleRef.Name,
	}
}
