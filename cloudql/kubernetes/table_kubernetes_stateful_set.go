package kubernetes

import (
	"context"
	opengovernance "github.com/opengovern/og-describer-kubernetes/discovery/pkg/es"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableKubernetesStatefulSet(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "k8_stateful_set",
		Description: "A statefulSet is the workload API object used to manage stateful applications.",
		List: &plugin.ListConfig{
			Hydrate: opengovernance.ListKubernetesStatefulSet,
		},
		// StatefulSet, is namespaced resource.
		Columns: commonColumns([]*plugin.Column{
			{
				Name:        "service_name",
				Type:        proto.ColumnType_STRING,
				Description: "The name of the service that governs this StatefulSet.",
				Transform:   transform.FromField("Description.StatefulSet.Spec.ServiceName"),
			},
			{
				Name:        "replicas",
				Type:        proto.ColumnType_INT,
				Description: "The desired number of replicas of the given Template.",
				Transform:   transform.FromField("Description.StatefulSet.Spec.Replicas"),
			},
			{
				Name:        "collision_count",
				Type:        proto.ColumnType_INT,
				Description: "The count of hash collisions for the StatefulSet.",
				Transform:   transform.FromField("Description.StatefulSet.Status.CollisionCount"),
			},
			{
				Name:        "available_replicas",
				Type:        proto.ColumnType_INT,
				Description: "Total number of available pods (ready for at least minReadySeconds) targeted by this statefulset.",
				Transform:   transform.FromField("Description.StatefulSet.Status.AvailableReplicas"),
			},
			{
				Name:        "current_replicas",
				Type:        proto.ColumnType_INT,
				Description: "The number of Pods created by the StatefulSet controller from the StatefulSet version indicated by currentRevision.",
				Transform:   transform.FromField("Description.StatefulSet.Status.CurrentReplicas"),
			},
			{
				Name:        "current_revision",
				Type:        proto.ColumnType_STRING,
				Description: "Indicates the version of the StatefulSet used to generate Pods in the sequence [0,currentReplicas).",
				Transform:   transform.FromField("Description.StatefulSet.Status.CurrentRevision"),
			},
			{
				Name:        "observed_generation",
				Type:        proto.ColumnType_INT,
				Description: "The most recent generation observed for this StatefulSet.",
				Transform:   transform.FromField("Description.StatefulSet.Status.ObservedGeneration"),
			},
			{
				Name:        "pod_management_policy",
				Type:        proto.ColumnType_STRING,
				Description: "Policy that controls how pods are created during initial scale up, when replacing pods on nodes, or when scaling down.",
				Transform:   transform.FromField("Description.StatefulSet.Spec.PodManagementPolicy").Transform(transform.ToString),
			},
			{
				Name:        "ready_replicas",
				Type:        proto.ColumnType_INT,
				Description: "The number of Pods created by the StatefulSet controller that have a Ready Condition.",
				Transform:   transform.FromField("Description.StatefulSet.Status.ReadyReplicas"),
			},
			{
				Name:        "revision_history_limit",
				Type:        proto.ColumnType_INT,
				Description: "The maximum number of revisions that will be maintained in the StatefulSet's revision history.",
				Transform:   transform.FromField("Description.StatefulSet.Spec.RevisionHistoryLimit"),
			},
			{
				Name:        "updated_replicas",
				Type:        proto.ColumnType_INT,
				Description: "The number of Pods created by the StatefulSet controller from the StatefulSet version indicated by updateRevision.",
				Transform:   transform.FromField("Description.StatefulSet.Status.UpdatedReplicas"),
			},
			{
				Name:        "update_revision",
				Type:        proto.ColumnType_STRING,
				Description: "Indicates the version of the StatefulSet used to generate Pods in the sequence [replicas-updatedReplicas,replicas).",
				Transform:   transform.FromField("Description.StatefulSet.Status.UpdateRevision"),
			},
			{
				Name:        "conditions",
				Type:        proto.ColumnType_JSON,
				Description: "Represents the latest available observations of a stateful set's current state.",
				Transform:   transform.FromField("Description.StatefulSet.Status.Conditions"),
			},
			{
				Name:        "template",
				Type:        proto.ColumnType_JSON,
				Description: "Template is the object that describes the pod that will be created if insufficient replicas are detected.",
				Transform:   transform.FromField("Description.StatefulSet.Spec.Template"),
			},
			{
				Name:        "update_strategy",
				Type:        proto.ColumnType_JSON,
				Description: "Indicates the StatefulSetUpdateStrategy that will be employed to update Pods in the StatefulSet when a revision is made to Template.",
				Transform:   transform.FromField("Description.StatefulSet.Spec.UpdateStrategy"),
			},
			{
				Name:        "title",
				Type:        proto.ColumnType_STRING,
				Description: ColumnDescriptionTitle,
				Transform:   transform.FromField("Description.StatefulSet.Name"),
			},
			{
				Name:        "tags",
				Type:        proto.ColumnType_JSON,
				Description: ColumnDescriptionTags,
				Transform:   transform.From(transformStatefulSetTags),
			},
		}),
	}
}

func transformStatefulSetTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	obj := d.HydrateItem.(opengovernance.KubernetesStatefulSet).Description.StatefulSet
	return mergeTags(obj.Labels, obj.Annotations), nil
}
