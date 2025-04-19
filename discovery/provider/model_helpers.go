package provider

import (
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

type FieldsV1 struct {
	Raw []byte `json:"-" protobuf:"bytes,1,opt,name=Raw"`
}

type OwnerReference struct {
	APIVersion         string
	Kind               string
	Name               string
	UID                types.UID
	Controller         *bool
	BlockOwnerDeletion *bool
}

type ManagedFieldsEntry struct {
	Manager     string
	Operation   ManagedFieldsOperationType
	APIVersion  string
	Time        *Time
	FieldsType  string
	FieldsV1    *FieldsV1
	Subresource string
}

type TypeMeta struct {
	Kind       string
	APIVersion string
}

type ObjectMeta struct {
	Name                       string
	GenerateName               string
	Namespace                  string
	SelfLink                   string
	UID                        types.UID
	ResourceVersion            string
	Generation                 int64
	CreationTimestamp          Time
	DeletionTimestamp          *Time
	DeletionGracePeriodSeconds *int64
	Labels                     map[string]string
	Annotations                map[string]string
	OwnerReferences            []OwnerReference
	Finalizers                 []string
	ManagedFields              []ManagedFieldsEntry
}

type Subject struct {
	Kind      string
	APIGroup  string
	Name      string
	Namespace string
}

type RoleRef struct {
	APIGroup string
	Kind     string
	Name     string
}

type ClusterRoleBinding struct {
	TypeMeta
	ObjectMeta
	Subjects []Subject
	RoleRef  RoleRef
}
