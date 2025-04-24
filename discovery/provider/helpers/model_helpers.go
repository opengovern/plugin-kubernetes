package helpers

import (
	"time"

	appsv1 "k8s.io/api/apps/v1"
	autoscalingv1 "k8s.io/api/autoscaling/v1"
	corev1 "k8s.io/api/core/v1"
	discoveryv1 "k8s.io/api/discovery/v1"
	networkingv1 "k8s.io/api/networking/v1"
	policyv1 "k8s.io/api/policy/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	storagev1 "k8s.io/api/storage/v1"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// --- Base Meta Types ---
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
	t := timestamp.Time // Get the time.Time value
	return &t           // Return its address
}

type FieldsV1 struct {
	Raw []byte
}

func ConvertFieldsV1(raw metav1.FieldsV1) FieldsV1 {
	return FieldsV1{
		Raw: raw.Raw,
	}
}
func ConvertFieldsV1Ptr(raw *metav1.FieldsV1) *FieldsV1 {
	if raw == nil {
		return nil
	}
	return &FieldsV1{
		Raw: raw.Raw,
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
	if ownerReferences == nil {
		return nil
	}
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
	Time        *time.Time // Uses helpers time via ConvertTimePtr
	FieldsType  string
	FieldsV1    *FieldsV1
	Subresource string
}

func ConvertManagedFieldsEntries(managedFieldsEntries []metav1.ManagedFieldsEntry) []ManagedFieldsEntry {
	if managedFieldsEntries == nil {
		return nil
	}
	result := make([]ManagedFieldsEntry, len(managedFieldsEntries))
	for i, entry := range managedFieldsEntries {
		result[i] = ManagedFieldsEntry{
			Manager:     entry.Manager,
			Operation:   ManagedFieldsOperationType(entry.Operation),
			APIVersion:  entry.APIVersion,
			Time:        ConvertTimePtr(entry.Time),
			FieldsType:  entry.FieldsType,
			FieldsV1:    ConvertFieldsV1Ptr(entry.FieldsV1),
			Subresource: entry.Subresource,
		}
	}
	return result
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
	CreationTimestamp          time.Time  // Uses helpers time via ConvertTime
	DeletionTimestamp          *time.Time // Uses helpers time via ConvertTimePtr
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
		SelfLink:                   obj.GetSelfLink(), // This method might not exist on the helper type if created later
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

// --- LabelSelector ---
type LabelSelectorRequirement struct {
	Key      string
	Operator string // metav1.LabelSelectorOperator
	Values   []string
}

func ConvertLabelSelectorOperator(op metav1.LabelSelectorOperator) string {
	return string(op)
}

func ConvertLabelSelectorRequirements(labelSelectorRequirements []metav1.LabelSelectorRequirement) []LabelSelectorRequirement {
	if labelSelectorRequirements == nil {
		return nil
	}
	requirements := make([]LabelSelectorRequirement, len(labelSelectorRequirements))
	for i, labelSelectorRequirement := range labelSelectorRequirements {
		requirements[i] = LabelSelectorRequirement{
			Key:      labelSelectorRequirement.Key,
			Operator: ConvertLabelSelectorOperator(labelSelectorRequirement.Operator),
			Values:   labelSelectorRequirement.Values,
		}
	}
	return requirements
}

type LabelSelector struct {
	MatchLabels      map[string]string
	MatchExpressions []LabelSelectorRequirement
}

func ConvertLabelSelector(ls *metav1.LabelSelector) *LabelSelector {
	if ls == nil {
		return nil
	}
	return &LabelSelector{
		MatchLabels:      ls.MatchLabels,
		MatchExpressions: ConvertLabelSelectorRequirements(ls.MatchExpressions),
	}
}

// --- ObjectReference (Single Correct Definition) ---
type ObjectReference struct {
	Kind            string
	Namespace       string
	Name            string
	UID             types.UID
	APIVersion      string
	ResourceVersion string
	FieldPath       string
}

// ConvertObjectReference converts a corev1.ObjectReference to helpers.ObjectReference
func ConvertObjectReference(ref corev1.ObjectReference) ObjectReference {
	return ObjectReference{
		Kind:            ref.Kind,
		Namespace:       ref.Namespace,
		Name:            ref.Name,
		UID:             ref.UID,
		APIVersion:      ref.APIVersion,
		ResourceVersion: ref.ResourceVersion,
		FieldPath:       ref.FieldPath,
	}
}

// --- ServiceAccount ---
type ServiceAccount struct {
	TypeMeta
	ObjectMeta
	Secrets                      []ObjectReference      // Uses helpers.ObjectReference
	ImagePullSecrets             []LocalObjectReference // Assumes LocalObjectReference defined elsewhere (e.g., volume.go)
	AutomountServiceAccountToken *bool
}

func ConvertServiceAccount(sa *corev1.ServiceAccount) ServiceAccount {
	return ServiceAccount{
		TypeMeta:                     ConvertTypeMeta(sa.TypeMeta),
		ObjectMeta:                   ConvertObjectMeta(&sa.ObjectMeta),
		Secrets:                      ConvertObjectReferences(sa.Secrets),
		ImagePullSecrets:             ConvertLocalObjectReferences(sa.ImagePullSecrets), // Assumes defined elsewhere (e.g., pod.go)
		AutomountServiceAccountToken: sa.AutomountServiceAccountToken,
	}
}

// --- Secret ---
type Secret struct {
	TypeMeta
	ObjectMeta
	Immutable  *bool
	Data       map[string][]byte
	StringData map[string]string
	Type       string // corev1.SecretType
}

func ConvertSecretType(st corev1.SecretType) string {
	return string(st)
}

func ConvertSecret(s *corev1.Secret) Secret {
	return Secret{
		TypeMeta:   ConvertTypeMeta(s.TypeMeta),
		ObjectMeta: ConvertObjectMeta(&s.ObjectMeta),
		Immutable:  s.Immutable,
		Data:       s.Data,
		StringData: s.StringData,
		Type:       ConvertSecretType(s.Type),
	}
}

// --- Namespace ---
type Namespace struct {
	TypeMeta
	ObjectMeta
	Spec   NamespaceSpec
	Status NamespaceStatus
}

type NamespaceSpec struct {
	Finalizers []string // corev1.FinalizerName
}

type NamespaceStatus struct {
	Phase      string // corev1.NamespacePhase
	Conditions []NamespaceCondition
}

type NamespaceCondition struct {
	Type               string     // corev1.NamespaceConditionType
	Status             string     // corev1.ConditionStatus
	LastTransitionTime *time.Time // Uses helpers time via ConvertTimePtr
	Reason             string
	Message            string
}

func ConvertFinalizerName(fn corev1.FinalizerName) string {
	return string(fn)
}

func ConvertFinalizers(fns []corev1.FinalizerName) []string {
	if fns == nil {
		return nil
	}
	result := make([]string, len(fns))
	for i, f := range fns {
		result[i] = ConvertFinalizerName(f)
	}
	return result
}

func ConvertNamespaceSpec(spec corev1.NamespaceSpec) NamespaceSpec {
	return NamespaceSpec{
		Finalizers: ConvertFinalizers(spec.Finalizers),
	}
}

func ConvertNamespacePhase(p corev1.NamespacePhase) string {
	return string(p)
}

func ConvertNamespaceConditionType(t corev1.NamespaceConditionType) string {
	return string(t)
}

func ConvertConditionStatus(s corev1.ConditionStatus) string {
	return string(s)
}

func ConvertNamespaceCondition(c corev1.NamespaceCondition) NamespaceCondition {
	return NamespaceCondition{
		Type:               ConvertNamespaceConditionType(c.Type),
		Status:             ConvertConditionStatus(c.Status),
		LastTransitionTime: ConvertTimePtr(&c.LastTransitionTime),
		Reason:             c.Reason,
		Message:            c.Message,
	}
}

func ConvertNamespaceConditions(conds []corev1.NamespaceCondition) []NamespaceCondition {
	if conds == nil {
		return nil
	}
	result := make([]NamespaceCondition, len(conds))
	for i, c := range conds {
		result[i] = ConvertNamespaceCondition(c)
	}
	return result
}

func ConvertNamespaceStatus(status corev1.NamespaceStatus) NamespaceStatus {
	return NamespaceStatus{
		Phase:      ConvertNamespacePhase(status.Phase),
		Conditions: ConvertNamespaceConditions(status.Conditions),
	}
}

func ConvertNamespace(ns *corev1.Namespace) Namespace {
	return Namespace{
		TypeMeta:   ConvertTypeMeta(ns.TypeMeta),
		ObjectMeta: ConvertObjectMeta(&ns.ObjectMeta),
		Spec:       ConvertNamespaceSpec(ns.Spec),
		Status:     ConvertNamespaceStatus(ns.Status),
	}
}

// --- Role (rbacv1) ---
type Role struct {
	TypeMeta
	ObjectMeta
	Rules []PolicyRule // Assumes PolicyRule defined elsewhere (e.g., cluster_role.go)
}

func ConvertRole(r *rbacv1.Role) Role {
	return Role{
		TypeMeta:   ConvertTypeMeta(r.TypeMeta),
		ObjectMeta: ConvertObjectMeta(&r.ObjectMeta),
		Rules:      ConvertPolicyRules(r.Rules), // Assumes ConvertPolicyRules defined elsewhere
	}
}

// --- RoleBinding (rbacv1) ---
type RoleBinding struct {
	TypeMeta
	ObjectMeta
	Subjects []Subject // Assumes Subject defined elsewhere (e.g., cluster_role_binding.go)
	RoleRef  RoleRef   // Assumes RoleRef defined elsewhere
}

func ConvertRoleBinding(rb *rbacv1.RoleBinding) RoleBinding {
	return RoleBinding{
		TypeMeta:   ConvertTypeMeta(rb.TypeMeta),
		ObjectMeta: ConvertObjectMeta(&rb.ObjectMeta),
		Subjects:   ConvertSubjects(rb.Subjects), // Assumes ConvertSubjects defined elsewhere
		RoleRef:    ConvertRoleRef(rb.RoleRef),   // Assumes ConvertRoleRef defined elsewhere
	}
}

// --- Adding appsv1 types ---

// --- Deployment ---
type Deployment struct {
	TypeMeta
	ObjectMeta
	Spec   DeploymentSpec
	Status DeploymentStatus
}

type DeploymentSpec struct {
	Replicas                *int32
	Selector                *LabelSelector  // Use helpers.LabelSelector
	Template                PodTemplateSpec // Use helpers.PodTemplateSpec
	Strategy                DeploymentStrategy
	MinReadySeconds         int32
	RevisionHistoryLimit    *int32
	Paused                  bool
	ProgressDeadlineSeconds *int32
}

type DeploymentStrategy struct {
	Type          string // appsv1.DeploymentStrategyType
	RollingUpdate *RollingUpdateDeployment
}

type RollingUpdateDeployment struct {
	MaxUnavailable *intstr.IntOrString // Keep as IntOrString for simplicity or map
	MaxSurge       *intstr.IntOrString // Keep as IntOrString for simplicity or map
}

type DeploymentStatus struct {
	ObservedGeneration  int64
	Replicas            int32
	UpdatedReplicas     int32
	ReadyReplicas       int32
	AvailableReplicas   int32
	UnavailableReplicas int32
	Conditions          []DeploymentCondition
	CollisionCount      *int32
}

type DeploymentCondition struct {
	Type               string    // appsv1.DeploymentConditionType
	Status             string    // corev1.ConditionStatus
	LastUpdateTime     time.Time // Use helpers.Time via ConvertTime
	LastTransitionTime time.Time // Use helpers.Time via ConvertTime
	Reason             string
	Message            string
}

func ConvertDeploymentStrategyType(dst appsv1.DeploymentStrategyType) string {
	return string(dst)
}

// ConvertIntOrString defined in pod.go? If not, define or handle here.
// For now, assume IntOrString can be used directly or needs mapping.
func ConvertRollingUpdateDeployment(rud *appsv1.RollingUpdateDeployment) *RollingUpdateDeployment {
	if rud == nil {
		return nil
	}
	return &RollingUpdateDeployment{
		MaxUnavailable: rud.MaxUnavailable,
		MaxSurge:       rud.MaxSurge,
	}
}

func ConvertDeploymentStrategy(strategy appsv1.DeploymentStrategy) DeploymentStrategy {
	return DeploymentStrategy{
		Type:          ConvertDeploymentStrategyType(strategy.Type),
		RollingUpdate: ConvertRollingUpdateDeployment(strategy.RollingUpdate),
	}
}

func ConvertDeploymentSpec(spec appsv1.DeploymentSpec) DeploymentSpec {
	return DeploymentSpec{
		Replicas:                spec.Replicas,
		Selector:                ConvertLabelSelector(spec.Selector),
		Template:                ConvertPodTemplateSpec(spec.Template),
		Strategy:                ConvertDeploymentStrategy(spec.Strategy),
		MinReadySeconds:         spec.MinReadySeconds,
		RevisionHistoryLimit:    spec.RevisionHistoryLimit,
		Paused:                  spec.Paused,
		ProgressDeadlineSeconds: spec.ProgressDeadlineSeconds,
	}
}

func ConvertDeploymentConditionType(dct appsv1.DeploymentConditionType) string {
	return string(dct)
}

func ConvertDeploymentCondition(c appsv1.DeploymentCondition) DeploymentCondition {
	return DeploymentCondition{
		Type:               ConvertDeploymentConditionType(c.Type),
		Status:             ConvertConditionStatus(c.Status),
		LastUpdateTime:     ConvertTime(c.LastUpdateTime),
		LastTransitionTime: ConvertTime(c.LastTransitionTime),
		Reason:             c.Reason,
		Message:            c.Message,
	}
}

func ConvertDeploymentConditions(conditions []appsv1.DeploymentCondition) []DeploymentCondition {
	if conditions == nil {
		return nil
	}
	result := make([]DeploymentCondition, len(conditions))
	for i, c := range conditions {
		result[i] = ConvertDeploymentCondition(c)
	}
	return result
}

func ConvertDeploymentStatus(status appsv1.DeploymentStatus) DeploymentStatus {
	return DeploymentStatus{
		ObservedGeneration:  status.ObservedGeneration,
		Replicas:            status.Replicas,
		UpdatedReplicas:     status.UpdatedReplicas,
		ReadyReplicas:       status.ReadyReplicas,
		AvailableReplicas:   status.AvailableReplicas,
		UnavailableReplicas: status.UnavailableReplicas,
		Conditions:          ConvertDeploymentConditions(status.Conditions),
		CollisionCount:      status.CollisionCount,
	}
}

func ConvertDeployment(dep *appsv1.Deployment) Deployment {
	return Deployment{
		TypeMeta:   ConvertTypeMeta(dep.TypeMeta),
		ObjectMeta: ConvertObjectMeta(&dep.ObjectMeta),
		Spec:       ConvertDeploymentSpec(dep.Spec),
		Status:     ConvertDeploymentStatus(dep.Status),
	}
}

// --- DaemonSet ---
type DaemonSet struct {
	TypeMeta
	ObjectMeta
	Spec   DaemonSetSpec
	Status DaemonSetStatus
}

type DaemonSetSpec struct {
	Selector             *LabelSelector
	Template             PodTemplateSpec
	UpdateStrategy       DaemonSetUpdateStrategy
	MinReadySeconds      int32
	RevisionHistoryLimit *int32
}

type DaemonSetUpdateStrategy struct {
	Type          string // appsv1.DaemonSetUpdateStrategyType
	RollingUpdate *RollingUpdateDaemonSet
}

type RollingUpdateDaemonSet struct {
	MaxUnavailable *intstr.IntOrString // Keep IntOrString
	MaxSurge       *intstr.IntOrString // Keep IntOrString (added in later k8s versions)
}

type DaemonSetStatus struct {
	CurrentNumberScheduled int32
	NumberMisscheduled     int32
	DesiredNumberScheduled int32
	NumberReady            int32
	ObservedGeneration     int64
	UpdatedNumberScheduled int32
	NumberAvailable        int32
	NumberUnavailable      int32
	CollisionCount         *int32
	Conditions             []DaemonSetCondition
}

type DaemonSetCondition struct {
	Type               string    // appsv1.DaemonSetConditionType
	Status             string    // corev1.ConditionStatus
	LastTransitionTime time.Time // Use helpers.Time via ConvertTime
	Reason             string
	Message            string
}

func ConvertDaemonSetUpdateStrategyType(dsust appsv1.DaemonSetUpdateStrategyType) string {
	return string(dsust)
}

func ConvertRollingUpdateDaemonSet(ruds *appsv1.RollingUpdateDaemonSet) *RollingUpdateDaemonSet {
	if ruds == nil {
		return nil
	}
	return &RollingUpdateDaemonSet{
		MaxUnavailable: ruds.MaxUnavailable,
		MaxSurge:       ruds.MaxSurge,
	}
}

func ConvertDaemonSetUpdateStrategy(strategy appsv1.DaemonSetUpdateStrategy) DaemonSetUpdateStrategy {
	return DaemonSetUpdateStrategy{
		Type:          ConvertDaemonSetUpdateStrategyType(strategy.Type),
		RollingUpdate: ConvertRollingUpdateDaemonSet(strategy.RollingUpdate),
	}
}

func ConvertDaemonSetSpec(spec appsv1.DaemonSetSpec) DaemonSetSpec {
	return DaemonSetSpec{
		Selector:             ConvertLabelSelector(spec.Selector),
		Template:             ConvertPodTemplateSpec(spec.Template),
		UpdateStrategy:       ConvertDaemonSetUpdateStrategy(spec.UpdateStrategy),
		MinReadySeconds:      spec.MinReadySeconds,
		RevisionHistoryLimit: spec.RevisionHistoryLimit,
	}
}

func ConvertDaemonSetConditionType(dsct appsv1.DaemonSetConditionType) string {
	return string(dsct)
}

func ConvertDaemonSetCondition(c appsv1.DaemonSetCondition) DaemonSetCondition {
	return DaemonSetCondition{
		Type:               ConvertDaemonSetConditionType(c.Type),
		Status:             ConvertConditionStatus(c.Status),
		LastTransitionTime: ConvertTime(c.LastTransitionTime),
		Reason:             c.Reason,
		Message:            c.Message,
	}
}

func ConvertDaemonSetConditions(conditions []appsv1.DaemonSetCondition) []DaemonSetCondition {
	if conditions == nil {
		return nil
	}
	result := make([]DaemonSetCondition, len(conditions))
	for i, c := range conditions {
		result[i] = ConvertDaemonSetCondition(c)
	}
	return result
}

func ConvertDaemonSetStatus(status appsv1.DaemonSetStatus) DaemonSetStatus {
	return DaemonSetStatus{
		CurrentNumberScheduled: status.CurrentNumberScheduled,
		NumberMisscheduled:     status.NumberMisscheduled,
		DesiredNumberScheduled: status.DesiredNumberScheduled,
		NumberReady:            status.NumberReady,
		ObservedGeneration:     status.ObservedGeneration,
		UpdatedNumberScheduled: status.UpdatedNumberScheduled,
		NumberAvailable:        status.NumberAvailable,
		NumberUnavailable:      status.NumberUnavailable,
		CollisionCount:         status.CollisionCount,
		Conditions:             ConvertDaemonSetConditions(status.Conditions),
	}
}

func ConvertDaemonSet(ds *appsv1.DaemonSet) DaemonSet {
	return DaemonSet{
		TypeMeta:   ConvertTypeMeta(ds.TypeMeta),
		ObjectMeta: ConvertObjectMeta(&ds.ObjectMeta),
		Spec:       ConvertDaemonSetSpec(ds.Spec),
		Status:     ConvertDaemonSetStatus(ds.Status),
	}
}

// --- ReplicaSet ---
type ReplicaSet struct {
	TypeMeta
	ObjectMeta
	Spec   ReplicaSetSpec
	Status ReplicaSetStatus
}

type ReplicaSetSpec struct {
	Replicas        *int32
	MinReadySeconds int32
	Selector        *LabelSelector
	Template        PodTemplateSpec
}

type ReplicaSetStatus struct {
	Replicas             int32
	FullyLabeledReplicas int32
	ReadyReplicas        int32
	AvailableReplicas    int32
	ObservedGeneration   int64
	Conditions           []ReplicaSetCondition
}

type ReplicaSetCondition struct {
	Type               string    // appsv1.ReplicaSetConditionType
	Status             string    // corev1.ConditionStatus
	LastTransitionTime time.Time // Use helpers.Time via ConvertTime
	Reason             string
	Message            string
}

func ConvertReplicaSetSpec(spec appsv1.ReplicaSetSpec) ReplicaSetSpec {
	return ReplicaSetSpec{
		Replicas:        spec.Replicas,
		MinReadySeconds: spec.MinReadySeconds,
		Selector:        ConvertLabelSelector(spec.Selector),
		Template:        ConvertPodTemplateSpec(spec.Template),
	}
}

func ConvertReplicaSetConditionType(rsct appsv1.ReplicaSetConditionType) string {
	return string(rsct)
}

func ConvertReplicaSetCondition(c appsv1.ReplicaSetCondition) ReplicaSetCondition {
	return ReplicaSetCondition{
		Type:               ConvertReplicaSetConditionType(c.Type),
		Status:             ConvertConditionStatus(c.Status),
		LastTransitionTime: ConvertTime(c.LastTransitionTime),
		Reason:             c.Reason,
		Message:            c.Message,
	}
}

func ConvertReplicaSetConditions(conditions []appsv1.ReplicaSetCondition) []ReplicaSetCondition {
	if conditions == nil {
		return nil
	}
	result := make([]ReplicaSetCondition, len(conditions))
	for i, c := range conditions {
		result[i] = ConvertReplicaSetCondition(c)
	}
	return result
}

func ConvertReplicaSetStatus(status appsv1.ReplicaSetStatus) ReplicaSetStatus {
	return ReplicaSetStatus{
		Replicas:             status.Replicas,
		FullyLabeledReplicas: status.FullyLabeledReplicas,
		ReadyReplicas:        status.ReadyReplicas,
		AvailableReplicas:    status.AvailableReplicas,
		ObservedGeneration:   status.ObservedGeneration,
		Conditions:           ConvertReplicaSetConditions(status.Conditions),
	}
}

func ConvertReplicaSet(rs *appsv1.ReplicaSet) ReplicaSet {
	return ReplicaSet{
		TypeMeta:   ConvertTypeMeta(rs.TypeMeta),
		ObjectMeta: ConvertObjectMeta(&rs.ObjectMeta),
		Spec:       ConvertReplicaSetSpec(rs.Spec),
		Status:     ConvertReplicaSetStatus(rs.Status),
	}
}

// --- StatefulSet ---
type StatefulSet struct {
	TypeMeta
	ObjectMeta
	Spec   StatefulSetSpec
	Status StatefulSetStatus
}

type StatefulSetSpec struct {
	Replicas                             *int32
	Selector                             *LabelSelector
	Template                             PodTemplateSpec
	VolumeClaimTemplates                 []PersistentVolumeClaim // Assumes PersistentVolumeClaim defined elsewhere (it's not standard PVC)
	ServiceName                          string
	PodManagementPolicy                  string // appsv1.PodManagementPolicyType
	UpdateStrategy                       StatefulSetUpdateStrategy
	RevisionHistoryLimit                 *int32
	MinReadySeconds                      int32                                            // Added later?
	PersistentVolumeClaimRetentionPolicy *StatefulSetPersistentVolumeClaimRetentionPolicy // Added later
	Ordinals                             *StatefulSetOrdinals                             // Added later
}

// NOTE: StatefulSet needs PersistentVolumeClaim helper type. Add if not present.
// type PersistentVolumeClaim struct { ... } // Need definition and conversion

type StatefulSetUpdateStrategy struct {
	Type          string // appsv1.StatefulSetUpdateStrategyType
	RollingUpdate *RollingUpdateStatefulSetStrategy
}

type RollingUpdateStatefulSetStrategy struct {
	Partition      *int32
	MaxUnavailable *intstr.IntOrString // Added later
}

type StatefulSetStatus struct {
	ObservedGeneration int64
	Replicas           int32
	ReadyReplicas      int32
	CurrentReplicas    int32
	UpdatedReplicas    int32
	CurrentRevision    string
	UpdateRevision     string
	CollisionCount     *int32
	Conditions         []StatefulSetCondition
	AvailableReplicas  int32 // Added later
}

type StatefulSetCondition struct {
	Type               string    // appsv1.StatefulSetConditionType
	Status             string    // corev1.ConditionStatus
	LastTransitionTime time.Time // Use helpers.Time via ConvertTime
	Reason             string
	Message            string
}

type StatefulSetPersistentVolumeClaimRetentionPolicy struct {
	WhenDeleted string // appsv1.PersistentVolumeClaimRetentionPolicyType
	WhenScaled  string // appsv1.PersistentVolumeClaimRetentionPolicyType
}

type StatefulSetOrdinals struct {
	Start int32
}

func ConvertPodManagementPolicyType(pmpt appsv1.PodManagementPolicyType) string {
	return string(pmpt)
}

func ConvertStatefulSetUpdateStrategyType(ssust appsv1.StatefulSetUpdateStrategyType) string {
	return string(ssust)
}

func ConvertRollingUpdateStatefulSetStrategy(russ *appsv1.RollingUpdateStatefulSetStrategy) *RollingUpdateStatefulSetStrategy {
	if russ == nil {
		return nil
	}
	return &RollingUpdateStatefulSetStrategy{
		Partition:      russ.Partition,
		MaxUnavailable: russ.MaxUnavailable,
	}
}

func ConvertStatefulSetUpdateStrategy(strategy appsv1.StatefulSetUpdateStrategy) StatefulSetUpdateStrategy {
	return StatefulSetUpdateStrategy{
		Type:          ConvertStatefulSetUpdateStrategyType(strategy.Type),
		RollingUpdate: ConvertRollingUpdateStatefulSetStrategy(strategy.RollingUpdate),
	}
}

func ConvertPersistentVolumeClaimRetentionPolicyType(pvcrpt appsv1.PersistentVolumeClaimRetentionPolicyType) string {
	return string(pvcrpt)
}

func ConvertStatefulSetPersistentVolumeClaimRetentionPolicy(policy *appsv1.StatefulSetPersistentVolumeClaimRetentionPolicy) *StatefulSetPersistentVolumeClaimRetentionPolicy {
	if policy == nil {
		return nil
	}
	return &StatefulSetPersistentVolumeClaimRetentionPolicy{
		WhenDeleted: ConvertPersistentVolumeClaimRetentionPolicyType(policy.WhenDeleted),
		WhenScaled:  ConvertPersistentVolumeClaimRetentionPolicyType(policy.WhenScaled),
	}
}

func ConvertStatefulSetOrdinals(ordinals *appsv1.StatefulSetOrdinals) *StatefulSetOrdinals {
	if ordinals == nil {
		return nil
	}
	return &StatefulSetOrdinals{
		Start: ordinals.Start,
	}
}

// Placeholder for converting []corev1.PersistentVolumeClaim - needs definition
func ConvertVolumeClaimTemplates(templates []corev1.PersistentVolumeClaim) []PersistentVolumeClaim {
	if templates == nil {
		return nil
	}
	result := make([]PersistentVolumeClaim, len(templates))
	for i := range templates { // Iterate by index to pass pointer
		result[i] = ConvertPersistentVolumeClaim(&templates[i])
	}
	return result
}

func ConvertStatefulSetSpec(spec appsv1.StatefulSetSpec) StatefulSetSpec {
	return StatefulSetSpec{
		Replicas:                             spec.Replicas,
		Selector:                             ConvertLabelSelector(spec.Selector),
		Template:                             ConvertPodTemplateSpec(spec.Template),
		VolumeClaimTemplates:                 ConvertVolumeClaimTemplates(spec.VolumeClaimTemplates), // Use updated function
		ServiceName:                          spec.ServiceName,
		PodManagementPolicy:                  ConvertPodManagementPolicyType(spec.PodManagementPolicy),
		UpdateStrategy:                       ConvertStatefulSetUpdateStrategy(spec.UpdateStrategy),
		RevisionHistoryLimit:                 spec.RevisionHistoryLimit,
		MinReadySeconds:                      spec.MinReadySeconds,
		PersistentVolumeClaimRetentionPolicy: ConvertStatefulSetPersistentVolumeClaimRetentionPolicy(spec.PersistentVolumeClaimRetentionPolicy),
		Ordinals:                             ConvertStatefulSetOrdinals(spec.Ordinals),
	}
}

func ConvertStatefulSetConditionType(ssct appsv1.StatefulSetConditionType) string {
	return string(ssct)
}

func ConvertStatefulSetCondition(c appsv1.StatefulSetCondition) StatefulSetCondition {
	return StatefulSetCondition{
		Type:               ConvertStatefulSetConditionType(c.Type),
		Status:             ConvertConditionStatus(c.Status),
		LastTransitionTime: ConvertTime(c.LastTransitionTime),
		Reason:             c.Reason,
		Message:            c.Message,
	}
}

func ConvertStatefulSetConditions(conditions []appsv1.StatefulSetCondition) []StatefulSetCondition {
	if conditions == nil {
		return nil
	}
	result := make([]StatefulSetCondition, len(conditions))
	for i, c := range conditions {
		result[i] = ConvertStatefulSetCondition(c)
	}
	return result
}

func ConvertStatefulSetStatus(status appsv1.StatefulSetStatus) StatefulSetStatus {
	return StatefulSetStatus{
		ObservedGeneration: status.ObservedGeneration,
		Replicas:           status.Replicas,
		ReadyReplicas:      status.ReadyReplicas,
		CurrentReplicas:    status.CurrentReplicas,
		UpdatedReplicas:    status.UpdatedReplicas,
		CurrentRevision:    status.CurrentRevision,
		UpdateRevision:     status.UpdateRevision,
		CollisionCount:     status.CollisionCount,
		Conditions:         ConvertStatefulSetConditions(status.Conditions),
		AvailableReplicas:  status.AvailableReplicas,
	}
}

func ConvertStatefulSet(ss *appsv1.StatefulSet) StatefulSet {
	return StatefulSet{
		TypeMeta:   ConvertTypeMeta(ss.TypeMeta),
		ObjectMeta: ConvertObjectMeta(&ss.ObjectMeta),
		Spec:       ConvertStatefulSetSpec(ss.Spec),
		Status:     ConvertStatefulSetStatus(ss.Status),
	}
}

// --- PersistentVolumeClaim (Needed for StatefulSet) ---
type PersistentVolumeClaim struct {
	TypeMeta
	ObjectMeta
	Spec   PersistentVolumeClaimSpec // Use helpers.PersistentVolumeClaimSpec
	Status PersistentVolumeClaimStatus
}

// Need PersistentVolumeClaimSpec definition if not already present
// Assuming it's defined in volume.go or needs to be added here
type PersistentVolumeClaimSpec struct {
	AccessModes      []string
	Selector         *LabelSelector             // Use helpers.LabelSelector
	Resources        VolumeResourceRequirements // Use helpers.VolumeResourceRequirements
	VolumeName       string
	StorageClassName *string
	VolumeMode       *string
	DataSource       *TypedLocalObjectReference // Use helpers.TypedLocalObjectReference
	DataSourceRef    *TypedObjectReference      // Use helpers.TypedObjectReference
}

// Need PersistentVolumeClaimStatus definition
type PersistentVolumeClaimStatus struct {
	Phase                            string                       // corev1.PersistentVolumeClaimPhase
	AccessModes                      []string                     // corev1.PersistentVolumeAccessMode
	Capacity                         map[string]resource.Quantity // corev1.ResourceList
	Conditions                       []PersistentVolumeClaimCondition
	AllocatedResources               map[string]resource.Quantity // corev1.ResourceList
	AllocatedResourceStatuses        map[string]string            // corev1.ClaimResourceStatus - New in 1.27?
	CurrentVolumeAttributesClassName *string
	ModifyVolumeStatus               *struct {
		TargetVolumeAttributesClassName string
		Status                          string
	}
}

func ConvertPersistentVolumeClaimStatus(status corev1.PersistentVolumeClaimStatus) PersistentVolumeClaimStatus {
	var accessModes []string
	if status.AccessModes != nil {
		accessModes = make([]string, len(status.AccessModes))
		for i, a := range status.AccessModes {
			accessModes[i] = string(a)
		}
	}
	capacity := make(map[string]resource.Quantity)
	if status.Capacity != nil {
		for k, v := range status.Capacity {
			capacity[string(k)] = v
		}
	}
	var conditions []PersistentVolumeClaimCondition
	if status.Conditions != nil {
		conditions = make([]PersistentVolumeClaimCondition, len(status.Conditions))
		for i, c := range status.Conditions {
			conditions[i] = ConvertPersistentVolumeClaimCondition(c)
		}
	}
	allocatedResources := make(map[string]resource.Quantity)
	if status.AllocatedResources != nil {
		for k, v := range status.AllocatedResources {
			allocatedResources[string(k)] = v
		}
	}
	allocatedResourceStatus := make(map[string]string)
	if status.AllocatedResourceStatuses != nil {
		allocatedResourceStatus = make(map[string]string)
		for k, v := range status.AllocatedResourceStatuses {
			allocatedResourceStatus[string(k)] = string(v)
		}
	}
	modifyVolumeStatus := struct {
		TargetVolumeAttributesClassName string
		Status                          string
	}{}
	if status.ModifyVolumeStatus != nil {
		modifyVolumeStatus.TargetVolumeAttributesClassName = status.ModifyVolumeStatus.TargetVolumeAttributesClassName
		modifyVolumeStatus.Status = string(status.ModifyVolumeStatus.Status)
	}
	return PersistentVolumeClaimStatus{
		Phase:                            string(status.Phase),
		AccessModes:                      accessModes,
		Capacity:                         capacity,
		Conditions:                       conditions,
		AllocatedResources:               allocatedResources,
		AllocatedResourceStatuses:        allocatedResourceStatus,
		CurrentVolumeAttributesClassName: status.CurrentVolumeAttributesClassName,
		ModifyVolumeStatus:               &modifyVolumeStatus,
	}
}

type PersistentVolumeClaimCondition struct {
	Type               string // corev1.PersistentVolumeClaimConditionType
	Status             string // corev1.ConditionStatus
	LastProbeTime      time.Time
	LastTransitionTime time.Time
	Reason             string
	Message            string
}

func ConvertPersistentVolumeClaimCondition(condition corev1.PersistentVolumeClaimCondition) PersistentVolumeClaimCondition {
	return PersistentVolumeClaimCondition{
		Type:               string(condition.Type),
		Status:             string(condition.Status),
		LastProbeTime:      ConvertTime(condition.LastProbeTime),
		LastTransitionTime: ConvertTime(condition.LastTransitionTime),
		Reason:             condition.Reason,
		Message:            condition.Message,
	}
}

// Add necessary conversion functions for PVC Status

// --- Adding More Types ---

// --- Service ---
type Service struct {
	TypeMeta
	ObjectMeta
	Spec   ServiceSpec
	Status ServiceStatus
}

type ServiceSpec struct {
	Ports                         []ServicePort
	Selector                      map[string]string // Keep as is
	ClusterIP                     string
	ClusterIPs                    []string
	Type                          string // corev1.ServiceType
	ExternalIPs                   []string
	SessionAffinity               string // corev1.ServiceAffinity
	LoadBalancerIP                string
	LoadBalancerSourceRanges      []string
	ExternalName                  string
	ExternalTrafficPolicy         string // corev1.ServiceExternalTrafficPolicyType
	HealthCheckNodePort           int32
	PublishNotReadyAddresses      bool
	SessionAffinityConfig         *SessionAffinityConfig
	IPFamilies                    []string // corev1.IPFamily
	IPFamilyPolicy                *string  // corev1.IPFamilyPolicyType
	AllocateLoadBalancerNodePorts *bool
	LoadBalancerClass             *string
	InternalTrafficPolicy         *string // corev1.ServiceInternalTrafficPolicyType
}

type ServicePort struct {
	Name        string
	Protocol    string // corev1.Protocol
	AppProtocol *string
	Port        int32
	TargetPort  intstr.IntOrString // Keep as IntOrString
	NodePort    int32
}

type SessionAffinityConfig struct {
	ClientIP *ClientIPConfig
}

type ClientIPConfig struct {
	TimeoutSeconds *int32
}

type ServiceStatus struct {
	LoadBalancer LoadBalancerStatus
	Conditions   []metav1.Condition // Use metav1.Condition directly for simplicity
}

type LoadBalancerStatus struct {
	Ingress []LoadBalancerIngress
}

type LoadBalancerIngress struct {
	IP       string
	Hostname string
	Ports    []PortStatus
}

type PortStatus struct {
	Port     int32
	Protocol string // corev1.Protocol
	Error    *string
}

func ConvertServiceType(st corev1.ServiceType) string {
	return string(st)
}

func ConvertServiceAffinity(sa corev1.ServiceAffinity) string {
	return string(sa)
}

func ConvertServiceExternalTrafficPolicyType(setpt corev1.ServiceExternalTrafficPolicyType) string {
	return string(setpt)
}

func ConvertIPFamily(ipf corev1.IPFamily) string {
	return string(ipf)
}

func ConvertIPFamilies(ipfs []corev1.IPFamily) []string {
	if ipfs == nil {
		return nil
	}
	res := make([]string, len(ipfs))
	for i, f := range ipfs {
		res[i] = ConvertIPFamily(f)
	}
	return res
}

func ConvertIPFamilyPolicyType(ipfpt *corev1.IPFamilyPolicyType) *string {
	if ipfpt == nil {
		return nil
	}
	s := string(*ipfpt)
	return &s
}

func ConvertServiceInternalTrafficPolicyType(sitpt *corev1.ServiceInternalTrafficPolicyType) *string {
	if sitpt == nil {
		return nil
	}
	s := string(*sitpt)
	return &s
}

func ConvertServicePort(sp corev1.ServicePort) ServicePort {
	return ServicePort{
		Name:        sp.Name,
		Protocol:    ConvertProtocol(sp.Protocol),
		AppProtocol: sp.AppProtocol,
		Port:        sp.Port,
		TargetPort:  sp.TargetPort, // Keep IntOrString
		NodePort:    sp.NodePort,
	}
}

func ConvertServicePorts(ports []corev1.ServicePort) []ServicePort {
	if ports == nil {
		return nil
	}
	res := make([]ServicePort, len(ports))
	for i, p := range ports {
		res[i] = ConvertServicePort(p)
	}
	return res
}

func ConvertClientIPConfig(cic *corev1.ClientIPConfig) *ClientIPConfig {
	if cic == nil {
		return nil
	}
	return &ClientIPConfig{
		TimeoutSeconds: cic.TimeoutSeconds,
	}
}

func ConvertSessionAffinityConfig(sac *corev1.SessionAffinityConfig) *SessionAffinityConfig {
	if sac == nil {
		return nil
	}
	return &SessionAffinityConfig{
		ClientIP: ConvertClientIPConfig(sac.ClientIP),
	}
}

func ConvertServiceSpec(spec corev1.ServiceSpec) ServiceSpec {
	return ServiceSpec{
		Ports:                         ConvertServicePorts(spec.Ports),
		Selector:                      spec.Selector,
		ClusterIP:                     spec.ClusterIP,
		ClusterIPs:                    spec.ClusterIPs,
		Type:                          ConvertServiceType(spec.Type),
		ExternalIPs:                   spec.ExternalIPs,
		SessionAffinity:               ConvertServiceAffinity(spec.SessionAffinity),
		LoadBalancerIP:                spec.LoadBalancerIP,
		LoadBalancerSourceRanges:      spec.LoadBalancerSourceRanges,
		ExternalName:                  spec.ExternalName,
		ExternalTrafficPolicy:         ConvertServiceExternalTrafficPolicyType(spec.ExternalTrafficPolicy),
		HealthCheckNodePort:           spec.HealthCheckNodePort,
		PublishNotReadyAddresses:      spec.PublishNotReadyAddresses,
		SessionAffinityConfig:         ConvertSessionAffinityConfig(spec.SessionAffinityConfig),
		IPFamilies:                    ConvertIPFamilies(spec.IPFamilies),
		IPFamilyPolicy:                ConvertIPFamilyPolicyType(spec.IPFamilyPolicy),
		AllocateLoadBalancerNodePorts: spec.AllocateLoadBalancerNodePorts,
		LoadBalancerClass:             spec.LoadBalancerClass,
		InternalTrafficPolicy:         ConvertServiceInternalTrafficPolicyType(spec.InternalTrafficPolicy),
	}
}

func ConvertPortStatus(ps corev1.PortStatus) PortStatus {
	return PortStatus{
		Port:     ps.Port,
		Protocol: ConvertProtocol(ps.Protocol),
		Error:    ps.Error,
	}
}

func ConvertPortStatuses(pss []corev1.PortStatus) []PortStatus {
	if pss == nil {
		return nil
	}
	res := make([]PortStatus, len(pss))
	for i, ps := range pss {
		res[i] = ConvertPortStatus(ps)
	}
	return res
}

func ConvertLoadBalancerIngress(lbi corev1.LoadBalancerIngress) LoadBalancerIngress {
	return LoadBalancerIngress{
		IP:       lbi.IP,
		Hostname: lbi.Hostname,
		Ports:    ConvertPortStatuses(lbi.Ports),
	}
}

func ConvertLoadBalancerIngresses(lbi []corev1.LoadBalancerIngress) []LoadBalancerIngress {
	if lbi == nil {
		return nil
	}
	res := make([]LoadBalancerIngress, len(lbi))
	for i, ing := range lbi {
		res[i] = ConvertLoadBalancerIngress(ing)
	}
	return res
}

func ConvertLoadBalancerStatus(lbs corev1.LoadBalancerStatus) LoadBalancerStatus {
	return LoadBalancerStatus{
		Ingress: ConvertLoadBalancerIngresses(lbs.Ingress),
	}
}

func ConvertServiceStatus(status corev1.ServiceStatus) ServiceStatus {
	return ServiceStatus{
		LoadBalancer: ConvertLoadBalancerStatus(status.LoadBalancer),
		Conditions:   status.Conditions, // Use metav1.Condition directly
	}
}

func ConvertService(svc *corev1.Service) Service {
	return Service{
		TypeMeta:   ConvertTypeMeta(svc.TypeMeta),
		ObjectMeta: ConvertObjectMeta(&svc.ObjectMeta),
		Spec:       ConvertServiceSpec(svc.Spec),
		Status:     ConvertServiceStatus(svc.Status),
	}
}

// --- Ingress (networkingv1) ---
type Ingress struct {
	TypeMeta
	ObjectMeta
	Spec   IngressSpec
	Status IngressStatus
}

type IngressSpec struct {
	IngressClassName *string
	DefaultBackend   *IngressBackend
	TLS              []IngressTLS
	Rules            []IngressRule
}

type IngressBackend struct {
	Service  *IngressServiceBackend
	Resource *TypedLocalObjectReference // Use helpers.TypedLocalObjectReference
}

type IngressServiceBackend struct {
	Name string
	Port ServiceBackendPort
}

func ConvertIngressServiceBackend(sb *networkingv1.IngressServiceBackend) *IngressServiceBackend {
	if sb == nil {
		return nil
	}
	return &IngressServiceBackend{
		Name: sb.Name,
		Port: ServiceBackendPort{
			Name:   sb.Port.Name,
			Number: sb.Port.Number,
		},
	}
}

type ServiceBackendPort struct {
	Name   string
	Number int32
}

type IngressTLS struct {
	Hosts      []string
	SecretName string
}

type IngressRule struct {
	Host string
	IngressRuleValue
}

type IngressRuleValue struct {
	HTTP *HTTPIngressRuleValue
}

func ConvertIngressRuleValue(rv networkingv1.IngressRuleValue) IngressRuleValue {
	return IngressRuleValue{
		HTTP: ConvertHTTPIngressRuleValue(rv.HTTP),
	}
}

type HTTPIngressRuleValue struct {
	Paths []HTTPIngressPath
}

func ConvertHTTPIngressRuleValue(rv *networkingv1.HTTPIngressRuleValue) *HTTPIngressRuleValue {
	if rv == nil {
		return nil
	}
	paths := make([]HTTPIngressPath, len(rv.Paths))
	for i, path := range rv.Paths {
		paths[i] = ConvertHTTPIngressPath(path)
	}
	return &HTTPIngressRuleValue{
		Paths: paths,
	}
}

type HTTPIngressPath struct {
	Path     string
	PathType *string // networkingv1.PathType
	Backend  IngressBackend
}

func ConvertHTTPIngressPath(p networkingv1.HTTPIngressPath) HTTPIngressPath {
	var pathType *string
	if p.PathType != nil {
		pathTypeTmp := string(*p.PathType)
		pathType = &pathTypeTmp
	}
	return HTTPIngressPath{
		Path:     p.Path,
		PathType: pathType,
		Backend:  *ConvertIngressBackend(&p.Backend),
	}
}

type IngressStatus struct {
	LoadBalancer IngressLoadBalancerStatus // Changed type
}

// Added IngressLoadBalancerStatus helper type
// ConvertIngressPortStatus converts networkingv1.IngressPortStatus to helpers.PortStatus
func ConvertIngressPortStatus(ps networkingv1.IngressPortStatus) PortStatus {
	return PortStatus{
		Port:     ps.Port,
		Protocol: ConvertProtocol(ps.Protocol),
		Error:    ps.Error,
	}
}

func ConvertIngressPortStatuses(pss []networkingv1.IngressPortStatus) []PortStatus {
	if pss == nil {
		return nil
	}
	res := make([]PortStatus, len(pss))
	for i, ps := range pss {
		res[i] = ConvertIngressPortStatus(ps)
	}
	return res
}

func ConvertIngressLoadBalancerIngress(lbi networkingv1.IngressLoadBalancerIngress) LoadBalancerIngress {
	return LoadBalancerIngress{
		IP:       lbi.IP,
		Hostname: lbi.Hostname,
		Ports:    ConvertIngressPortStatuses(lbi.Ports),
	}
}

func ConvertIngressLoadBalancerIngresses(lbi []networkingv1.IngressLoadBalancerIngress) []LoadBalancerIngress {
	if lbi == nil {
		return nil
	}
	res := make([]LoadBalancerIngress, len(lbi))
	for i, ing := range lbi {
		res[i] = ConvertIngressLoadBalancerIngress(ing)
	}
	return res
}

func ConvertIngressLoadBalancerStatus(lbs networkingv1.IngressLoadBalancerStatus) IngressLoadBalancerStatus {
	return IngressLoadBalancerStatus{
		Ingress: ConvertIngressLoadBalancerIngresses(lbs.Ingress),
	}
}

func ConvertIngressSpec(spec networkingv1.IngressSpec) IngressSpec {
	return IngressSpec{
		IngressClassName: spec.IngressClassName,
		DefaultBackend:   ConvertIngressBackend(spec.DefaultBackend),
		TLS:              ConvertIngressTLSs(spec.TLS),
		Rules:            ConvertIngressRules(spec.Rules),
	}
}

func ConvertIngressStatus(status networkingv1.IngressStatus) IngressStatus {
	return IngressStatus{
		LoadBalancer: ConvertIngressLoadBalancerStatus(status.LoadBalancer), // Use new function
	}
}

func ConvertIngress(ing *networkingv1.Ingress) Ingress {
	return Ingress{
		TypeMeta:   ConvertTypeMeta(ing.TypeMeta),
		ObjectMeta: ConvertObjectMeta(&ing.ObjectMeta),
		Spec:       ConvertIngressSpec(ing.Spec),
		Status:     ConvertIngressStatus(ing.Status),
	}
}

// --- PersistentVolume (corev1) ---
type PersistentVolume struct {
	TypeMeta
	ObjectMeta
	Spec   PersistentVolumeSpec
	Status PersistentVolumeStatus
}

type PersistentVolumeSpec struct {
	Capacity                      map[string]resource.Quantity // corev1.ResourceList
	PersistentVolumeSource        VolumeSource                 // Use helpers.VolumeSource
	AccessModes                   []string                     // corev1.PersistentVolumeAccessMode
	ClaimRef                      *ObjectReference             // Use helpers.ObjectReference
	PersistentVolumeReclaimPolicy string                       // corev1.PersistentVolumeReclaimPolicy
	StorageClassName              string
	MountOptions                  []string
	VolumeMode                    *string // corev1.PersistentVolumeMode
	NodeAffinity                  *VolumeNodeAffinity
}

type VolumeNodeAffinity struct {
	Required *NodeSelector // Use helpers.NodeSelector
}

type PersistentVolumeStatus struct {
	Phase   string // corev1.PersistentVolumePhase
	Message string
	Reason  string
}

func ConvertPersistentVolumeReclaimPolicy(pvcp corev1.PersistentVolumeReclaimPolicy) string {
	return string(pvcp)
}

func ConvertVolumeNodeAffinity(vna *corev1.VolumeNodeAffinity) *VolumeNodeAffinity {
	if vna == nil {
		return nil
	}
	return &VolumeNodeAffinity{
		Required: ConvertNodeSelector(vna.Required),
	}
}

// Assumes ConvertVolumeSource exists in volume.go or needs definition
func ConvertPersistentVolumeSource(pvs corev1.PersistentVolumeSource) VolumeSource {
	// NOTE: Assumes Convert* functions for each volume type exist elsewhere (e.g., volume.go)
	//       and return the corresponding helper type.
	return VolumeSource{
		// Restore all volume source type conversions
		HostPath: ConvertHostPathVolumeSource(pvs.HostPath),
		NFS:      ConvertNFSVolumeSource(pvs.NFS),
		FC:       ConvertFCVolumeSource(pvs.FC),
	}
}

func ConvertPersistentVolumeSpec(spec corev1.PersistentVolumeSpec) PersistentVolumeSpec {
	return PersistentVolumeSpec{
		Capacity:                      ConvertResourceList(spec.Capacity),
		PersistentVolumeSource:        ConvertPersistentVolumeSource(spec.PersistentVolumeSource), // Needs careful implementation
		AccessModes:                   ConvertAccessModes(spec.AccessModes),
		ClaimRef:                      ConvertObjectReferencePointer(spec.ClaimRef),
		PersistentVolumeReclaimPolicy: ConvertPersistentVolumeReclaimPolicy(spec.PersistentVolumeReclaimPolicy),
		StorageClassName:              spec.StorageClassName,
		MountOptions:                  spec.MountOptions,
		VolumeMode:                    ConvertVolumeMode(spec.VolumeMode),
		NodeAffinity:                  ConvertVolumeNodeAffinity(spec.NodeAffinity),
	}
}

func ConvertPersistentVolumePhase(pvp corev1.PersistentVolumePhase) string {
	return string(pvp)
}

func ConvertPersistentVolumeStatus(status corev1.PersistentVolumeStatus) PersistentVolumeStatus {
	return PersistentVolumeStatus{
		Phase:   ConvertPersistentVolumePhase(status.Phase),
		Message: status.Message,
		Reason:  status.Reason,
	}
}

func ConvertPersistentVolume(pv *corev1.PersistentVolume) PersistentVolume {
	return PersistentVolume{
		TypeMeta:   ConvertTypeMeta(pv.TypeMeta),
		ObjectMeta: ConvertObjectMeta(&pv.ObjectMeta),
		Spec:       ConvertPersistentVolumeSpec(pv.Spec),
		Status:     ConvertPersistentVolumeStatus(pv.Status),
	}
}

// Helper needed for ClaimRef which is *corev1.ObjectReference
func ConvertObjectReferencePointer(ref *corev1.ObjectReference) *ObjectReference {
	if ref == nil {
		return nil
	}
	convertedRef := ConvertObjectReference(*ref)
	return &convertedRef
}

// --- Node (corev1) ---
type Node struct {
	TypeMeta
	ObjectMeta
	Spec   NodeSpec
	Status NodeStatus
}

type NodeSpec struct {
	PodCIDR       string
	PodCIDRs      []string
	ProviderID    string
	Unschedulable bool
	Taints        []Taint
	ConfigSource  *NodeConfigSource
}

type Taint struct {
	Key       string
	Value     string
	Effect    string     // corev1.TaintEffect
	TimeAdded *time.Time // Use helpers.Time via ConvertTimePtr
}

type NodeConfigSource struct {
	ConfigMap *ConfigMapNodeConfigSource
}

type ConfigMapNodeConfigSource struct {
	Namespace        string
	Name             string
	UID              types.UID
	KubeletConfigKey string
	ResourceVersion  string
}

type NodeStatus struct {
	Capacity        map[string]resource.Quantity // corev1.ResourceList
	Allocatable     map[string]resource.Quantity // corev1.ResourceList
	Phase           string                       // corev1.NodePhase
	Conditions      []NodeCondition
	Addresses       []NodeAddress
	DaemonEndpoints NodeDaemonEndpoints
	NodeInfo        NodeSystemInfo
	Images          []ContainerImage
	VolumesInUse    []string // corev1.UniqueVolumeName
	VolumesAttached []AttachedVolume
	Config          *NodeConfigStatus
	RuntimeHandlers []NodeRuntimeHandler
}

type NodeCondition struct {
	Type               string     // corev1.NodeConditionType
	Status             string     // corev1.ConditionStatus
	LastHeartbeatTime  *time.Time // Use helpers.Time via ConvertTimePtr
	LastTransitionTime *time.Time // Use helpers.Time via ConvertTimePtr
	Reason             string
	Message            string
}

type NodeAddress struct {
	Type    string // corev1.NodeAddressType
	Address string
}

type NodeDaemonEndpoints struct {
	KubeletEndpoint DaemonEndpoint
}

type DaemonEndpoint struct {
	Port int32
}

type NodeSystemInfo struct {
	MachineID               string
	SystemUUID              string
	BootID                  string
	KernelVersion           string
	OSImage                 string
	ContainerRuntimeVersion string
	KubeletVersion          string
	KubeProxyVersion        string
	OperatingSystem         string
	Architecture            string
}

type ContainerImage struct {
	Names     []string
	SizeBytes int64
}

type AttachedVolume struct {
	Name       string // corev1.UniqueVolumeName
	DevicePath string
}

type NodeConfigStatus struct {
	Assigned      *NodeConfigSource
	Active        *NodeConfigSource
	LastKnownGood *NodeConfigSource
	Error         string
}

type NodeRuntimeHandler struct {
	Name     string
	Features *NodeRuntimeHandlerFeatures
}

type NodeRuntimeHandlerFeatures struct {
	RecursiveReadOnlyMounts *bool
}

func ConvertTaintEffect(te corev1.TaintEffect) string {
	return string(te)
}

func ConvertTaint(t corev1.Taint) Taint {
	return Taint{
		Key:       t.Key,
		Value:     t.Value,
		Effect:    ConvertTaintEffect(t.Effect),
		TimeAdded: ConvertTimePtr(t.TimeAdded),
	}
}

func ConvertTaints(ts []corev1.Taint) []Taint {
	if ts == nil {
		return nil
	}
	res := make([]Taint, len(ts))
	for i, t := range ts {
		res[i] = ConvertTaint(t)
	}
	return res
}

func ConvertConfigMapNodeConfigSource(cmns *corev1.ConfigMapNodeConfigSource) *ConfigMapNodeConfigSource {
	if cmns == nil {
		return nil
	}
	return &ConfigMapNodeConfigSource{
		Namespace:        cmns.Namespace,
		Name:             cmns.Name,
		UID:              cmns.UID,
		KubeletConfigKey: cmns.KubeletConfigKey,
		ResourceVersion:  cmns.ResourceVersion,
	}
}

func ConvertNodeConfigSource(ncs *corev1.NodeConfigSource) *NodeConfigSource {
	if ncs == nil {
		return nil
	}
	return &NodeConfigSource{
		ConfigMap: ConvertConfigMapNodeConfigSource(ncs.ConfigMap),
	}
}

func ConvertNodeSpec(spec corev1.NodeSpec) NodeSpec {
	return NodeSpec{
		PodCIDR:       spec.PodCIDR,
		PodCIDRs:      spec.PodCIDRs,
		ProviderID:    spec.ProviderID,
		Unschedulable: spec.Unschedulable,
		Taints:        ConvertTaints(spec.Taints),
		ConfigSource:  ConvertNodeConfigSource(spec.ConfigSource),
	}
}

func ConvertNodePhase(np corev1.NodePhase) string {
	return string(np)
}

func ConvertNodeConditionType(nct corev1.NodeConditionType) string {
	return string(nct)
}

func ConvertNodeCondition(c corev1.NodeCondition) NodeCondition {
	return NodeCondition{
		Type:               ConvertNodeConditionType(c.Type),
		Status:             ConvertConditionStatus(c.Status),
		LastHeartbeatTime:  ConvertTimePtr(&c.LastHeartbeatTime),
		LastTransitionTime: ConvertTimePtr(&c.LastTransitionTime),
		Reason:             c.Reason,
		Message:            c.Message,
	}
}

func ConvertNodeConditions(conds []corev1.NodeCondition) []NodeCondition {
	if conds == nil {
		return nil
	}
	res := make([]NodeCondition, len(conds))
	for i, c := range conds {
		res[i] = ConvertNodeCondition(c)
	}
	return res
}

func ConvertNodeAddressType(nat corev1.NodeAddressType) string {
	return string(nat)
}

func ConvertNodeAddress(na corev1.NodeAddress) NodeAddress {
	return NodeAddress{
		Type:    ConvertNodeAddressType(na.Type),
		Address: na.Address,
	}
}

func ConvertNodeAddresses(nas []corev1.NodeAddress) []NodeAddress {
	if nas == nil {
		return nil
	}
	res := make([]NodeAddress, len(nas))
	for i, na := range nas {
		res[i] = ConvertNodeAddress(na)
	}
	return res
}

func ConvertDaemonEndpoint(de corev1.DaemonEndpoint) DaemonEndpoint {
	return DaemonEndpoint{
		Port: de.Port,
	}
}

func ConvertNodeDaemonEndpoints(nde corev1.NodeDaemonEndpoints) NodeDaemonEndpoints {
	return NodeDaemonEndpoints{
		KubeletEndpoint: ConvertDaemonEndpoint(nde.KubeletEndpoint),
	}
}

func ConvertNodeSystemInfo(nsi corev1.NodeSystemInfo) NodeSystemInfo {
	return NodeSystemInfo{
		MachineID:               nsi.MachineID,
		SystemUUID:              nsi.SystemUUID,
		BootID:                  nsi.BootID,
		KernelVersion:           nsi.KernelVersion,
		OSImage:                 nsi.OSImage,
		ContainerRuntimeVersion: nsi.ContainerRuntimeVersion,
		KubeletVersion:          nsi.KubeletVersion,
		KubeProxyVersion:        nsi.KubeProxyVersion,
		OperatingSystem:         nsi.OperatingSystem,
		Architecture:            nsi.Architecture,
	}
}

func ConvertContainerImage(ci corev1.ContainerImage) ContainerImage {
	return ContainerImage{
		Names:     ci.Names,
		SizeBytes: ci.SizeBytes,
	}
}

func ConvertContainerImages(cis []corev1.ContainerImage) []ContainerImage {
	if cis == nil {
		return nil
	}
	res := make([]ContainerImage, len(cis))
	for i, ci := range cis {
		res[i] = ConvertContainerImage(ci)
	}
	return res
}

func ConvertUniqueVolumeName(uvn corev1.UniqueVolumeName) string {
	return string(uvn)
}

func ConvertUniqueVolumeNames(uvns []corev1.UniqueVolumeName) []string {
	if uvns == nil {
		return nil
	}
	res := make([]string, len(uvns))
	for i, uvn := range uvns {
		res[i] = ConvertUniqueVolumeName(uvn)
	}
	return res
}

func ConvertAttachedVolume(av corev1.AttachedVolume) AttachedVolume {
	return AttachedVolume{
		Name:       ConvertUniqueVolumeName(av.Name),
		DevicePath: av.DevicePath,
	}
}

func ConvertAttachedVolumes(avs []corev1.AttachedVolume) []AttachedVolume {
	if avs == nil {
		return nil
	}
	res := make([]AttachedVolume, len(avs))
	for i, av := range avs {
		res[i] = ConvertAttachedVolume(av)
	}
	return res
}

func ConvertNodeConfigStatus(ncs *corev1.NodeConfigStatus) *NodeConfigStatus {
	if ncs == nil {
		return nil
	}
	return &NodeConfigStatus{
		Assigned:      ConvertNodeConfigSource(ncs.Assigned),
		Active:        ConvertNodeConfigSource(ncs.Active),
		LastKnownGood: ConvertNodeConfigSource(ncs.LastKnownGood),
		Error:         ncs.Error,
	}
}

func ConvertNodeRuntimeHandlerFeatures(f *corev1.NodeRuntimeHandlerFeatures) *NodeRuntimeHandlerFeatures {
	if f == nil {
		return nil
	}
	return &NodeRuntimeHandlerFeatures{
		RecursiveReadOnlyMounts: f.RecursiveReadOnlyMounts,
	}
}

func ConvertNodeRuntimeHandler(h corev1.NodeRuntimeHandler) NodeRuntimeHandler {
	return NodeRuntimeHandler{
		Name:     h.Name,
		Features: ConvertNodeRuntimeHandlerFeatures(h.Features),
	}
}

func ConvertNodeRuntimeHandlers(hs []corev1.NodeRuntimeHandler) []NodeRuntimeHandler {
	if hs == nil {
		return nil
	}
	res := make([]NodeRuntimeHandler, len(hs))
	for i, h := range hs {
		res[i] = ConvertNodeRuntimeHandler(h)
	}
	return res
}

func ConvertNodeStatus(status corev1.NodeStatus) NodeStatus {
	return NodeStatus{
		Capacity:        ConvertResourceList(status.Capacity),
		Allocatable:     ConvertResourceList(status.Allocatable),
		Phase:           ConvertNodePhase(status.Phase),
		Conditions:      ConvertNodeConditions(status.Conditions),
		Addresses:       ConvertNodeAddresses(status.Addresses),
		DaemonEndpoints: ConvertNodeDaemonEndpoints(status.DaemonEndpoints),
		NodeInfo:        ConvertNodeSystemInfo(status.NodeInfo),
		Images:          ConvertContainerImages(status.Images),
		VolumesInUse:    ConvertUniqueVolumeNames(status.VolumesInUse),
		VolumesAttached: ConvertAttachedVolumes(status.VolumesAttached),
		Config:          ConvertNodeConfigStatus(status.Config),
		RuntimeHandlers: ConvertNodeRuntimeHandlers(status.RuntimeHandlers),
	}
}

func ConvertNode(node *corev1.Node) Node {
	return Node{
		TypeMeta:   ConvertTypeMeta(node.TypeMeta),
		ObjectMeta: ConvertObjectMeta(&node.ObjectMeta),
		Spec:       ConvertNodeSpec(node.Spec),
		Status:     ConvertNodeStatus(node.Status),
	}
}

// --- LimitRange (corev1) ---
type LimitRange struct {
	TypeMeta
	ObjectMeta
	Spec LimitRangeSpec
}

type LimitRangeSpec struct {
	Limits []LimitRangeItem
}

type LimitRangeItem struct {
	Type                 string                       // corev1.LimitType
	Max                  map[string]resource.Quantity // corev1.ResourceList
	Min                  map[string]resource.Quantity // corev1.ResourceList
	Default              map[string]resource.Quantity // corev1.ResourceList
	DefaultRequest       map[string]resource.Quantity // corev1.ResourceList
	MaxLimitRequestRatio map[string]resource.Quantity // corev1.ResourceList
}

func ConvertLimitType(lt corev1.LimitType) string {
	return string(lt)
}

func ConvertLimitRangeItem(lri corev1.LimitRangeItem) LimitRangeItem {
	return LimitRangeItem{
		Type:                 ConvertLimitType(lri.Type),
		Max:                  ConvertResourceList(lri.Max),
		Min:                  ConvertResourceList(lri.Min),
		Default:              ConvertResourceList(lri.Default),
		DefaultRequest:       ConvertResourceList(lri.DefaultRequest),
		MaxLimitRequestRatio: ConvertResourceList(lri.MaxLimitRequestRatio),
	}
}

func ConvertLimitRangeItems(lris []corev1.LimitRangeItem) []LimitRangeItem {
	if lris == nil {
		return nil
	}
	res := make([]LimitRangeItem, len(lris))
	for i, lri := range lris {
		res[i] = ConvertLimitRangeItem(lri)
	}
	return res
}

func ConvertLimitRangeSpec(spec corev1.LimitRangeSpec) LimitRangeSpec {
	return LimitRangeSpec{
		Limits: ConvertLimitRangeItems(spec.Limits),
	}
}

func ConvertLimitRange(lr *corev1.LimitRange) LimitRange {
	return LimitRange{
		TypeMeta:   ConvertTypeMeta(lr.TypeMeta),
		ObjectMeta: ConvertObjectMeta(&lr.ObjectMeta),
		Spec:       ConvertLimitRangeSpec(lr.Spec),
	}
}

// --- ResourceQuota (corev1) ---
type ResourceQuota struct {
	TypeMeta
	ObjectMeta
	Spec   ResourceQuotaSpec
	Status ResourceQuotaStatus
}

type ResourceQuotaSpec struct {
	Hard          map[string]resource.Quantity // corev1.ResourceList
	Scopes        []string                     // corev1.ResourceQuotaScope
	ScopeSelector *ScopeSelector
}

// Completing ScopeSelector definition
type ScopeSelector struct {
	MatchExpressions []ScopedResourceSelectorRequirement
}

type ScopedResourceSelectorRequirement struct {
	ScopeName string // corev1.ResourceQuotaScope
	Operator  string // corev1.ScopeSelectorOperator
	Values    []string
}

type ResourceQuotaStatus struct {
	Hard map[string]resource.Quantity // corev1.ResourceList
	Used map[string]resource.Quantity // corev1.ResourceList
}

func ConvertResourceQuotaScope(rqs corev1.ResourceQuotaScope) string {
	return string(rqs)
}

func ConvertResourceQuotaScopes(rqss []corev1.ResourceQuotaScope) []string {
	if rqss == nil {
		return nil
	}
	res := make([]string, len(rqss))
	for i, s := range rqss {
		res[i] = ConvertResourceQuotaScope(s)
	}
	return res
}

func ConvertScopeSelectorOperator(sso corev1.ScopeSelectorOperator) string {
	return string(sso)
}

func ConvertScopedResourceSelectorRequirement(sr corev1.ScopedResourceSelectorRequirement) ScopedResourceSelectorRequirement {
	return ScopedResourceSelectorRequirement{
		ScopeName: ConvertResourceQuotaScope(sr.ScopeName),
		Operator:  ConvertScopeSelectorOperator(sr.Operator),
		Values:    sr.Values,
	}
}

func ConvertScopedResourceSelectorRequirements(srs []corev1.ScopedResourceSelectorRequirement) []ScopedResourceSelectorRequirement {
	if srs == nil {
		return nil
	}
	res := make([]ScopedResourceSelectorRequirement, len(srs))
	for i, sr := range srs {
		res[i] = ConvertScopedResourceSelectorRequirement(sr)
	}
	return res
}

func ConvertScopeSelector(ss *corev1.ScopeSelector) *ScopeSelector {
	if ss == nil {
		return nil
	}
	return &ScopeSelector{
		MatchExpressions: ConvertScopedResourceSelectorRequirements(ss.MatchExpressions),
	}
}

func ConvertResourceQuotaSpec(spec corev1.ResourceQuotaSpec) ResourceQuotaSpec {
	return ResourceQuotaSpec{
		Hard:          ConvertResourceList(spec.Hard),
		Scopes:        ConvertResourceQuotaScopes(spec.Scopes),
		ScopeSelector: ConvertScopeSelector(spec.ScopeSelector),
	}
}

func ConvertResourceQuotaStatus(status corev1.ResourceQuotaStatus) ResourceQuotaStatus {
	return ResourceQuotaStatus{
		Hard: ConvertResourceList(status.Hard),
		Used: ConvertResourceList(status.Used),
	}
}

func ConvertResourceQuota(rq *corev1.ResourceQuota) ResourceQuota {
	return ResourceQuota{
		TypeMeta:   ConvertTypeMeta(rq.TypeMeta),
		ObjectMeta: ConvertObjectMeta(&rq.ObjectMeta),
		Spec:       ConvertResourceQuotaSpec(rq.Spec),
		Status:     ConvertResourceQuotaStatus(rq.Status),
	}
}

// --- PodTemplate (corev1) ---
type PodTemplate struct {
	TypeMeta
	ObjectMeta
	Template PodTemplateSpec // Use helpers.PodTemplateSpec
}

func ConvertPodTemplate(pt *corev1.PodTemplate) PodTemplate {
	return PodTemplate{
		TypeMeta:   ConvertTypeMeta(pt.TypeMeta),
		ObjectMeta: ConvertObjectMeta(&pt.ObjectMeta),
		Template:   ConvertPodTemplateSpec(pt.Template), // Assumes this exists (pod.go)
	}
}

// --- ReplicationController (corev1) ---
type ReplicationController struct {
	TypeMeta
	ObjectMeta
	Spec   ReplicationControllerSpec
	Status ReplicationControllerStatus
}

type ReplicationControllerSpec struct {
	Replicas        *int32
	MinReadySeconds int32
	Selector        map[string]string // Keep as is
	Template        *PodTemplateSpec  // Use helpers.PodTemplateSpec (Pointer)
}

type ReplicationControllerStatus struct {
	Replicas             int32
	FullyLabeledReplicas int32
	ReadyReplicas        int32
	AvailableReplicas    int32
	ObservedGeneration   int64
	Conditions           []ReplicationControllerCondition
}

type ReplicationControllerCondition struct {
	Type               string    // corev1.ReplicationControllerConditionType
	Status             string    // corev1.ConditionStatus
	LastTransitionTime time.Time // Use helpers.Time via ConvertTime
	Reason             string
	Message            string
}

func ConvertReplicationControllerSpec(spec corev1.ReplicationControllerSpec) ReplicationControllerSpec {
	// Need to handle pointer for PodTemplateSpec
	var templateSpecPtr *PodTemplateSpec
	if spec.Template != nil {
		converted := ConvertPodTemplateSpec(*spec.Template) // Assumes exists (pod.go)
		templateSpecPtr = &converted
	}
	return ReplicationControllerSpec{
		Replicas:        spec.Replicas,
		MinReadySeconds: spec.MinReadySeconds,
		Selector:        spec.Selector,
		Template:        templateSpecPtr,
	}
}

func ConvertReplicationControllerConditionType(rcct corev1.ReplicationControllerConditionType) string {
	return string(rcct)
}

func ConvertReplicationControllerCondition(c corev1.ReplicationControllerCondition) ReplicationControllerCondition {
	return ReplicationControllerCondition{
		Type:               ConvertReplicationControllerConditionType(c.Type),
		Status:             ConvertConditionStatus(c.Status),
		LastTransitionTime: ConvertTime(c.LastTransitionTime),
		Reason:             c.Reason,
		Message:            c.Message,
	}
}

func ConvertReplicationControllerConditions(conds []corev1.ReplicationControllerCondition) []ReplicationControllerCondition {
	if conds == nil {
		return nil
	}
	res := make([]ReplicationControllerCondition, len(conds))
	for i, c := range conds {
		res[i] = ConvertReplicationControllerCondition(c)
	}
	return res
}

func ConvertReplicationControllerStatus(status corev1.ReplicationControllerStatus) ReplicationControllerStatus {
	return ReplicationControllerStatus{
		Replicas:             status.Replicas,
		FullyLabeledReplicas: status.FullyLabeledReplicas,
		ReadyReplicas:        status.ReadyReplicas,
		AvailableReplicas:    status.AvailableReplicas,
		ObservedGeneration:   status.ObservedGeneration,
		Conditions:           ConvertReplicationControllerConditions(status.Conditions),
	}
}

func ConvertReplicationController(rc *corev1.ReplicationController) ReplicationController {
	return ReplicationController{
		TypeMeta:   ConvertTypeMeta(rc.TypeMeta),
		ObjectMeta: ConvertObjectMeta(&rc.ObjectMeta),
		Spec:       ConvertReplicationControllerSpec(rc.Spec),
		Status:     ConvertReplicationControllerStatus(rc.Status),
	}
}

// --- StorageClass (storagev1) ---
type StorageClass struct {
	TypeMeta
	ObjectMeta
	Provisioner          string
	Parameters           map[string]string
	ReclaimPolicy        *string // corev1.PersistentVolumeReclaimPolicy
	MountOptions         []string
	AllowVolumeExpansion *bool
	VolumeBindingMode    *string // storagev1.VolumeBindingMode
	AllowedTopologies    []TopologySelectorTerm
}

type TopologySelectorTerm struct {
	MatchLabelExpressions []TopologySelectorLabelRequirement
}

type TopologySelectorLabelRequirement struct {
	Key    string
	Values []string
}

// Assumes ConvertPersistentVolumeReclaimPolicy defined with PV types

func ConvertVolumeBindingMode(vbm *storagev1.VolumeBindingMode) *string {
	if vbm == nil {
		return nil
	}
	s := string(*vbm)
	return &s
}

func ConvertStorageClass(sc *storagev1.StorageClass) StorageClass {
	// Need to handle *corev1.PersistentVolumeReclaimPolicy
	var reclaimPolicyPtr *string
	if sc.ReclaimPolicy != nil {
		rp := ConvertPersistentVolumeReclaimPolicy(*sc.ReclaimPolicy) // Assumes exists
		reclaimPolicyPtr = &rp
	}

	return StorageClass{
		TypeMeta:             ConvertTypeMeta(sc.TypeMeta),
		ObjectMeta:           ConvertObjectMeta(&sc.ObjectMeta),
		Provisioner:          sc.Provisioner,
		Parameters:           sc.Parameters,
		ReclaimPolicy:        reclaimPolicyPtr,
		MountOptions:         sc.MountOptions,
		AllowVolumeExpansion: sc.AllowVolumeExpansion,
		VolumeBindingMode:    ConvertVolumeBindingMode(sc.VolumeBindingMode),
	}
}

// --- Endpoints (corev1) ---
type Endpoints struct {
	TypeMeta
	ObjectMeta
	Subsets []EndpointSubset
}

type EndpointSubset struct {
	Addresses         []EndpointAddress
	NotReadyAddresses []EndpointAddress
	Ports             []EndpointPort
}

type EndpointAddress struct {
	IP        string
	Hostname  string
	NodeName  *string
	TargetRef *ObjectReference // Use helpers.ObjectReference
}

type EndpointPort struct {
	Name        string
	Port        int32
	Protocol    string // corev1.Protocol
	AppProtocol *string
}

func ConvertEndpointPort(ep corev1.EndpointPort) EndpointPort {
	return EndpointPort{
		Name:        ep.Name,
		Port:        ep.Port,
		Protocol:    ConvertProtocol(ep.Protocol), // Assumes exists (pod.go)
		AppProtocol: ep.AppProtocol,
	}
}

func ConvertEndpointPorts(eps []corev1.EndpointPort) []EndpointPort {
	if eps == nil {
		return nil
	}
	res := make([]EndpointPort, len(eps))
	for i, ep := range eps {
		res[i] = ConvertEndpointPort(ep)
	}
	return res
}

func ConvertEndpointAddress(ea corev1.EndpointAddress) EndpointAddress {
	return EndpointAddress{
		IP:        ea.IP,
		Hostname:  ea.Hostname,
		NodeName:  ea.NodeName,
		TargetRef: ConvertObjectReferencePointer(ea.TargetRef), // Assumes exists
	}
}

func ConvertEndpointAddresses(eas []corev1.EndpointAddress) []EndpointAddress {
	if eas == nil {
		return nil
	}
	res := make([]EndpointAddress, len(eas))
	for i, ea := range eas {
		res[i] = ConvertEndpointAddress(ea)
	}
	return res
}

func ConvertEndpointSubset(es corev1.EndpointSubset) EndpointSubset {
	return EndpointSubset{
		Addresses:         ConvertEndpointAddresses(es.Addresses),
		NotReadyAddresses: ConvertEndpointAddresses(es.NotReadyAddresses),
		Ports:             ConvertEndpointPorts(es.Ports),
	}
}

func ConvertEndpointSubsets(ess []corev1.EndpointSubset) []EndpointSubset {
	if ess == nil {
		return nil
	}
	res := make([]EndpointSubset, len(ess))
	for i, es := range ess {
		res[i] = ConvertEndpointSubset(es)
	}
	return res
}

func ConvertEndpoints(ep *corev1.Endpoints) Endpoints {
	return Endpoints{
		TypeMeta:   ConvertTypeMeta(ep.TypeMeta),
		ObjectMeta: ConvertObjectMeta(&ep.ObjectMeta),
		Subsets:    ConvertEndpointSubsets(ep.Subsets),
	}
}

// --- Event (corev1) ---
type Event struct {
	ObjectMeta
	InvolvedObject ObjectReference // Use helpers.ObjectReference
	Reason         string
	Message        string
	Source         EventSource
	FirstTimestamp time.Time // Use helpers.Time via ConvertTime
	LastTimestamp  time.Time // Use helpers.Time via ConvertTime
	Count          int32
	Type           string           // Warning, Normal etc.
	EventTime      metav1.MicroTime // Keep as metav1.MicroTime for now
	// Series and Action added later
	Series              *EventSeries
	Action              string
	Related             *ObjectReference // Use helpers.ObjectReference
	ReportingController string
	ReportingInstance   string
}

type EventSource struct {
	Component string
	Host      string
}

type EventSeries struct {
	Count            int32
	LastObservedTime metav1.MicroTime // Keep as metav1.MicroTime
}

func ConvertEventSource(es corev1.EventSource) EventSource {
	return EventSource{
		Component: es.Component,
		Host:      es.Host,
	}
}

func ConvertEventSeries(es *corev1.EventSeries) *EventSeries {
	if es == nil {
		return nil
	}
	return &EventSeries{
		Count:            es.Count,
		LastObservedTime: es.LastObservedTime, // Keep MicroTime
	}
}

func ConvertEvent(ev *corev1.Event) Event {
	return Event{
		ObjectMeta:          ConvertObjectMeta(&ev.ObjectMeta),
		InvolvedObject:      ConvertObjectReference(ev.InvolvedObject),
		Reason:              ev.Reason,
		Message:             ev.Message,
		Source:              ConvertEventSource(ev.Source),
		FirstTimestamp:      ConvertTime(ev.FirstTimestamp),
		LastTimestamp:       ConvertTime(ev.LastTimestamp),
		Count:               ev.Count,
		Type:                ev.Type,
		EventTime:           ev.EventTime, // Keep MicroTime
		Series:              ConvertEventSeries(ev.Series),
		Action:              ev.Action,
		Related:             ConvertObjectReferencePointer(ev.Related),
		ReportingController: ev.ReportingController,
		ReportingInstance:   ev.ReportingInstance,
	}
}

// --- NetworkPolicy (networkingv1) ---
type NetworkPolicy struct {
	TypeMeta
	ObjectMeta
	Spec NetworkPolicySpec
}

type NetworkPolicySpec struct {
	PodSelector LabelSelector // Use helpers.LabelSelector
	PolicyTypes []string      // networkingv1.PolicyType
	Ingress     []NetworkPolicyIngressRule
	Egress      []NetworkPolicyEgressRule
}

type NetworkPolicyIngressRule struct {
	Ports []NetworkPolicyPort
	From  []NetworkPolicyPeer
}

type NetworkPolicyEgressRule struct {
	Ports []NetworkPolicyPort
	To    []NetworkPolicyPeer
}

type NetworkPolicyPort struct {
	Protocol *string             // corev1.Protocol
	Port     *intstr.IntOrString // Keep IntOrString
	EndPort  *int32
}

type NetworkPolicyPeer struct {
	PodSelector       *LabelSelector // Use helpers.LabelSelector
	NamespaceSelector *LabelSelector // Use helpers.LabelSelector
	IPBlock           *IPBlock
}

type IPBlock struct {
	CIDR   string
	Except []string
}

type NetworkPolicyStatus struct {
	Conditions []metav1.Condition // Use metav1.Condition directly
}

func ConvertPolicyType(pt networkingv1.PolicyType) string {
	return string(pt)
}

func ConvertPolicyTypes(pts []networkingv1.PolicyType) []string {
	if pts == nil {
		return nil
	}
	res := make([]string, len(pts))
	for i, pt := range pts {
		res[i] = ConvertPolicyType(pt)
	}
	return res
}

func ConvertProtocolPointer(p *corev1.Protocol) *string {
	if p == nil {
		return nil
	}
	s := string(*p)
	return &s
}

func ConvertNetworkPolicyPort(npp networkingv1.NetworkPolicyPort) NetworkPolicyPort {
	return NetworkPolicyPort{
		Protocol: ConvertProtocolPointer(npp.Protocol),
		Port:     npp.Port, // Keep IntOrString
		EndPort:  npp.EndPort,
	}
}

func ConvertNetworkPolicyPorts(npps []networkingv1.NetworkPolicyPort) []NetworkPolicyPort {
	if npps == nil {
		return nil
	}
	res := make([]NetworkPolicyPort, len(npps))
	for i, p := range npps {
		res[i] = ConvertNetworkPolicyPort(p)
	}
	return res
}

func ConvertIPBlock(ipb *networkingv1.IPBlock) *IPBlock {
	if ipb == nil {
		return nil
	}
	return &IPBlock{
		CIDR:   ipb.CIDR,
		Except: ipb.Except,
	}
}

func ConvertNetworkPolicyPeer(npp networkingv1.NetworkPolicyPeer) NetworkPolicyPeer {
	return NetworkPolicyPeer{
		PodSelector:       ConvertLabelSelector(npp.PodSelector),
		NamespaceSelector: ConvertLabelSelector(npp.NamespaceSelector),
		IPBlock:           ConvertIPBlock(npp.IPBlock),
	}
}

func ConvertNetworkPolicyPeers(npps []networkingv1.NetworkPolicyPeer) []NetworkPolicyPeer {
	if npps == nil {
		return nil
	}
	res := make([]NetworkPolicyPeer, len(npps))
	for i, p := range npps {
		res[i] = ConvertNetworkPolicyPeer(p)
	}
	return res
}

func ConvertNetworkPolicyIngressRule(npir networkingv1.NetworkPolicyIngressRule) NetworkPolicyIngressRule {
	return NetworkPolicyIngressRule{
		Ports: ConvertNetworkPolicyPorts(npir.Ports),
		From:  ConvertNetworkPolicyPeers(npir.From),
	}
}

func ConvertNetworkPolicyIngressRules(npirs []networkingv1.NetworkPolicyIngressRule) []NetworkPolicyIngressRule {
	if npirs == nil {
		return nil
	}
	res := make([]NetworkPolicyIngressRule, len(npirs))
	for i, r := range npirs {
		res[i] = ConvertNetworkPolicyIngressRule(r)
	}
	return res
}

func ConvertNetworkPolicyEgressRule(nper networkingv1.NetworkPolicyEgressRule) NetworkPolicyEgressRule {
	return NetworkPolicyEgressRule{
		Ports: ConvertNetworkPolicyPorts(nper.Ports),
		To:    ConvertNetworkPolicyPeers(nper.To),
	}
}

func ConvertNetworkPolicyEgressRules(npers []networkingv1.NetworkPolicyEgressRule) []NetworkPolicyEgressRule {
	if npers == nil {
		return nil
	}
	res := make([]NetworkPolicyEgressRule, len(npers))
	for i, r := range npers {
		res[i] = ConvertNetworkPolicyEgressRule(r)
	}
	return res
}

func ConvertNetworkPolicySpec(spec networkingv1.NetworkPolicySpec) NetworkPolicySpec {
	return NetworkPolicySpec{
		PodSelector: *ConvertLabelSelector(&spec.PodSelector), // Selector is not pointer
		PolicyTypes: ConvertPolicyTypes(spec.PolicyTypes),
		Ingress:     ConvertNetworkPolicyIngressRules(spec.Ingress),
		Egress:      ConvertNetworkPolicyEgressRules(spec.Egress),
	}
}

func ConvertNetworkPolicy(np *networkingv1.NetworkPolicy) NetworkPolicy {
	return NetworkPolicy{
		TypeMeta:   ConvertTypeMeta(np.TypeMeta),
		ObjectMeta: ConvertObjectMeta(&np.ObjectMeta),
		Spec:       ConvertNetworkPolicySpec(np.Spec),
	}
}

// --- PodDisruptionBudget (policyv1) ---
type PodDisruptionBudget struct {
	TypeMeta
	ObjectMeta
	Spec   PodDisruptionBudgetSpec
	Status PodDisruptionBudgetStatus
}

type PodDisruptionBudgetSpec struct {
	MinAvailable               *intstr.IntOrString // Keep IntOrString
	Selector                   *LabelSelector      // Use helpers.LabelSelector
	MaxUnavailable             *intstr.IntOrString // Keep IntOrString
	UnhealthyPodEvictionPolicy *string             // policyv1.UnhealthyPodEvictionPolicyType
}

type PodDisruptionBudgetStatus struct {
	ObservedGeneration int64
	DisruptedPods      map[string]time.Time // Use helpers.Time via ConvertTime
	DisruptionsAllowed int32
	CurrentHealthy     int32
	DesiredHealthy     int32
	ExpectedPods       int32
	Conditions         []metav1.Condition // Use metav1.Condition directly
}

func ConvertUnhealthyPodEvictionPolicyType(upet *policyv1.UnhealthyPodEvictionPolicyType) *string {
	if upet == nil {
		return nil
	}
	s := string(*upet)
	return &s
}

func ConvertPodDisruptionBudgetSpec(spec policyv1.PodDisruptionBudgetSpec) PodDisruptionBudgetSpec {
	return PodDisruptionBudgetSpec{
		MinAvailable:               spec.MinAvailable, // Keep IntOrString
		Selector:                   ConvertLabelSelector(spec.Selector),
		MaxUnavailable:             spec.MaxUnavailable, // Keep IntOrString
		UnhealthyPodEvictionPolicy: ConvertUnhealthyPodEvictionPolicyType(spec.UnhealthyPodEvictionPolicy),
	}
}

// Convert metav1.Time map to time.Time map
func ConvertDisruptedPodsMap(pods map[string]metav1.Time) map[string]time.Time {
	if pods == nil {
		return nil
	}
	res := make(map[string]time.Time)
	for k, v := range pods {
		res[k] = ConvertTime(v)
	}
	return res
}

func ConvertPodDisruptionBudgetStatus(status policyv1.PodDisruptionBudgetStatus) PodDisruptionBudgetStatus {
	return PodDisruptionBudgetStatus{
		ObservedGeneration: status.ObservedGeneration,
		DisruptedPods:      ConvertDisruptedPodsMap(status.DisruptedPods),
		DisruptionsAllowed: status.DisruptionsAllowed,
		CurrentHealthy:     status.CurrentHealthy,
		DesiredHealthy:     status.DesiredHealthy,
		ExpectedPods:       status.ExpectedPods,
		Conditions:         status.Conditions, // Use metav1.Condition directly
	}
}

func ConvertPodDisruptionBudget(pdb *policyv1.PodDisruptionBudget) PodDisruptionBudget {
	return PodDisruptionBudget{
		TypeMeta:   ConvertTypeMeta(pdb.TypeMeta),
		ObjectMeta: ConvertObjectMeta(&pdb.ObjectMeta),
		Spec:       ConvertPodDisruptionBudgetSpec(pdb.Spec),
		Status:     ConvertPodDisruptionBudgetStatus(pdb.Status),
	}
}

// --- HorizontalPodAutoscaler (autoscalingv1) ---
type HorizontalPodAutoscaler struct {
	TypeMeta
	ObjectMeta
	Spec   HorizontalPodAutoscalerSpec
	Status HorizontalPodAutoscalerStatus
}

type HorizontalPodAutoscalerSpec struct {
	ScaleTargetRef                 CrossVersionObjectReference
	MinReplicas                    *int32
	MaxReplicas                    int32
	TargetCPUUtilizationPercentage *int32 // Deprecated in v2
}

type CrossVersionObjectReference struct {
	Kind       string
	Name       string
	APIVersion string
}

type HorizontalPodAutoscalerStatus struct {
	ObservedGeneration              *int64
	LastScaleTime                   *time.Time // Use helpers.Time via ConvertTimePtr
	CurrentReplicas                 int32
	DesiredReplicas                 int32
	CurrentCPUUtilizationPercentage *int32 // Deprecated in v2
}

func ConvertCrossVersionObjectReference(ref autoscalingv1.CrossVersionObjectReference) CrossVersionObjectReference {
	return CrossVersionObjectReference{
		Kind:       ref.Kind,
		Name:       ref.Name,
		APIVersion: ref.APIVersion,
	}
}

func ConvertHorizontalPodAutoscalerSpec(spec autoscalingv1.HorizontalPodAutoscalerSpec) HorizontalPodAutoscalerSpec {
	return HorizontalPodAutoscalerSpec{
		ScaleTargetRef:                 ConvertCrossVersionObjectReference(spec.ScaleTargetRef),
		MinReplicas:                    spec.MinReplicas,
		MaxReplicas:                    spec.MaxReplicas,
		TargetCPUUtilizationPercentage: spec.TargetCPUUtilizationPercentage,
	}
}

func ConvertHorizontalPodAutoscalerStatus(status autoscalingv1.HorizontalPodAutoscalerStatus) HorizontalPodAutoscalerStatus {
	return HorizontalPodAutoscalerStatus{
		ObservedGeneration:              status.ObservedGeneration,
		LastScaleTime:                   ConvertTimePtr(status.LastScaleTime),
		CurrentReplicas:                 status.CurrentReplicas,
		DesiredReplicas:                 status.DesiredReplicas,
		CurrentCPUUtilizationPercentage: status.CurrentCPUUtilizationPercentage,
	}
}

func ConvertHorizontalPodAutoscaler(hpa *autoscalingv1.HorizontalPodAutoscaler) HorizontalPodAutoscaler {
	return HorizontalPodAutoscaler{
		TypeMeta:   ConvertTypeMeta(hpa.TypeMeta),
		ObjectMeta: ConvertObjectMeta(&hpa.ObjectMeta),
		Spec:       ConvertHorizontalPodAutoscalerSpec(hpa.Spec),
		Status:     ConvertHorizontalPodAutoscalerStatus(hpa.Status),
	}
}

// --- CustomResourceDefinition (apiextensionsv1) ---
// NOTE: This is a very complex type. Mapping selectively.
type CustomResourceDefinition struct {
	TypeMeta
	ObjectMeta
	Spec   CustomResourceDefinitionSpec
	Status CustomResourceDefinitionStatus
}

type CustomResourceDefinitionSpec struct {
	Group                 string
	Names                 CustomResourceDefinitionNames
	Scope                 string // apiextensionsv1.ResourceScope
	Versions              []CustomResourceDefinitionVersion
	Conversion            *CustomResourceConversion
	PreserveUnknownFields bool
}

type CustomResourceDefinitionNames struct {
	Plural     string
	Singular   string
	ShortNames []string
	Kind       string
	ListKind   string
	Categories []string
}

type CustomResourceDefinitionVersion struct {
	Name                     string
	Served                   bool
	Storage                  bool
	Deprecated               bool
	DeprecationWarning       *string
	Schema                   *CustomResourceValidation
	Subresources             *CustomResourceSubresources
	AdditionalPrinterColumns []CustomResourceColumnDefinition
	SelectableFields         []SelectableField // Added later
}

type CustomResourceValidation struct {
	OpenAPIV3Schema *JSONSchemaProps // Represents JSONSchemaProps from apiextensionsv1
}

// JSONSchemaProps is extremely complex, providing a minimal placeholder
type JSONSchemaProps struct {
	// Add fields as needed, e.g., Type, Format, Properties, Items etc.
	Type string
	// ... other fields
}

type CustomResourceSubresources struct {
	Status *CustomResourceSubresourceStatus
	Scale  *CustomResourceSubresourceScale
}

type CustomResourceSubresourceStatus struct{}

type CustomResourceSubresourceScale struct {
	SpecReplicasPath   string
	StatusReplicasPath string
	LabelSelectorPath  *string
}

type CustomResourceColumnDefinition struct {
	Name        string
	Type        string
	Format      string
	Description string
	Priority    int32
	JSONPath    string
}

type CustomResourceConversion struct {
	Strategy string // apiextensionsv1.ConversionStrategyType
	Webhook  *WebhookConversion
}

type WebhookConversion struct {
	ClientConfig             *WebhookClientConfig
	ConversionReviewVersions []string
}

type WebhookClientConfig struct {
	URL      *string // Pointer in v1
	Service  *ServiceReference
	CABundle []byte
}

type ServiceReference struct {
	Namespace string
	Name      string
	Path      *string
	Port      *int32
}

type CustomResourceDefinitionStatus struct {
	Conditions     []CustomResourceDefinitionCondition
	AcceptedNames  CustomResourceDefinitionNames
	StoredVersions []string
}

type CustomResourceDefinitionCondition struct {
	Type               string     // apiextensionsv1.CustomResourceDefinitionConditionType
	Status             string     // apiextensionsv1.ConditionStatus
	LastTransitionTime *time.Time // Use helpers.Time via ConvertTimePtr
	Reason             string
	Message            string
}

type SelectableField struct {
	JSONPath string
}

// NOTE: Conversion functions for CRD are highly complex and omitted for brevity.
// Add specific conversions as needed.
func ConvertCustomResourceDefinitionNames(names apiextensionsv1.CustomResourceDefinitionNames) CustomResourceDefinitionNames {
	return CustomResourceDefinitionNames{ /* ... map fields ... */ }
}

// ... and many more conversion functions ...

func ConvertCustomResourceDefinition(crd *apiextensionsv1.CustomResourceDefinition) CustomResourceDefinition {
	// Simplified - requires full implementation of sub-conversions
	return CustomResourceDefinition{
		TypeMeta:   ConvertTypeMeta(crd.TypeMeta),
		ObjectMeta: ConvertObjectMeta(&crd.ObjectMeta),
		// Spec:    ConvertCustomResourceDefinitionSpec(crd.Spec), // Requires full conversion chain
		// Status:  ConvertCustomResourceDefinitionStatus(crd.Status),
	}
}

// --- EndpointSlice (discoveryv1) ---
type EndpointSlice struct {
	TypeMeta
	ObjectMeta
	AddressType string         // discoveryv1.AddressType
	Endpoints   []Endpoint     // Use helpers.Endpoint
	Ports       []EndpointPort // Use helpers.EndpointPort
}

type Endpoint struct {
	Addresses          []string
	Conditions         EndpointConditions
	Hostname           *string
	TargetRef          *ObjectReference  // Use helpers.ObjectReference
	Topology           map[string]string // Deprecated
	NodeName           *string           // Added later
	Zone               *string           // Added later
	Hints              *EndpointHints    // Added later
	DeprecatedTopology map[string]string
}

type EndpointConditions struct {
	Ready       *bool
	Serving     *bool
	Terminating *bool
}

type EndpointHints struct {
	ForZones []ForZone
}

type ForZone struct {
	Name string
}

// Convert EndpointPort is likely already defined for corev1.Endpoints

func ConvertAddressType(at discoveryv1.AddressType) string {
	return string(at)
}

func ConvertEndpointConditions(ec discoveryv1.EndpointConditions) EndpointConditions {
	return EndpointConditions{
		Ready:       ec.Ready,
		Serving:     ec.Serving,
		Terminating: ec.Terminating,
	}
}

func ConvertForZone(fz discoveryv1.ForZone) ForZone {
	return ForZone{Name: fz.Name}
}

func ConvertForZones(fzs []discoveryv1.ForZone) []ForZone {
	if fzs == nil {
		return nil
	}
	res := make([]ForZone, len(fzs))
	for i, fz := range fzs {
		res[i] = ConvertForZone(fz)
	}
	return res
}

func ConvertEndpointHints(eh *discoveryv1.EndpointHints) *EndpointHints {
	if eh == nil {
		return nil
	}
	return &EndpointHints{ForZones: ConvertForZones(eh.ForZones)}
}

func ConvertEndpoint(ep discoveryv1.Endpoint) Endpoint {
	return Endpoint{
		Addresses:          ep.Addresses,
		Conditions:         ConvertEndpointConditions(ep.Conditions),
		Hostname:           ep.Hostname,
		TargetRef:          ConvertObjectReferencePointer(ep.TargetRef),
		NodeName:           ep.NodeName,
		Zone:               ep.Zone,
		Hints:              ConvertEndpointHints(ep.Hints),
		DeprecatedTopology: ep.DeprecatedTopology,
	}
}

func ConvertEndpointsSliceEndpoints(eps []discoveryv1.Endpoint) []Endpoint {
	if eps == nil {
		return nil
	}
	res := make([]Endpoint, len(eps))
	for i, ep := range eps {
		res[i] = ConvertEndpoint(ep)
	}
	return res
}

func ConvertEndpointPortsV2(eps []discoveryv1.EndpointPort) []EndpointPort {
	if eps == nil {
		return nil
	}

	res := make([]EndpointPort, len(eps))
	for i, ep := range eps {
		var name string
		var port int32
		var protocol string
		if ep.Name != nil {
			name = *ep.Name
		}
		if ep.Port != nil {
			port = *ep.Port
		}
		if ep.Protocol != nil {
			protocol = string(*ep.Protocol)
		}
		res[i] = EndpointPort{
			Name:        name,
			Port:        port,
			Protocol:    protocol,
			AppProtocol: ep.AppProtocol,
		}
	}
	return res
}

func ConvertEndpointSlice(es *discoveryv1.EndpointSlice) EndpointSlice {
	return EndpointSlice{
		TypeMeta:    ConvertTypeMeta(es.TypeMeta),
		ObjectMeta:  ConvertObjectMeta(&es.ObjectMeta),
		AddressType: ConvertAddressType(es.AddressType),
		Endpoints:   ConvertEndpointsSliceEndpoints(es.Endpoints),
		Ports:       ConvertEndpointPortsV2(es.Ports), // Assumes converter exists
	}
}

// Add necessary conversion functions for PVC Status
func ConvertPersistentVolumeClaim(pvc *corev1.PersistentVolumeClaim) PersistentVolumeClaim {
	if pvc == nil {
		// Return zero value or handle error as appropriate
		return PersistentVolumeClaim{}
	}
	return PersistentVolumeClaim{
		TypeMeta:   ConvertTypeMeta(pvc.TypeMeta),
		ObjectMeta: ConvertObjectMeta(&pvc.ObjectMeta),
		Spec:       ConvertPersistentVolumeClaimSpec(pvc.Spec),
		Status:     ConvertPersistentVolumeClaimStatus(pvc.Status),
	}
}

// Added IngressLoadBalancerStatus helper type
type IngressLoadBalancerStatus struct {
	Ingress []LoadBalancerIngress // Use helpers.LoadBalancerIngress
}

// ConvertIngressBackend converts *networkingv1.IngressBackend to *IngressBackend
// Assumes IngressBackend, ConvertIngressServiceBackend and ConvertTypedLocalObjectReference exist.
func ConvertIngressBackend(ib *networkingv1.IngressBackend) *IngressBackend {
	if ib == nil {
		return nil
	}
	return &IngressBackend{
		Service:  ConvertIngressServiceBackend(ib.Service),
		Resource: ConvertTypedLocalObjectReference(ib.Resource),
	}
}

// ConvertIngressTLS converts networkingv1.IngressTLS to IngressTLS
// Assumes IngressTLS helper type exists
func ConvertIngressTLS(itls networkingv1.IngressTLS) IngressTLS {
	return IngressTLS{
		Hosts:      itls.Hosts,
		SecretName: itls.SecretName,
	}
}

// ConvertIngressTLSs converts []networkingv1.IngressTLS to []IngressTLS
// Assumes IngressTLS helper type and ConvertIngressTLS exist.
func ConvertIngressTLSs(itls []networkingv1.IngressTLS) []IngressTLS {
	if itls == nil {
		return nil
	}
	res := make([]IngressTLS, len(itls))
	for i, tls := range itls {
		res[i] = ConvertIngressTLS(tls) // Assumes ConvertIngressTLS exists
	}
	return res
}

// ConvertIngressRule converts networkingv1.IngressRule to IngressRule
func ConvertIngressRule(ir networkingv1.IngressRule) IngressRule {
	return IngressRule{
		Host:             ir.Host,
		IngressRuleValue: ConvertIngressRuleValue(ir.IngressRuleValue),
	}
}

// ConvertIngressRules converts []networkingv1.IngressRule to []IngressRule
// Assumes IngressRule helper type and ConvertIngressRule exist.
func ConvertIngressRules(irs []networkingv1.IngressRule) []IngressRule {
	if irs == nil {
		return nil
	}
	res := make([]IngressRule, len(irs))
	for i, r := range irs {
		res[i] = ConvertIngressRule(r) // Assumes ConvertIngressRule exists
	}
	return res
}
