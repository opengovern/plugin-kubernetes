package kubernetes

import (
	"context"
	opengovernance "github.com/opengovern/og-describer-kubernetes/discovery/pkg/es"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableKubernetesCustomResourceDefinition(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "k8_custom_resource_definition",
		Description: "Kubernetes Custom Resource Definition.",
		List: &plugin.ListConfig{
			Hydrate: opengovernance.ListKubernetesCustomResourceDefinition,
		},
		Columns: commonColumns([]*plugin.Column{
			{
				Name:        "spec",
				Description: "Spec describes how the user wants the resources to appear.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Description.CustomResourceDefinition.Spec"),
			},
			{
				Name:        "status",
				Description: "Status indicates the actual state of the CustomResourceDefinition.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Description.CustomResourceDefinition.Status"),
			},
		}),
	}
}
