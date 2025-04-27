package kubernetes

import (
	"context"
	opengovernance "github.com/opengovern/og-describer-kubernetes/discovery/pkg/es"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableKubernetesResource(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "k8_resource",
		Description: "ClusterRole contains rules that represent a set of permissions.",
		List: &plugin.ListConfig{
			Hydrate: opengovernance.ListKubernetesResource,
		},
		Columns: commonGeneralColumns([]*plugin.Column{
			{
				Name:        "kind",
				Type:        proto.ColumnType_STRING,
				Description: "Resource Kind",
				Transform:   transform.FromField("Description.Kind"),
			},
			{
				Name:        "object_name",
				Type:        proto.ColumnType_STRING,
				Description: "object name.",
				Transform:   transform.FromField("Description.ObjectName"),
			},
			{
				Name:        "namespace",
				Type:        proto.ColumnType_STRING,
				Description: "namespace.",
				Transform:   transform.FromField("Description.Namespace"),
			},
			{
				Name:        "uid",
				Type:        proto.ColumnType_STRING,
				Description: "uid.",
				Transform:   transform.FromField("Description.UID"),
			},
			{
				Name:        "creation_timestamp",
				Type:        proto.ColumnType_STRING,
				Description: "creation timestamp.",
				Transform:   transform.FromField("Description.CreationTimestamp"),
			},
			{
				Name:        "resource_version",
				Type:        proto.ColumnType_STRING,
				Description: "resource version.",
				Transform:   transform.FromField("Description.ResourceVersion"),
			},
			{
				Name:        "resource_table",
				Type:        proto.ColumnType_STRING,
				Description: "creation timestamp.",
				Transform:   transform.FromField("Description.ResourceTable"),
			},
			{
				Name:        "api_version",
				Type:        proto.ColumnType_STRING,
				Description: "creation timestamp.",
				Transform:   transform.FromField("Description.ApiVersion"),
			},
		}),
	}
}
