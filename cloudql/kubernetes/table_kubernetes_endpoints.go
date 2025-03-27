package kubernetes

import (
	"context"
	opengovernance "github.com/opengovern/og-describer-kubernetes/discovery/pkg/es"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableKubernetesEndpoints(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "kubernetes_endpoints",
		Description: "Set of addresses and ports that comprise a service. More info: https://kubernetes.io/docs/concepts/services-networking/service/#services-without-selectors.",
		List: &plugin.ListConfig{
			Hydrate: opengovernance.ListKubernetesEndpoint,
		},
		Columns: commonColumns([]*plugin.Column{
			{
				Name:        "subsets",
				Type:        proto.ColumnType_JSON,
				Description: "List of addresses and ports that comprise a service.",
				Transform:   transform.FromField("Description.Endpoint.Subsets"),
			},
			{
				Name:        "title",
				Type:        proto.ColumnType_STRING,
				Description: ColumnDescriptionTitle,
				Transform:   transform.FromField("Description.Endpoint.Name"),
			},
			{
				Name:        "tags",
				Type:        proto.ColumnType_JSON,
				Description: ColumnDescriptionTags,
				Transform:   transform.From(transformEndpointTags),
			},
		}),
	}
}

func transformEndpointTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	obj := d.HydrateItem.(opengovernance.KubernetesEndpoint).Description.Endpoint
	return mergeTags(obj.Labels, obj.Annotations), nil
}
