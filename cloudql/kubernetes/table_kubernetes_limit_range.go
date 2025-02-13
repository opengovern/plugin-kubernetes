package kubernetes

import (
	"context"
	opengovernance "github.com/opengovern/og-describer-kubernetes/discovery/pkg/es"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableKubernetesLimitRange(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "kubernetes_limit_range",
		Description: "Kubernetes Limit Range",
		Get: &plugin.GetConfig{
			Hydrate: opengovernance.GetKubernetesLimitRange,
		},
		List: &plugin.ListConfig{
			Hydrate: opengovernance.ListKubernetesLimitRange,
		},
		Columns: commonColumns([]*plugin.Column{
			{
				Name:        "spec_limits",
				Type:        proto.ColumnType_JSON,
				Description: "List of limit range item objects that are enforced.",
				Transform:   transform.FromField("Description.LimitRange.Spec.Limits"),
			},
			{
				Name:        "title",
				Type:        proto.ColumnType_STRING,
				Description: ColumnDescriptionTitle,
				Transform:   transform.FromField("Description.LimitRange.Name"),
			},
			{
				Name:        "tags",
				Type:        proto.ColumnType_JSON,
				Description: ColumnDescriptionTags,
				Transform:   transform.From(transformLimitRangeTags),
			},
		}),
	}
}

func transformLimitRangeTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	obj := d.HydrateItem.(opengovernance.KubernetesLimitRange).Description.LimitRange
	return mergeTags(obj.Labels, obj.Annotations), nil
}
