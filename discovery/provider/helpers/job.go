package helpers

import (
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1" // Needed for PodFailurePolicyOnPodConditionsPattern status
)

// --- PodFailurePolicy ---
type PodFailurePolicyOnExitCodesRequirement struct {
	ContainerName *string
	Operator      string // PodFailurePolicyOnExitCodesOperator
	Values        []int32
}

func ConvertPodFailurePolicyOnExitCodesOperator(op batchv1.PodFailurePolicyOnExitCodesOperator) string {
	return string(op)
}

func ConvertPodFailurePolicyOnExitCodesRequirement(r *batchv1.PodFailurePolicyOnExitCodesRequirement) *PodFailurePolicyOnExitCodesRequirement {
	if r == nil {
		return nil
	}
	return &PodFailurePolicyOnExitCodesRequirement{
		ContainerName: r.ContainerName,
		Operator:      ConvertPodFailurePolicyOnExitCodesOperator(r.Operator),
		Values:        r.Values,
	}
}

type PodFailurePolicyOnPodConditionsPattern struct {
	Type   string // corev1.PodConditionType
	Status string // corev1.ConditionStatus
}

func ConvertPodConditionTypeToString(pct corev1.PodConditionType) string {
	return string(pct)
}

func ConvertConditionStatusToString(cs corev1.ConditionStatus) string {
	return string(cs)
}

// Renamed from ConvertPodFailurePolicyOnPodConditionsPattern to ConvertPodFailurePolicyOnPodConditionsPatternSingle
func ConvertPodFailurePolicyOnPodConditionsPatternSingle(p batchv1.PodFailurePolicyOnPodConditionsPattern) PodFailurePolicyOnPodConditionsPattern {
	return PodFailurePolicyOnPodConditionsPattern{
		Type:   ConvertPodConditionTypeToString(p.Type),
		Status: ConvertConditionStatusToString(p.Status),
	}
}

// Renamed from ConvertPodFailurePolicyOnPodConditionsPatterns
func ConvertPodFailurePolicyOnPodConditionsPatterns(pp []batchv1.PodFailurePolicyOnPodConditionsPattern) []PodFailurePolicyOnPodConditionsPattern {
	if pp == nil {
		return nil
	}
	ps := make([]PodFailurePolicyOnPodConditionsPattern, len(pp))
	for i, p := range pp {
		ps[i] = ConvertPodFailurePolicyOnPodConditionsPatternSingle(p) // Call the single converter
	}
	return ps
}

type PodFailurePolicyRule struct {
	Action          string // PodFailurePolicyAction
	OnExitCodes     *PodFailurePolicyOnExitCodesRequirement
	OnPodConditions []PodFailurePolicyOnPodConditionsPattern
}

func ConvertPodFailurePolicyAction(a batchv1.PodFailurePolicyAction) string {
	return string(a)
}

func ConvertPodFailurePolicyRule(r batchv1.PodFailurePolicyRule) PodFailurePolicyRule {
	return PodFailurePolicyRule{
		Action:          ConvertPodFailurePolicyAction(r.Action),
		OnExitCodes:     ConvertPodFailurePolicyOnExitCodesRequirement(r.OnExitCodes),
		OnPodConditions: ConvertPodFailurePolicyOnPodConditionsPatterns(r.OnPodConditions),
	}
}

type PodFailurePolicy struct {
	Rules []PodFailurePolicyRule
}

func ConvertPodFailurePolicy(policy *batchv1.PodFailurePolicy) *PodFailurePolicy {
	if policy == nil {
		return nil
	}
	rules := make([]PodFailurePolicyRule, len(policy.Rules))
	for i, r := range policy.Rules {
		rules[i] = ConvertPodFailurePolicyRule(r)
	}
	return &PodFailurePolicy{
		Rules: rules,
	}
}

// --- SuccessPolicy ---
type SuccessPolicyRule struct {
	SucceededIndexes *string
	SucceededCount   *int32
}

// ConvertSuccessPolicyRuleSingle converts a single rule
func ConvertSuccessPolicyRuleSingle(r batchv1.SuccessPolicyRule) SuccessPolicyRule {
	return SuccessPolicyRule{
		SucceededIndexes: r.SucceededIndexes,
		SucceededCount:   r.SucceededCount,
	}
}

// ConvertSuccessPolicyRules converts a slice of rules
func ConvertSuccessPolicyRules(srcRules []batchv1.SuccessPolicyRule) []SuccessPolicyRule {
	if srcRules == nil {
		return nil
	}
	rules := make([]SuccessPolicyRule, len(srcRules))
	for i, r := range srcRules {
		rules[i] = ConvertSuccessPolicyRuleSingle(r)
	}
	return rules
}

type SuccessPolicy struct {
	Rules []SuccessPolicyRule
}

func ConvertSuccessPolicy(p *batchv1.SuccessPolicy) *SuccessPolicy {
	if p == nil {
		return nil
	}
	return &SuccessPolicy{
		Rules: ConvertSuccessPolicyRules(p.Rules),
	}
}

// --- Job Completion/Replacement Policies ---
type CompletionMode string       // batchv1.CompletionMode
type PodReplacementPolicy string // batchv1.PodReplacementPolicy

func ConvertCompletionMode(cm *batchv1.CompletionMode) *CompletionMode {
	if cm == nil {
		return nil
	}
	mode := CompletionMode(*cm)
	return &mode
}

func ConvertPodReplacementPolicy(prp *batchv1.PodReplacementPolicy) *PodReplacementPolicy {
	if prp == nil {
		return nil
	}
	policy := PodReplacementPolicy(*prp)
	return &policy
}

// --- JobSpec ---
type JobSpec struct {
	Parallelism             *int32
	Completions             *int32
	ActiveDeadlineSeconds   *int64
	PodFailurePolicy        *PodFailurePolicy
	SuccessPolicy           *SuccessPolicy
	BackoffLimit            *int32
	BackoffLimitPerIndex    *int32
	MaxFailedIndexes        *int32
	Selector                *LabelSelector // Assumes LabelSelector in model_helpers.go
	ManualSelector          *bool
	Template                PodTemplateSpec // Assumes PodTemplateSpec in pod.go
	TTLSecondsAfterFinished *int32
	CompletionMode          *CompletionMode
	Suspend                 *bool
	PodReplacementPolicy    *PodReplacementPolicy
	ManagedBy               *string
}

func ConvertJobSpec(spec batchv1.JobSpec) JobSpec {
	return JobSpec{
		Parallelism:             spec.Parallelism,
		Completions:             spec.Completions,
		ActiveDeadlineSeconds:   spec.ActiveDeadlineSeconds,
		PodFailurePolicy:        ConvertPodFailurePolicy(spec.PodFailurePolicy),
		SuccessPolicy:           ConvertSuccessPolicy(spec.SuccessPolicy),
		BackoffLimit:            spec.BackoffLimit,
		BackoffLimitPerIndex:    spec.BackoffLimitPerIndex,
		MaxFailedIndexes:        spec.MaxFailedIndexes,
		Selector:                ConvertLabelSelector(spec.Selector), // Assumes ConvertLabelSelector in model_helpers.go
		ManualSelector:          spec.ManualSelector,
		Template:                ConvertPodTemplateSpec(spec.Template), // Assumes ConvertPodTemplateSpec in pod.go
		TTLSecondsAfterFinished: spec.TTLSecondsAfterFinished,
		CompletionMode:          ConvertCompletionMode(spec.CompletionMode),
		Suspend:                 spec.Suspend,
		PodReplacementPolicy:    ConvertPodReplacementPolicy(spec.PodReplacementPolicy),
		ManagedBy:               spec.ManagedBy,
	}
}
