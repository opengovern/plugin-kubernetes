package kubernetes

import (
	"context"
	opengovernance "github.com/opengovern/og-describer-kubernetes/discovery/pkg/es"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableKubernetesPersistentVolumeClaim(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "kubernetes_persistent_volume_claim",
		Description: "A PersistentVolumeClaim (PVC) is a request for storage by a user.",
		List: &plugin.ListConfig{
			Hydrate: opengovernance.ListKubernetesPersistentVolumeClaim,
		},
		Columns: commonColumns([]*plugin.Column{
			{
				Name:        "volume_name",
				Type:        proto.ColumnType_STRING,
				Description: "The binding reference to the PersistentVolume backing this claim.",
				Transform:   transform.FromField("Description.PVC.Spec.VolumeName"),
			},
			{
				Name:        "volume_mode",
				Type:        proto.ColumnType_STRING,
				Description: "Defines if a volume is intended to be used with a formatted filesystem or to remain in raw block state.",
				Transform:   transform.FromField("Description.PVC.Spec.VolumeMode"),
			},
			{
				Name:        "storage_class",
				Type:        proto.ColumnType_STRING,
				Description: "Name of StorageClass to which this persistent volume belongs. Empty value means that this volume does not belong to any StorageClass.",
				Transform:   transform.FromField("Description.PVC.Spec.StorageClassName"),
			},
			{
				Name:        "access_modes",
				Type:        proto.ColumnType_JSON,
				Description: "List of ways the volume can be mounted.",
				Transform:   transform.FromField("Description.PVC.Spec.AccessModes"),
			},
			{
				Name:        "data_source",
				Type:        proto.ColumnType_JSON,
				Description: "The source of the volume. This can be used to specify either: an existing VolumeSnapshot object (snapshot.storage.k8s.io/VolumeSnapshot), an existing PVC (PersistentVolumeClaim) or an existing custom resource that implements data population (Alpha).",
				Transform:   transform.FromField("Description.PVC.Spec.DataSource"),
			},
			{
				Name:        "resources",
				Type:        proto.ColumnType_JSON,
				Description: "Represents the minimum resources the volume should have.",
				Transform:   transform.FromField("Description.PVC.Spec.Resources"),
			},
			{
				Name:        "selector",
				Type:        proto.ColumnType_JSON,
				Description: "The actual volume backing the persistent volume.",
				Transform:   transform.FromField("Description.PVC.Spec.Selector"),
			},

			//// PersistentVolumeClaimStatus columns
			{
				Name:        "phase",
				Type:        proto.ColumnType_STRING,
				Description: "Phase indicates the current phase of PersistentVolumeClaim.",
				Transform:   transform.FromField("Description.PVC.Status.Phase"),
			},
			{
				Name:        "status_access_modes",
				Type:        proto.ColumnType_JSON,
				Description: "The actual access modes the volume backing the PVC has.",
				Transform:   transform.FromField("Description.PVC.Status.AccessModes"),
			},
			{
				Name:        "capacity",
				Type:        proto.ColumnType_JSON,
				Description: "The actual resources of the underlying volume.",
				Transform:   transform.FromField("Description.PVC.Status.Capacity"),
			},
			{
				Name:        "conditions",
				Type:        proto.ColumnType_JSON,
				Description: "The Condition of persistent volume claim.",
				Transform:   transform.FromField("Description.PVC.Status.Conditions"),
			},
			//// Steampipe Standard Columns
			{
				Name:        "title",
				Type:        proto.ColumnType_STRING,
				Description: ColumnDescriptionTitle,
				Transform:   transform.FromField("Description.PVC.Name"),
			},
			{
				Name:        "tags",
				Type:        proto.ColumnType_JSON,
				Description: ColumnDescriptionTags,
				Transform:   transform.From(transformPVCTags),
			},
		}),
	}
}

func transformPVCTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	obj := d.HydrateItem.(opengovernance.KubernetesPersistentVolumeClaim).Description.PVC
	return mergeTags(obj.Labels, obj.Annotations), nil
}
