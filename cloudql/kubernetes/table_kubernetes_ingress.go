package kubernetes

import (
	"context"
	opengovernance "github.com/opengovern/og-describer-kubernetes/discovery/pkg/es"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableKubernetesIngress(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "kubernetes_ingress",
		Description: "Ingress exposes HTTP and HTTPS routes from outside the cluster to services within the cluster. Traffic routing is controlled by rules defined on the Ingress resource.",
		Get: &plugin.GetConfig{
			Hydrate: opengovernance.GetKubernetesIngress,
		},
		List: &plugin.ListConfig{
			Hydrate: opengovernance.ListKubernetesIngress,
		},
		Columns: commonColumns([]*plugin.Column{
			{
				Name:        "ingress_class_name",
				Type:        proto.ColumnType_STRING,
				Description: "Name of the IngressClass cluster resource. The associated IngressClass defines which controller will implement the resource.",
				Transform:   transform.FromField("Description.Ingress.Spec.IngressClassName"),
			},
			{
				Name:        "default_backend",
				Type:        proto.ColumnType_JSON,
				Description: "A default backend capable of servicing requests that don't match any rule. At least one of 'backend' or 'rules' must be specified.",
				Transform:   transform.FromField("Description.Ingress.Spec.DefaultBackend"),
			},
			{
				Name:        "tls",
				Type:        proto.ColumnType_JSON,
				Description: "TLS configuration.",
				Transform:   transform.FromField("Description.Ingress.Spec.TLS"),
			},
			{
				Name:        "rules",
				Type:        proto.ColumnType_JSON,
				Description: "A list of host rules used to configure the Ingress.",
				Transform:   transform.FromField("Description.Ingress.Spec.Rules"),
			},
			{
				Name:        "load_balancer",
				Type:        proto.ColumnType_JSON,
				Description: "a list containing ingress points for the load-balancer. Traffic intended for the service should be sent to these ingress points.",
				Transform:   transform.FromField("Description.Ingress.Status.LoadBalancer.Ingress"),
			},
			{
				Name:        "title",
				Type:        proto.ColumnType_STRING,
				Description: ColumnDescriptionTitle,
				Transform:   transform.FromField("Description.Ingress.Name"),
			},
			{
				Name:        "tags",
				Type:        proto.ColumnType_JSON,
				Description: ColumnDescriptionTags,
				Transform:   transform.From(transformIngressTags),
			},
		}),
	}
}

func transformIngressTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	obj := d.HydrateItem.(opengovernance.KubernetesIngress).Description.Ingress
	return mergeTags(obj.Labels, obj.Annotations), nil
}
