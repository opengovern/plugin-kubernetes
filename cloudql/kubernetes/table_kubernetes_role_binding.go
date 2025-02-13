package kubernetes

import (
	"context"
	opengovernance "github.com/opengovern/og-describer-kubernetes/discovery/pkg/es"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableKubernetesRoleBinding(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "kubernetes_role_binding",
		Description: "A role binding grants the permissions defined in a role to a user or set of users. It holds a list of subjects (users, groups, or service accounts), and a reference to the role being granted.",
		Get: &plugin.GetConfig{
			Hydrate: opengovernance.GetKubernetesRoleBinding,
		},
		List: &plugin.ListConfig{
			Hydrate: opengovernance.ListKubernetesRoleBinding,
		},
		Columns: commonColumns([]*plugin.Column{
			{
				Name:        "subjects",
				Type:        proto.ColumnType_JSON,
				Description: "List of references to the objects the role applies to.",
				Transform:   transform.FromField("Description.RoleBinding.Subjects"),
			},

			//// RoleRef columns
			{
				Name:        "role_name",
				Type:        proto.ColumnType_STRING,
				Description: "Name of the role for which access is granted to subjects.",
				Transform:   transform.FromField("Description.RoleBinding.RoleRef.Name"),
			},
			{
				Name:        "role_api_group",
				Type:        proto.ColumnType_STRING,
				Description: "The group for the referenced role.",
				Transform:   transform.FromField("Description.RoleBinding.RoleRef.APIGroup"),
			},
			{
				Name:        "role_kind",
				Type:        proto.ColumnType_STRING,
				Description: "Type of the role referenced.",
				Transform:   transform.FromField("Description.RoleBinding.RoleRef.Kind"),
			},
			{
				Name:        "title",
				Type:        proto.ColumnType_STRING,
				Description: ColumnDescriptionTitle,
				Transform:   transform.FromField("Description.RoleBinding.Name"),
			},
			{
				Name:        "tags",
				Type:        proto.ColumnType_JSON,
				Description: ColumnDescriptionTags,
				Transform:   transform.From(transformRoleBindingTags),
			},
		}),
	}
}

func transformRoleBindingTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	obj := d.HydrateItem.(opengovernance.KubernetesRoleBinding).Description.RoleBinding
	return mergeTags(obj.Labels, obj.Annotations), nil
}
