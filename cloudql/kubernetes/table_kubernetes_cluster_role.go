package kubernetes

import (
	"context"
	opengovernance "github.com/opengovern/og-describer-kubernetes/discovery/pkg/es"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableKubernetesClusterRole(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "kubernetes_cluster_role",
		Description: "ClusterRole contains rules that represent a set of permissions.",
		Get: &plugin.GetConfig{
			Hydrate: opengovernance.GetKubernetesClusterRole,
		},
		List: &plugin.ListConfig{
			Hydrate: opengovernance.ListKubernetesClusterRole,
		},
		// ClusterRole, is a non-namespaced resource.
		Columns: commonColumns([]*plugin.Column{
			{
				Name:        "rules",
				Type:        proto.ColumnType_JSON,
				Description: "List of the PolicyRules for this Role.",
				Transform:   transform.FromField("Description.ClusterRole.Rules"),
			},
			{
				Name:        "aggregation_rule",
				Type:        proto.ColumnType_JSON,
				Description: "An optional field that describes how to build the Rules for this ClusterRole",
				Transform:   transform.FromField("Description.ClusterRole.AggregationRule"),
			},

			//// Steampipe Standard Columns
			{
				Name:        "title",
				Type:        proto.ColumnType_STRING,
				Description: ColumnDescriptionTitle,
				Transform:   transform.FromField("Description.ClusterRole.Name"),
			},
			{
				Name:        "tags",
				Type:        proto.ColumnType_JSON,
				Description: ColumnDescriptionTags,
				Transform:   transform.From(transformClusterRoleTags),
			},
		}),
	}
}

func transformClusterRoleTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	obj := d.HydrateItem.(opengovernance.KubernetesClusterRole).Description.ClusterRole
	return mergeTags(obj.Labels, obj.Annotations), nil
}
