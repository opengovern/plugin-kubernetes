package kubernetes

import (
	"context"
	opengovernance "github.com/opengovern/og-describer-kubernetes/discovery/pkg/es"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableKubernetesServiceAccount(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "kubernetes_service_account",
		Description: "A service account provides an identity for processes that run in a Pod.",
		Get: &plugin.GetConfig{
			Hydrate: opengovernance.GetKubernetesServiceAccount,
		},
		List: &plugin.ListConfig{
			Hydrate: opengovernance.ListKubernetesServiceAccount,
		},
		// Service Account, is namespaced resource.
		Columns: commonColumns([]*plugin.Column{
			{
				Name:        "automount_service_account_token",
				Type:        proto.ColumnType_BOOL,
				Description: "Indicates whether pods running as this service account should have an API token automatically mounted. Can be overridden at the pod level.",
				Transform:   transform.FromField("Description.ServiceAccount.AutomountServiceAccountToken"),
			},
			{
				Name:        "image_pull_secrets",
				Type:        proto.ColumnType_JSON,
				Description: "List of references to secrets in the same namespace to use for pulling any images in pods that reference this ServiceAccount. ImagePullSecrets are distinct from Secrets because Secrets can be mounted in the pod, but ImagePullSecrets are only accessed by the kubelet.",
				Transform:   transform.FromField("Description.ServiceAccount.ImagePullSecrets"),
			},
			{
				Name:        "secrets",
				Type:        proto.ColumnType_JSON,
				Description: "Secrets is the list of secrets allowed to be used by pods running using this ServiceAccount.",
				Transform:   transform.FromField("Description.ServiceAccount.Secrets"),
			},
			{
				Name:        "title",
				Type:        proto.ColumnType_STRING,
				Description: ColumnDescriptionTitle,
				Transform:   transform.FromField("Description.ServiceAccount.Name"),
			},
			{
				Name:        "tags",
				Type:        proto.ColumnType_JSON,
				Description: ColumnDescriptionTags,
				Transform:   transform.From(transformServiceAccountTags),
			},
		}),
	}
}

func transformServiceAccountTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	obj := d.HydrateItem.(opengovernance.KubernetesServiceAccount).Description.ServiceAccount
	return mergeTags(obj.Labels, obj.Annotations), nil
}
