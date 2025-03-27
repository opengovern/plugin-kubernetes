package kubernetes

import (
	"context"
	opengovernance "github.com/opengovern/og-describer-kubernetes/discovery/pkg/es"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableKubernetesStorageClass(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "kubernetes_storage_class",
		Description: "Storage class provides a way for administrators to describe the classes of storage they offer.",
		List: &plugin.ListConfig{
			Hydrate: opengovernance.ListKubernetesStorageClass,
		},
		Columns: commonColumns([]*plugin.Column{
			{
				Name:        "provisioner",
				Type:        proto.ColumnType_STRING,
				Description: "Provisioner indicates the type of the provisioner.",
				Transform:   transform.FromField("Description.StorageClass.Provisioner"),
			},
			{
				Name:        "reclaim_policy",
				Type:        proto.ColumnType_STRING,
				Description: "Dynamically provisioned PersistentVolumes of this storage class are created with this reclaimPolicy. Defaults to Delete.",
				Transform:   transform.FromField("Description.StorageClass.ReclaimPolicy"),
			},
			{
				Name:        "allow_volume_expansion",
				Type:        proto.ColumnType_BOOL,
				Description: "AllowVolumeExpansion shows whether the storage class allows volume expand.",
				Transform:   transform.FromField("Description.StorageClass.AllowVolumeExpansion"),
			},
			{
				Name:        "volume_binding_mode",
				Type:        proto.ColumnType_STRING,
				Description: "VolumeBindingMode indicates how PersistentVolumeClaims should be provisioned and bound. When unset, VolumeBindingImmediate is used. This field is only honored by servers that enable the VolumeScheduling feature.",
				Transform:   transform.FromField("Description.StorageClass.VolumeBindingMode"),
			},
			{
				Name:        "allowed_topologies",
				Type:        proto.ColumnType_JSON,
				Description: "Restrict the node topologies where volumes can be dynamically provisioned. Each volume plugin defines its own supported topology specifications. An empty TopologySelectorTerm list means there is no topology restriction.",
				Transform:   transform.FromField("Description.StorageClass.AllowedTopologies"),
			},
			{
				Name:        "mount_options",
				Type:        proto.ColumnType_JSON,
				Description: "Dynamically provisioned PersistentVolumes of this storage class are created with these mountOptions, e.g. ['ro', 'soft']. Not validated - mount of the PVs will simply fail if one is invalid.",
				Transform:   transform.FromField("Description.StorageClass.MountOptions"),
			},
			{
				Name:        "parameters",
				Type:        proto.ColumnType_JSON,
				Description: "Parameters holds the parameters for the provisioner that should create volumes of this storage class.",
				Transform:   transform.FromField("Description.StorageClass.Parameters"),
			},
			{
				Name:        "title",
				Type:        proto.ColumnType_STRING,
				Description: ColumnDescriptionTitle,
				Transform:   transform.FromField("Description.StorageClass.Name"),
			},
		}),
	}
}
