package kubernetes

import (
	"context"
	opengovernance "github.com/opengovern/og-describer-kubernetes/discovery/pkg/es"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableKubernetesNetworkPolicy(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "kubernetes_network_policy",
		Description: "Network policy specifiy how pods are allowed to communicate with each other and with other network endpoints.",
		List: &plugin.ListConfig{
			Hydrate: opengovernance.ListKubernetesNetworkPolicy,
		},
		Columns: commonColumns([]*plugin.Column{
			{
				Name:        "pod_selector",
				Type:        proto.ColumnType_JSON,
				Description: "Selects the pods to which this NetworkPolicy object applies. The array of ingress rules is applied to any pods selected by this field. An empty podSelector matches all pods in this namespace.",
				Transform:   transform.FromField("Description.NetworkPolicy.Spec.PodSelector"),
			},
			{
				Name:        "ingress",
				Type:        proto.ColumnType_JSON,
				Description: "List of ingress rules to be applied to the selected pods. If this field is empty then this NetworkPolicy does not allow any traffic (and serves solely to ensure that the pods it selects are isolated by default)",
				Transform:   transform.FromField("Description.NetworkPolicy.Spec.Ingress"),
			},
			{
				Name:        "egress",
				Type:        proto.ColumnType_JSON,
				Description: "List of egress rules to be applied to the selected pods. If this field is empty then this NetworkPolicy limits all outgoing traffic (and serves solely to ensure that the pods it selects are isolated by default).",
				Transform:   transform.FromField("Description.NetworkPolicy.Spec.Egress"),
			},
			{
				Name:        "policy_types",
				Type:        proto.ColumnType_JSON,
				Description: "List of rule types that the NetworkPolicy relates to. Valid options are \"Ingress\", \"Egress\", or \"Ingress,Egress\". If this field is not specified, it will default based on the existence of Ingress or Egress rules.",
				Transform:   transform.FromField("Description.NetworkPolicy.Spec.PolicyTypes"),
			},
			{
				Name:        "title",
				Type:        proto.ColumnType_STRING,
				Description: ColumnDescriptionTitle,
				Transform:   transform.FromField("Description.NetworkPolicy.Name"),
			},
			{
				Name:        "tags",
				Type:        proto.ColumnType_JSON,
				Description: ColumnDescriptionTags,
				Transform:   transform.From(transformNetworkPolicyTags),
			},
		}),
	}
}

func transformNetworkPolicyTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	obj := d.HydrateItem.(opengovernance.KubernetesNetworkPolicy).Description.NetworkPolicy
	return mergeTags(obj.Labels, obj.Annotations), nil
}
