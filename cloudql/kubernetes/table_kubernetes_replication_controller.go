package kubernetes

import (
	"context"
	opengovernance "github.com/opengovern/og-describer-kubernetes/discovery/pkg/es"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableKubernetesReplicaController(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "kubernetes_replication_controller",
		Description: "A Replication Controller makes sure that a pod or homogeneous set of pods are always up and available. If there are too many pods, it will kill some. If there are too few, the Replication Controller will start more.",
		Get: &plugin.GetConfig{
			Hydrate: opengovernance.GetKubernetesReplicationController,
		},
		List: &plugin.ListConfig{
			Hydrate: opengovernance.ListKubernetesReplicationController,
		},
		Columns: commonColumns([]*plugin.Column{
			{
				Name:        "replicas",
				Type:        proto.ColumnType_INT,
				Description: "Replicas is the number of desired replicas. Defaults to 1.",
				Transform:   transform.FromField("Description.ReplicationController.Spec.Replicas"),
			},
			{
				Name:        "min_ready_seconds",
				Type:        proto.ColumnType_INT,
				Description: "Minimum number of seconds for which a newly created pod should be ready without any of its container crashing, for it to be considered available. Defaults to 0",
				Transform:   transform.FromField("Description.ReplicationController.Spec.MinReadySeconds"),
			},
			{
				Name:        "selector_query",
				Type:        proto.ColumnType_STRING,
				Description: "A query string representation of the selector.",
				Transform:   transform.FromField("Description.ReplicationController.Spec.Selector").Transform(selectorMapToString),
			},
			{
				Name:        "selector",
				Type:        proto.ColumnType_JSON,
				Description: "Selector is a label query over pods that should match the replica count. Label keys and values that must match in order to be controlled by this replica set.",
				Transform:   transform.FromField("Description.ReplicationController.Spec.Selector"),
			},
			{
				Name:        "template",
				Type:        proto.ColumnType_JSON,
				Description: "Template is the object that describes the pod that will be created if insufficient replicas are detected.",
				Transform:   transform.FromField("Description.ReplicationController.Spec.Template"),
			},

			//// Status Columns
			{
				Name:        "status_replicas",
				Type:        proto.ColumnType_INT,
				Description: "The most recently oberved number of replicas.",
				Transform:   transform.FromField("Description.ReplicationController.Status.Replicas"),
			},
			{
				Name:        "fully_labeled_replicas",
				Type:        proto.ColumnType_INT,
				Description: "The number of pods that have labels matching the labels of the pod template of the replicaset.",
				Transform:   transform.FromField("Description.ReplicationController.Status.FullyLabeledReplicas"),
			},
			{
				Name:        "ready_replicas",
				Type:        proto.ColumnType_INT,
				Description: "The number of ready replicas for this replica set.",
				Transform:   transform.FromField("Description.ReplicationController.Status.ReadyReplicas"),
			},
			{
				Name:        "available_replicas",
				Type:        proto.ColumnType_INT,
				Description: "The number of available replicas (ready for at least minReadySeconds) for this replica set.",
				Transform:   transform.FromField("Description.ReplicationController.Status.AvailableReplicas"),
			},
			{
				Name:        "observed_generation",
				Type:        proto.ColumnType_INT,
				Description: "Reflects the generation of the most recently observed replication controller.",
				Transform:   transform.FromField("Description.ReplicationController.Status.ObservedGeneration"),
			},
			{
				Name:        "conditions",
				Type:        proto.ColumnType_JSON,
				Description: "Represents the latest available observations of a replication controller's current state.",
				Transform:   transform.FromField("Description.ReplicationController.Status.Conditions"),
			},
			{
				Name:        "title",
				Type:        proto.ColumnType_STRING,
				Description: ColumnDescriptionTitle,
				Transform:   transform.FromField("Description.ReplicationController.Name"),
			},
			{
				Name:        "tags",
				Type:        proto.ColumnType_JSON,
				Description: ColumnDescriptionTags,
				Transform:   transform.From(transformReplicaControllerTags),
			},
		}),
	}
}

func transformReplicaControllerTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	obj := d.HydrateItem.(opengovernance.KubernetesReplicationController).Description.ReplicationController
	return mergeTags(obj.Labels, obj.Annotations), nil
}
