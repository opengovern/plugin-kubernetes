package kubernetes

import (
	"context"
	opengovernance "github.com/opengovern/og-describer-kubernetes/discovery/pkg/es"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableKubernetesEndpointSlice(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "kubernetes_endpoint_slice",
		Description: "EndpointSlice represents a subset of the endpoints that implement a service.",
		Get: &plugin.GetConfig{
			Hydrate: opengovernance.GetKubernetesEndpointSlice,
		},
		List: &plugin.ListConfig{
			Hydrate: opengovernance.ListKubernetesEndpointSlice,
		},
		Columns: commonColumns([]*plugin.Column{
			{
				Name:        "address_type",
				Type:        proto.ColumnType_STRING,
				Description: "Type of address carried by this EndpointSlice. All addresses in the slice are of the same type. Supported types are IPv4, IPv6, and FQDN.",
				Transform:   transform.FromField("Description.EndpointSlice.AddressType"),
			},
			{
				Name:        "endpoints",
				Type:        proto.ColumnType_JSON,
				Description: "List of unique endpoints in this slice.",
				Transform:   transform.FromField("Description.EndpointSlice.Endpoints"),
			},
			{
				Name:        "ports",
				Type:        proto.ColumnType_JSON,
				Description: "List of network ports exposed by each endpoint in this slice. Each port must have a unique name. When ports is empty, it indicates that there are no defined ports.",
				Transform:   transform.FromField("Description.EndpointSlice.Ports"),
			},
			//// Steampipe Standard Columns
			{
				Name:        "title",
				Type:        proto.ColumnType_STRING,
				Description: ColumnDescriptionTitle,
				Transform:   transform.FromField("Description.EndpointSlice.Name"),
			},
			{
				Name:        "tags",
				Type:        proto.ColumnType_JSON,
				Description: ColumnDescriptionTags,
				Transform:   transform.From(transformEndpointSliceTags),
			},
		}),
	}
}

func transformEndpointSliceTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	obj := d.HydrateItem.(opengovernance.KubernetesEndpointSlice).Description.EndpointSlice
	return mergeTags(obj.Labels, obj.Annotations), nil
}
