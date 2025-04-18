package kubernetes

import (
	"context"
	opengovernance "github.com/opengovern/og-describer-kubernetes/discovery/pkg/es"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableKubernetesDaemonset(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "k8_daemonset",
		Description: "A DaemonSet ensures that all (or some) Nodes run a copy of a Pod.",
		List: &plugin.ListConfig{
			Hydrate: opengovernance.ListKubernetesDaemonSet,
		},
		Columns: commonColumns([]*plugin.Column{
			{
				Name:        "min_ready_seconds",
				Type:        proto.ColumnType_INT,
				Description: "The minimum number of seconds for which a newly created DaemonSet pod should be ready without any of its container crashing, for it to be considered available. Defaults to 0",
				Transform:   transform.FromField("Description.DaemonSet.Spec.MinReadySeconds"),
			},
			{
				Name:        "revision_history_limit",
				Type:        proto.ColumnType_INT,
				Description: "The number of old history to retain to allow rollback. This is a pointer to distinguish between explicit zero and not specified. Defaults to 10.",
				Transform:   transform.FromField("Description.DaemonSet.Spec.RevisionHistoryLimit"),
			},
			{
				Name:        "selector_query",
				Type:        proto.ColumnType_STRING,
				Description: "A query string representation of the selector.",
				Transform:   transform.FromField("Description.DaemonSet.Spec.Selector").Transform(labelSelectorToString),
			},
			{
				Name:        "selector",
				Type:        proto.ColumnType_JSON,
				Description: "A label query over pods that are managed by the daemon set.",
				Transform:   transform.FromField("Description.DaemonSet.Spec.Volumes"),
			},
			{
				Name:        "template",
				Type:        proto.ColumnType_JSON,
				Description: "An object that describes the pod that will be created.",
				Transform:   transform.FromField("Description.DaemonSet.Spec.Template"),
			},
			{
				Name:        "update_strategy",
				Type:        proto.ColumnType_JSON,
				Description: "An update strategy to replace existing DaemonSet pods with new pods.",
				Transform:   transform.FromField("Description.DaemonSet.Spec.UpdateStrategy"),
			},

			//// DaemonSetStatus Columns
			{
				Name:        "current_number_scheduled",
				Type:        proto.ColumnType_INT,
				Description: "The number of nodes that are running at least 1 daemon pod and are supposed to run the daemon pod.",
				Transform:   transform.FromField("Description.DaemonSet.Status.CurrentNumberScheduled"),
			},
			{
				Name:        "number_misscheduled",
				Type:        proto.ColumnType_INT,
				Description: "The number of nodes that are running the daemon pod, but are not supposed to run the daemon pod.",
				Transform:   transform.FromField("Description.DaemonSet.Status.NumberMisscheduled"),
			},
			{
				Name:        "desired_number_scheduled",
				Type:        proto.ColumnType_INT,
				Description: "The total number of nodes that should be running the daemon pod (including nodes correctly running the daemon pod).",
				Transform:   transform.FromField("Description.DaemonSet.Status.DesiredNumberScheduled"),
			},
			{
				Name:        "number_ready",
				Type:        proto.ColumnType_INT,
				Description: "The number of nodes that should be running the daemon pod and have one or more of the daemon pod running and ready.",
				Transform:   transform.FromField("Description.DaemonSet.Status.NumberReady"),
			},
			{
				Name:        "observed_generation",
				Type:        proto.ColumnType_INT,
				Description: "The most recent generation observed by the daemon set controller.",
				Transform:   transform.FromField("Description.DaemonSet.Status.ObservedGeneration"),
			},
			{
				Name:        "updated_number_scheduled",
				Type:        proto.ColumnType_INT,
				Description: "The total number of nodes that are running updated daemon pod.",
				Transform:   transform.FromField("Description.DaemonSet.Status.UpdatedNumberScheduled"),
			},
			{
				Name:        "number_available",
				Type:        proto.ColumnType_INT,
				Description: "The number of nodes that should be running the daemon pod and have one or more of the daemon pod running and available (ready for at least spec.minReadySeconds).",
				Transform:   transform.FromField("Description.DaemonSet.Status.NumberAvailable"),
			},
			{
				Name:        "number_unavailable",
				Type:        proto.ColumnType_INT,
				Description: "The number of nodes that should be running the daemon pod and have none of the daemon pod running and available (ready for at least spec.minReadySeconds).",
				Transform:   transform.FromField("Description.DaemonSet.Status.NumberUnavailable"),
			},
			{
				Name:        "collision_count",
				Type:        proto.ColumnType_INT,
				Description: "Count of hash collisions for the DaemonSet. The DaemonSet controller uses this field as a collision avoidance mechanism when it needs to create the name for the newest ControllerRevision.",
				Transform:   transform.FromField("Description.DaemonSet.Status.CollisionCount"),
			},
			{
				Name:        "conditions",
				Type:        proto.ColumnType_JSON,
				Description: "Represents the latest available observations of a DaemonSet's current state.",
				Transform:   transform.FromField("Description.DaemonSet.Status.Conditions"),
			},
			{
				Name:        "title",
				Type:        proto.ColumnType_STRING,
				Description: ColumnDescriptionTitle,
				Transform:   transform.FromField("Description.DaemonSet.Name"),
			},
			{
				Name:        "tags",
				Type:        proto.ColumnType_JSON,
				Description: ColumnDescriptionTags,
				Transform:   transform.From(transformDaemonSetTags),
			},
		}),
	}
}

func transformDaemonSetTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	obj := d.HydrateItem.(opengovernance.KubernetesDaemonSet).Description.DaemonSet
	return mergeTags(obj.Labels, obj.Annotations), nil
}
