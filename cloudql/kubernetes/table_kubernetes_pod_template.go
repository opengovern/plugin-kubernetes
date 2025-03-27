package kubernetes

import (
	"context"
	opengovernance "github.com/opengovern/og-describer-kubernetes/discovery/pkg/es"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableKubernetesPodTemplate(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "kubernetes_pod_template",
		Description: "Kubernetes Pod Template is a collection of templates for creating copies of a predefined pod.",
		List: &plugin.ListConfig{
			Hydrate: opengovernance.ListKubernetesPodTemplate,
		},
		Columns: commonColumns([]*plugin.Column{
			{
				Name:        "template",
				Type:        proto.ColumnType_JSON,
				Description: "Template describes the pods that will be created.",
				Transform:   transform.FromField("Description.PodTemplate.Template"),
			},
			{
				Name:        "title",
				Type:        proto.ColumnType_STRING,
				Description: ColumnDescriptionTitle,
				Transform:   transform.FromField("Description.PodTemplate.Name"),
			},
			{
				Name:        "tags",
				Type:        proto.ColumnType_JSON,
				Description: ColumnDescriptionTags,
				Transform:   transform.From(transformPodTemplateTags),
			},
		}),
	}
}

func transformPodTemplateTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	obj := d.HydrateItem.(opengovernance.KubernetesPodTemplate).Description.PodTemplate
	return mergeTags(obj.Labels, obj.Annotations), nil
}
