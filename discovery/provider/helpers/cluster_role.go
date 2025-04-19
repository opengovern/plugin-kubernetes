package helpers

import (
	rbacv1 "k8s.io/api/rbac/v1"
)

// --- ClusterRole ---
type ClusterRole struct {
	TypeMeta
	ObjectMeta
	Rules           []PolicyRule
	AggregationRule *AggregationRule
}

// ConvertClusterRole creates a helper ClusterRole from a rbacv1 ClusterRole
func ConvertClusterRole(cr *rbacv1.ClusterRole) ClusterRole {
	return ClusterRole{
		TypeMeta:        ConvertTypeMeta(cr.TypeMeta),
		ObjectMeta:      ConvertObjectMeta(&cr.ObjectMeta),
		Rules:           ConvertPolicyRules(cr.Rules),
		AggregationRule: ConvertAggregationRule(cr.AggregationRule),
	}
}

// --- PolicyRule ---
type PolicyRule struct {
	Verbs           []string
	APIGroups       []string
	Resources       []string
	ResourceNames   []string
	NonResourceURLs []string
}

func ConvertPolicyRules(policyRules []rbacv1.PolicyRule) []PolicyRule {
	if policyRules == nil {
		return nil
	}
	rules := make([]PolicyRule, len(policyRules))
	for i, policyRule := range policyRules {
		rules[i] = PolicyRule{
			Verbs:           policyRule.Verbs,
			APIGroups:       policyRule.APIGroups,
			Resources:       policyRule.Resources,
			ResourceNames:   policyRule.ResourceNames,
			NonResourceURLs: policyRule.NonResourceURLs,
		}
	}
	return rules
}

// --- AggregationRule ---
type AggregationRule struct {
	ClusterRoleSelectors []LabelSelector // Assumes LabelSelector in model_helpers.go
}

func ConvertAggregationRule(aggregationRule *rbacv1.AggregationRule) *AggregationRule {
	if aggregationRule == nil {
		return nil
	}
	selectors := make([]LabelSelector, 0, len(aggregationRule.ClusterRoleSelectors))
	for _, sel := range aggregationRule.ClusterRoleSelectors {
		// Assuming ConvertLabelSelector exists in model_helpers.go
		convertedSel := ConvertLabelSelector(&sel)
		if convertedSel != nil {
			selectors = append(selectors, *convertedSel)
		}
	}
	return &AggregationRule{
		ClusterRoleSelectors: selectors,
	}
}
