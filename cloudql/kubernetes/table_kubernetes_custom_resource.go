package kubernetes

import (
	"context"
	opengovernance "github.com/opengovern/og-describer-kubernetes/discovery/pkg/es"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableKubernetesCustomResource(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "kubernetes_custom_resource",
		Description: "Custom resources are extensions of the Kubernetes API.",
		Get: &plugin.GetConfig{
			Hydrate: opengovernance.GetKubernetesCustomResource,
		},
		List: &plugin.ListConfig{
			Hydrate: opengovernance.ListKubernetesCustomResource,
		},
		Columns: commonColumns([]*plugin.Column{
			{
				Name:        "fully_qualified_name",
				Type:        proto.ColumnType_STRING,
				Description: "The fully qualified name of the custom resource.",
				Transform:   transform.FromField("Description.FullyQualifiedName"),
			},
			{
				Name:        "title",
				Type:        proto.ColumnType_STRING,
				Description: ColumnDescriptionTitle,
				Transform:   transform.FromField("Description.MetaObject.Name"),
			},
			{
				Name:        "tags",
				Type:        proto.ColumnType_JSON,
				Description: ColumnDescriptionTags,
				Transform:   transform.From(transformCustomReousrceTags),
			},
		}),
	}
}

func transformCustomReousrceTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	obj := d.HydrateItem.(opengovernance.KubernetesCustomResource).Description.MetaObject
	return mergeTags(obj.Labels, obj.Annotations), nil
}
