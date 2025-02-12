package kubernetes

import (
	"context"
	opengovernance "github.com/opengovern/og-describer-kubernetes/discovery/pkg/es"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableKubernetesConfigMap(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "kubernetes_config_map",
		Description: "Config Map can be used to store fine-grained information like individual properties or coarse-grained information like entire config files or JSON blobs.",
		Get: &plugin.GetConfig{
			Hydrate: opengovernance.GetKubernetesConfigMap,
		},
		List: &plugin.ListConfig{
			Hydrate: opengovernance.ListKubernetesConfigMap,
		},
		// ClusterRole, is a non-namespaced resource.
		Columns: commonColumns([]*plugin.Column{
			{
				Name:        "immutable",
				Type:        proto.ColumnType_BOOL,
				Description: "If set to true, ensures that data stored in the ConfigMap cannot be updated (only object metadata can be modified). If not set to true, the field can be modified at any time. Defaulted to nil.",
				Transform:   transform.FromField("Description.ConfigMap.Immutable"),
			},
			//// Steampipe Standard Columns
			{
				Name:        "title",
				Type:        proto.ColumnType_STRING,
				Description: ColumnDescriptionTitle,
				Transform:   transform.FromField("Description.ConfigMap.Name"),
			},
			{
				Name:        "tags",
				Type:        proto.ColumnType_JSON,
				Description: ColumnDescriptionTags,
				Transform:   transform.From(transformConfigMapTags),
			},
		}),
	}
}

func transformConfigMapTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	obj := d.HydrateItem.(opengovernance.KubernetesConfigMap).Description.ConfigMap
	return mergeTags(obj.Labels, obj.Annotations), nil
}
