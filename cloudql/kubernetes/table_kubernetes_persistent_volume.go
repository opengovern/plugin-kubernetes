package kubernetes

import (
	"context"
	opengovernance "github.com/opengovern/og-describer-kubernetes/discovery/pkg/es"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableKubernetesPersistentVolume(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "kubernetes_persistent_volume",
		Description: "A PersistentVolume (PV) is a piece of storage in the cluster that has been provisioned by an administrator or dynamically provisioned using Storage Classes. PVs are volume plugins like Volumes, but have a lifecycle independent of any individual Pod that uses the PV.",
		Get: &plugin.GetConfig{
			Hydrate: opengovernance.GetKubernetesPersistentVolume,
		},
		List: &plugin.ListConfig{
			Hydrate: opengovernance.ListKubernetesPersistentVolume,
		},
		Columns: commonColumns([]*plugin.Column{
			//// PersistentVolumeSpec columns
			{
				Name:        "storage_class",
				Type:        proto.ColumnType_STRING,
				Description: "Name of StorageClass to which this persistent volume belongs. Empty value means that this volume does not belong to any StorageClass.",
				Transform:   transform.FromField("Description.PV.Spec.StorageClassName"),
			},
			{
				Name:        "volume_mode",
				Type:        proto.ColumnType_STRING,
				Description: "Defines if a volume is intended to be used with a formatted filesystem or to remain in raw block state.",
				Transform:   transform.FromField("Description.PV.Spec.VolumeMode"),
			},
			{
				Name:        "persistent_volume_reclaim_policy",
				Type:        proto.ColumnType_STRING,
				Description: "What happens to a persistent volume when released from its claim. Valid options are Retain (default for manually created PersistentVolumes), Delete (default for dynamically provisioned PersistentVolumes), and Recycle (deprecated). Recycle must be supported by the volume plugin underlying this PersistentVolume.",
				Transform:   transform.FromField("Description.PV.Spec.PersistentVolumeReclaimPolicy"),
			},
			{
				Name:        "access_modes",
				Type:        proto.ColumnType_JSON,
				Description: "List of ways the volume can be mounted.",
				Transform:   transform.FromField("Description.PV.Spec.AccessModes"),
			},
			{
				Name:        "capacity",
				Type:        proto.ColumnType_JSON,
				Description: "A description of the persistent volume's resources and capacity.",
				Transform:   transform.FromField("Description.PV.Spec.Capacity"),
			},
			{
				Name:        "claim_ref",
				Type:        proto.ColumnType_JSON,
				Description: "ClaimRef is part of a bi-directional binding between PersistentVolume and PersistentVolumeClaim. Expected to be non-nil when bound.",
				Transform:   transform.FromField("Description.PV.Spec.ClaimRef"),
			},
			{
				Name:        "mount_options",
				Type:        proto.ColumnType_JSON,
				Description: "A list of mount options, e.g. [\"ro\", \"soft\"].",
				Transform:   transform.FromField("Description.PV.Spec.MountOptions"),
			},
			{
				Name:        "node_affinity",
				Type:        proto.ColumnType_JSON,
				Description: "Defines constraints that limit what nodes this volume can be accessed from.",
				Transform:   transform.FromField("Description.PV.Spec.NodeAffinity"),
			},
			{
				Name:        "persistent_volume_source",
				Type:        proto.ColumnType_JSON,
				Description: "The actual volume backing the persistent volume.",
				Transform:   transform.FromField("Description.PV.Spec.PersistentVolumeSource"),
			},

			//// PersistentVolumeStatus columns
			{
				Name:        "phase",
				Type:        proto.ColumnType_STRING,
				Description: "Phase indicates if a volume is available, bound to a claim, or released by a claim.",
				Transform:   transform.FromField("Description.PV.Status.Phase"),
			},
			{
				Name:        "message",
				Type:        proto.ColumnType_STRING,
				Description: "A human-readable message indicating details about why the volume is in this state.",
				Transform:   transform.FromField("Description.PV.Status.Message"),
			},
			{
				Name:        "reason",
				Type:        proto.ColumnType_STRING,
				Description: "Reason is a brief CamelCase string that describes any failure and is meant for machine parsing and tidy display in the CLI.",
				Transform:   transform.FromField("Description.PV.Status.Reason"),
			},
			{
				Name:        "title",
				Type:        proto.ColumnType_STRING,
				Description: ColumnDescriptionTitle,
				Transform:   transform.FromField("Description.PV.Name"),
			},
			{
				Name:        "tags",
				Type:        proto.ColumnType_JSON,
				Description: ColumnDescriptionTags,
				Transform:   transform.From(transformPVTags),
			},
		}),
	}
}

//// TRANSFORM FUNCTIONS

func transformPVTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	obj := d.HydrateItem.(opengovernance.KubernetesPersistentVolume).Description.PV
	return mergeTags(obj.Labels, obj.Annotations), nil
}
