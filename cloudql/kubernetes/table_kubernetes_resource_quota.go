package kubernetes

import (
	"context"
	opengovernance "github.com/opengovern/og-describer-kubernetes/discovery/pkg/es"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableKubernetesResourceQuota(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "kubernetes_resource_quota",
		Description: "Kubernetes Resource Quota",
		Get: &plugin.GetConfig{
			Hydrate: opengovernance.GetKubernetesResourceQuota,
		},
		List: &plugin.ListConfig{
			Hydrate: opengovernance.ListKubernetesResourceQuota,
		},
		Columns: commonColumns([]*plugin.Column{

			//// ResourceQuotaSpec Columns
			{
				Name:        "spec_hard",
				Type:        proto.ColumnType_JSON,
				Description: "Spec hard is the set of desired hard limits for each named resource.",
				Transform:   transform.FromField("Description.ResourceQuota.Spec.Hard"),
			},
			{
				Name:        "spec_scopes",
				Type:        proto.ColumnType_JSON,
				Description: "A collection of filters that must match each object tracked by a quota.",
				Transform:   transform.FromField("Description.ResourceQuota.Spec.Scopes"),
			},
			{
				Name:        "spec_scope_selector",
				Type:        proto.ColumnType_JSON,
				Description: "A collection of filters like scopes that must match each object tracked by a quota but expressed using ScopeSelectorOperator in combination with possible values.",
				Transform:   transform.FromField("Description.ResourceQuota.Spec.ScopeSelector"),
			},

			//// ResourceQuotaStatus Columns
			{
				Name:        "status_hard",
				Type:        proto.ColumnType_JSON,
				Description: "Status hard is the set of enforced hard limits for each named resource.",
				Transform:   transform.FromField("Description.ResourceQuota.Status.Hard"),
			},
			{
				Name:        "status_used",
				Type:        proto.ColumnType_JSON,
				Description: "Indicates current observed total usage of the resource in the namespace.",
				Transform:   transform.FromField("Description.ResourceQuota.Status.Used"),
			},
			{
				Name:        "title",
				Type:        proto.ColumnType_STRING,
				Description: ColumnDescriptionTitle,
				Transform:   transform.FromField("Description.ResourceQuota.Name"),
			},
			{
				Name:        "tags",
				Type:        proto.ColumnType_JSON,
				Description: ColumnDescriptionTags,
				Transform:   transform.From(transformResourceQuotaTags),
			},
		}),
	}
}

func transformResourceQuotaTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	obj := d.HydrateItem.(opengovernance.KubernetesResourceQuota).Description.ResourceQuota
	return mergeTags(obj.Labels, obj.Annotations), nil
}
