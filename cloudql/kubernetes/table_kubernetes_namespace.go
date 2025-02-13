package kubernetes

import (
	"context"
	opengovernance "github.com/opengovern/og-describer-kubernetes/discovery/pkg/es"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableKubernetesNamespace(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "kubernetes_namespace",
		Description: "Kubernetes Namespace provides a scope for Names.",
		Get: &plugin.GetConfig{
			Hydrate: opengovernance.GetKubernetesNamespace,
		},
		List: &plugin.ListConfig{
			Hydrate: opengovernance.ListKubernetesNamespace,
		},
		Columns: commonColumns([]*plugin.Column{
			{
				Name:        "spec_finalizers",
				Type:        proto.ColumnType_JSON,
				Description: "Finalizers is an opaque list of values that must be empty to permanently remove object from storage.",
				Transform:   transform.FromField("Description.Namespace.Spec.Finalizers"),
			},

			//// NamespaceStatus Columns
			{
				Name:        "phase",
				Type:        proto.ColumnType_STRING,
				Description: "The current lifecycle phase of the namespace.",
				Transform:   transform.FromField("Description.Namespace.Status.Phase"),
			},
			{
				Name:        "conditions",
				Type:        proto.ColumnType_JSON,
				Description: "The latest available observations of namespace's current state.",
				Transform:   transform.FromField("Description.Namespace.Status.NamespaceCondition"),
			},
			{
				Name:        "title",
				Type:        proto.ColumnType_STRING,
				Description: ColumnDescriptionTitle,
				Transform:   transform.FromField("Description.Namespace.Name"),
			},
			{
				Name:        "tags",
				Type:        proto.ColumnType_JSON,
				Description: ColumnDescriptionTags,
				Transform:   transform.From(transformNamespaceTags),
			},
		}),
	}
}

func transformNamespaceTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	obj := d.HydrateItem.(opengovernance.KubernetesNamespace).Description.Namespace
	return mergeTags(obj.Labels, obj.Annotations), nil
}
