package provider

import (
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"time"
)

type ManagedFieldsOperationType string

const (
	ManagedFieldsOperationApply  ManagedFieldsOperationType = "Apply"
	ManagedFieldsOperationUpdate ManagedFieldsOperationType = "Update"
)

type Time struct {
	time.Time
}

func ConvertTime(timestamp metav1.Time) time.Time {
	return timestamp.Time
}

func ConvertTimePtr(timestamp *metav1.Time) *time.Time {
	if timestamp == nil {
		return nil
	}
	return &timestamp.Time
}

type FieldsV1 struct {
	Raw []byte
}

func ConvertFieldsV1(raw metav1.FieldsV1) FieldsV1 {
	return FieldsV1{
		raw.Raw,
	}
}
func ConvertFieldsV1Ptr(raw *metav1.FieldsV1) *FieldsV1 {
	if raw == nil {
		return nil
	}
	return &FieldsV1{
		raw.Raw,
	}
}

type OwnerReference struct {
	APIVersion         string
	Kind               string
	Name               string
	UID                types.UID
	Controller         *bool
	BlockOwnerDeletion *bool
}

func ConvertOwnerReferences(ownerReferences []metav1.OwnerReference) []OwnerReference {
	ownerRefs := make([]OwnerReference, len(ownerReferences))
	for i, ownerRef := range ownerReferences {
		ownerRefs[i] = OwnerReference{
			APIVersion:         ownerRef.APIVersion,
			Kind:               ownerRef.Kind,
			Name:               ownerRef.Name,
			UID:                ownerRef.UID,
			Controller:         ownerRef.Controller,
			BlockOwnerDeletion: ownerRef.BlockOwnerDeletion,
		}
	}
	return ownerRefs
}

type ManagedFieldsEntry struct {
	Manager     string
	Operation   ManagedFieldsOperationType
	APIVersion  string
	Time        *time.Time
	FieldsType  string
	FieldsV1    *FieldsV1
	Subresource string
}

func ConvertManagedFieldsEntries(managedFieldsEntry []metav1.ManagedFieldsEntry) []ManagedFieldsEntry {
	managedFieldsEntries := make([]ManagedFieldsEntry, len(managedFieldsEntry))
	for i, entry := range managedFieldsEntry {
		managedFieldsEntries[i] = ManagedFieldsEntry{
			Manager:     entry.Manager,
			Operation:   ManagedFieldsOperationType(entry.Operation),
			APIVersion:  entry.APIVersion,
			Time:        ConvertTimePtr(entry.Time),
			FieldsType:  entry.FieldsType,
			FieldsV1:    ConvertFieldsV1Ptr(entry.FieldsV1),
			Subresource: entry.Subresource,
		}
	}
	return managedFieldsEntries
}

type TypeMeta struct {
	Kind       string
	APIVersion string
}

func ConvertTypeMeta(typeMeta metav1.TypeMeta) TypeMeta {
	return TypeMeta{
		Kind:       typeMeta.Kind,
		APIVersion: typeMeta.APIVersion,
	}
}

type ObjectMeta struct {
	Name                       string
	GenerateName               string
	Namespace                  string
	SelfLink                   string
	UID                        types.UID
	ResourceVersion            string
	Generation                 int64
	CreationTimestamp          time.Time
	DeletionTimestamp          *time.Time
	DeletionGracePeriodSeconds *int64
	Labels                     map[string]string
	Annotations                map[string]string
	OwnerReferences            []OwnerReference
	Finalizers                 []string
	ManagedFields              []ManagedFieldsEntry
}

func ConvertObjectMeta(obj *metav1.ObjectMeta) ObjectMeta {
	return ObjectMeta{
		Name:                       obj.Name,
		GenerateName:               obj.GenerateName,
		Namespace:                  obj.Namespace,
		SelfLink:                   obj.GetSelfLink(),
		UID:                        obj.UID,
		ResourceVersion:            obj.ResourceVersion,
		Generation:                 obj.Generation,
		CreationTimestamp:          ConvertTime(obj.CreationTimestamp),
		DeletionTimestamp:          ConvertTimePtr(obj.DeletionTimestamp),
		DeletionGracePeriodSeconds: obj.DeletionGracePeriodSeconds,
		Labels:                     obj.Labels,
		Annotations:                obj.Annotations,
		OwnerReferences:            ConvertOwnerReferences(obj.OwnerReferences),
		Finalizers:                 obj.Finalizers,
		ManagedFields:              ConvertManagedFieldsEntries(obj.ManagedFields),
	}
}

type Subject struct {
	Kind      string
	APIGroup  string
	Name      string
	Namespace string
}

func ConvertSubject(srcSubjects []rbacv1.Subject) []Subject {
	subjects := make([]Subject, len(srcSubjects))
	for i, srcSubject := range srcSubjects {
		subjects[i] = Subject{
			Kind:      srcSubject.Kind,
			APIGroup:  srcSubject.APIGroup,
			Name:      srcSubject.Name,
			Namespace: srcSubject.Namespace,
		}
	}
	return subjects
}

type RoleRef struct {
	APIGroup string
	Kind     string
	Name     string
}

func ConvertRoleRef(roleRef rbacv1.RoleRef) RoleRef {
	return RoleRef{
		APIGroup: roleRef.APIGroup,
		Kind:     roleRef.Kind,
		Name:     roleRef.Name,
	}
}

type ClusterRoleBinding struct {
	TypeMeta
	ObjectMeta
	Subjects []Subject
	RoleRef  RoleRef
}
