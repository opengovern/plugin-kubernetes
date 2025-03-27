package kubernetes

import (
	"context"
	opengovernance "github.com/opengovern/og-describer-kubernetes/discovery/pkg/es"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableKubernetesPDB(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "kubernetes_pod_disruption_budget",
		Description: "A Pod Disruption Budget limits the number of Pods of a replicated application that are down simultaneously from voluntary disruptions.",
		List: &plugin.ListConfig{
			Hydrate: opengovernance.ListKubernetesPodDisruptionBudget,
		},
		Columns: commonColumns([]*plugin.Column{
			{
				Name:        "min_available",
				Type:        proto.ColumnType_STRING,
				Description: "An eviction is allowed if at least 'minAvailable' pods selected by 'selector' will still be available after the eviction.",
				Transform:   transform.FromField("Description.PodDisruptionBudget.Spec.MinAvailable"),
			},
			{
				Name:        "selector",
				Type:        proto.ColumnType_JSON,
				Description: "Label query over pods whose evictions are managed by the disruption budget.",
				Transform:   transform.FromField("Description.PodDisruptionBudget.Spec.Selector"),
			},
			{
				Name:        "max_unavailable",
				Type:        proto.ColumnType_STRING,
				Description: "An eviction is allowed if at most 'maxAvailable' pods selected by 'selector' will still be unavailable after the eviction.",
				Transform:   transform.FromField("Description.PodDisruptionBudget.Spec.MaxUnavailable"),
			},
			{
				Name:        "title",
				Type:        proto.ColumnType_STRING,
				Description: ColumnDescriptionTitle,
				Transform:   transform.FromField("Description.PodDisruptionBudget.Name"),
			},
			{
				Name:        "tags",
				Type:        proto.ColumnType_JSON,
				Description: ColumnDescriptionTags,
				Transform:   transform.From(transformPDBTags),
			},
		}),
	}
}

func transformPDBTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	obj := d.HydrateItem.(opengovernance.KubernetesPodDisruptionBudget).Description.PodDisruptionBudget
	return mergeTags(obj.Labels, obj.Annotations), nil
}
