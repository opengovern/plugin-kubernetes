package kubernetes

import (
	"context"
	opengovernance "github.com/opengovern/og-describer-kubernetes/discovery/pkg/es"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableKubernetesSecret(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "k8_secret",
		Description: "Secrets can be used to store sensitive information either as individual properties or coarse-grained entries like entire files or JSON blobs.",
		List: &plugin.ListConfig{
			Hydrate: opengovernance.ListKubernetesSecret,
		},
		Columns: commonColumns([]*plugin.Column{
			{
				Name:        "immutable",
				Type:        proto.ColumnType_BOOL,
				Description: "If set to true, ensures that data stored in the Secret cannot be updated (only object metadata can be modified). If not set to true, the field can be modified at any time. Defaulted to nil.",
				Transform:   transform.FromField("Description.Secret.Immutable"),
			},
			{
				Name:        "type",
				Type:        proto.ColumnType_STRING,
				Description: "Type of the secret data.",
				Transform:   transform.FromField("Description.Secret.Type"),
			},
			{
				Name:        "data",
				Type:        proto.ColumnType_JSON,
				Description: "Type of the secret data.",
				Transform:   transform.FromField("Description.Secret.Data"),
			},
			{
				Name:        "string_data",
				Type:        proto.ColumnType_JSON,
				Description: "Type of the secret data.",
				Transform:   transform.FromField("Description.Secret.StringData"),
			},

			//// Steampipe Standard Columns
			{
				Name:        "title",
				Type:        proto.ColumnType_STRING,
				Description: ColumnDescriptionTitle,
				Transform:   transform.FromField("Description.Secret.Name"),
			},
			{
				Name:        "tags",
				Type:        proto.ColumnType_JSON,
				Description: ColumnDescriptionTags,
				Transform:   transform.From(transformSecretTags),
			},
		}),
	}
}

func transformSecretTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	obj := d.HydrateItem.(opengovernance.KubernetesSecret).Description.Secret
	return mergeTags(obj.Labels, obj.Annotations), nil
}
