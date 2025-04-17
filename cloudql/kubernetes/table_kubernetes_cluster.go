package kubernetes

import (
	"context"
	opengovernance "github.com/opengovern/og-describer-kubernetes/discovery/pkg/es"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableKubernetesCluster(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "kubernetes_cluster",
		Description: "ClusterRole contains rules that represent a set of permissions.",
		List: &plugin.ListConfig{
			Hydrate: opengovernance.ListKubernetesCluster,
		},
		Columns: commonColumns([]*plugin.Column{
			{
				Name:        "context_name",
				Type:        proto.ColumnType_STRING,
				Description: "Context name of the cluster.",
				Transform:   transform.FromField("Description.ContextName"),
			},
			{
				Name:        "endpoint",
				Type:        proto.ColumnType_STRING,
				Description: "endpoint of the cluster.",
				Transform:   transform.FromField("Description.Endpoint"),
			},
			{
				Name:        "auth_method",
				Type:        proto.ColumnType_STRING,
				Description: "cluster auth method.",
				Transform:   transform.FromField("Description.AuthMethod"),
			},
			{
				Name:        "tls_server_verification",
				Type:        proto.ColumnType_BOOL,
				Description: "endpoint of the cluster.",
				Transform:   transform.FromField("Description.TLSServerVerification"),
			},
			{
				Name:        "server_version",
				Type:        proto.ColumnType_STRING,
				Description: "endpoint of the cluster.",
				Transform:   transform.FromField("Description.ServerVersion"),
			},
		}),
	}
}
