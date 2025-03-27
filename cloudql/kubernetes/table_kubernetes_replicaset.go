package kubernetes

import (
	"context"
	opengovernance "github.com/opengovern/og-describer-kubernetes/discovery/pkg/es"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableKubernetesReplicaSet(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "kubernetes_replicaset",
		Description: "Kubernetes replica set ensures that a specified number of pod replicas are running at any given time.",
		List: &plugin.ListConfig{
			Hydrate: opengovernance.ListKubernetesReplicaSet,
		},
		Columns: commonColumns([]*plugin.Column{
			{
				Name:        "replicas",
				Type:        proto.ColumnType_INT,
				Description: "Replicas is the number of desired replicas. Defaults to 1.",
				Transform:   transform.FromField("Description.ReplicaSet.Spec.Replicas"),
			},
			{
				Name:        "min_ready_seconds",
				Type:        proto.ColumnType_INT,
				Description: "Minimum number of seconds for which a newly created pod should be ready without any of its container crashing, for it to be considered available. Defaults to 0",
				Transform:   transform.FromField("Description.ReplicaSet.Spec.MinReadySeconds"),
			},
			{
				Name:        "selector_query",
				Type:        proto.ColumnType_STRING,
				Description: "A query string representation of the selector.",
				Transform:   transform.FromField("Description.ReplicaSet.Spec.Selector").Transform(labelSelectorToString),
			},
			{
				Name:        "selector",
				Type:        proto.ColumnType_JSON,
				Description: "Selector is a label query over pods that should match the replica count. Label keys and values that must match in order to be controlled by this replica set.",
				Transform:   transform.FromField("Description.ReplicaSet.Spec.Selector"),
			},
			{
				Name:        "template",
				Type:        proto.ColumnType_JSON,
				Description: "Template is the object that describes the pod that will be created if insufficient replicas are detected.",
				Transform:   transform.FromField("Description.ReplicaSet.Spec.Template"),
			},

			//// Status Columns
			{
				Name:        "status_replicas",
				Type:        proto.ColumnType_INT,
				Description: "The most recently oberved number of replicas.",
				Transform:   transform.FromField("Description.ReplicaSet.Status.Replicas"),
			},
			{
				Name:        "fully_labeled_replicas",
				Type:        proto.ColumnType_INT,
				Description: "The number of pods that have labels matching the labels of the pod template of the replicaset.",
				Transform:   transform.FromField("Description.ReplicaSet.Status.FullyLabeledReplicas"),
			},
			{
				Name:        "ready_replicas",
				Type:        proto.ColumnType_INT,
				Description: "The number of ready replicas for this replica set.",
				Transform:   transform.FromField("Description.ReplicaSet.Status.ReadyReplicas"),
			},
			{
				Name:        "available_replicas",
				Type:        proto.ColumnType_INT,
				Description: "The number of available replicas (ready for at least minReadySeconds) for this replica set.",
				Transform:   transform.FromField("Description.ReplicaSet.Status.AvailableReplicas"),
			},
			{
				Name:        "observed_generation",
				Type:        proto.ColumnType_INT,
				Description: "ObservedGeneration reflects the generation of the most recently observed ReplicaSet.",
				Transform:   transform.FromField("Description.ReplicaSet.Status.ObservedGeneration"),
			},
			{
				Name:        "conditions",
				Type:        proto.ColumnType_JSON,
				Description: "Represents the latest available observations of a replica set's current state.",
				Transform:   transform.FromField("Description.ReplicaSet.Status.Conditions"),
			},
			{
				Name:        "title",
				Type:        proto.ColumnType_STRING,
				Description: ColumnDescriptionTitle,
				Transform:   transform.FromField("Description.ReplicaSet.Name"),
			},
			{
				Name:        "tags",
				Type:        proto.ColumnType_JSON,
				Description: ColumnDescriptionTags,
				Transform:   transform.From(transformReplicaSetTags),
			},
		}),
	}
}

func transformReplicaSetTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	obj := d.HydrateItem.(opengovernance.KubernetesReplicaSet).Description.ReplicaSet
	return mergeTags(obj.Labels, obj.Annotations), nil
}
