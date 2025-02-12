package kubernetes

import (
	"context"
	opengovernance "github.com/opengovern/og-describer-kubernetes/discovery/pkg/es"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableKubernetesDeployment(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "kubernetes_deployment",
		Description: "Kubernetes Deployment enables declarative updates for Pods and ReplicaSets.",
		Get: &plugin.GetConfig{
			Hydrate: opengovernance.GetKubernetesDeployment,
		},
		List: &plugin.ListConfig{
			Hydrate: opengovernance.ListKubernetesDeployment,
		},
		Columns: commonColumns([]*plugin.Column{
			//// Spec Columns
			{
				Name:        "replicas",
				Type:        proto.ColumnType_INT,
				Description: "Number of desired pods. Defaults to 1.",
				Transform:   transform.FromField("Description.Deployment.Spec.Replicas"),
			},
			{
				Name:        "selector_query",
				Type:        proto.ColumnType_STRING,
				Description: "A query string representation of the selector.",
				Transform:   transform.FromField("Description.Deployment.Spec.Selector").Transform(labelSelectorToString),
			},
			{
				Name:        "selector",
				Type:        proto.ColumnType_JSON,
				Description: "Label selector for pods. A label selector is a label query over a set of resources.",
				Transform:   transform.FromField("Description.Deployment.Spec.Selector"),
			},
			{
				Name:        "template",
				Type:        proto.ColumnType_JSON,
				Description: "Template describes the pods that will be created.",
				Transform:   transform.FromField("Description.Deployment.Spec.Template"),
			},
			{
				Name:        "strategy",
				Type:        proto.ColumnType_JSON,
				Description: "The deployment strategy to use to replace existing pods with new ones.",
				Transform:   transform.FromField("Description.Deployment.Spec.Strategy"),
			},
			{
				Name:        "min_ready_seconds",
				Type:        proto.ColumnType_INT,
				Description: "Minimum number of seconds for which a newly created pod should be ready without any of its container crashing, for it to be considered available. Defaults to 0.",
				Transform:   transform.FromField("Description.Deployment.Spec.MinReadySeconds"),
			},
			{
				Name:        "revision_history_limit",
				Type:        proto.ColumnType_INT,
				Description: "The number of old ReplicaSets to retain to allow rollback.",
				Transform:   transform.FromField("Description.Deployment.Spec.RevisionHistoryLimit"),
			},
			{
				Name:        "paused",
				Type:        proto.ColumnType_BOOL,
				Description: "Indicates that the deployment is paused.",
				Transform:   transform.FromField("Description.Deployment.Spec.Paused"),
			},
			{
				Name:        "progress_deadline_seconds",
				Type:        proto.ColumnType_INT,
				Description: "The maximum time in seconds for a deployment to make progress before it is considered to be failed.",
				Transform:   transform.FromField("Description.Deployment.Spec.ProgressDeadlineSeconds"),
			},

			//// Status Columns
			{
				Name:        "observed_generation",
				Type:        proto.ColumnType_INT,
				Description: "The generation observed by the deployment controller.",
				Transform:   transform.FromField("Description.Deployment.Status.ObservedGeneration"),
			},
			{
				Name:        "status_replicas",
				Type:        proto.ColumnType_INT,
				Description: "Total number of non-terminated pods targeted by this deployment (their labels match the selector).",
				Transform:   transform.FromField("Description.Deployment.Status.Replicas"),
			},
			{
				Name:        "updated_replicas",
				Type:        proto.ColumnType_INT,
				Description: "Total number of non-terminated pods targeted by this deployment that have the desired template spec.",
				Transform:   transform.FromField("Description.Deployment.Status.UpdatedReplicas"),
			},
			{
				Name:        "ready_replicas",
				Type:        proto.ColumnType_INT,
				Description: "Total number of ready pods targeted by this deployment.",
				Transform:   transform.FromField("Description.Deployment.Status.ReadyReplicas"),
			},
			{
				Name:        "available_replicas",
				Type:        proto.ColumnType_INT,
				Description: "Total number of available pods (ready for at least minReadySeconds) targeted by this deployment.",
				Transform:   transform.FromField("Description.Deployment.Status.AvailableReplicas"),
			},
			{
				Name:        "unavailable_replicas",
				Type:        proto.ColumnType_INT,
				Description: "Total number of unavailable pods targeted by this deployment.",
				Transform:   transform.FromField("Description.Deployment.Status.UnavailableReplicas"),
			},
			{
				Name:        "conditions",
				Type:        proto.ColumnType_JSON,
				Description: "Represents the latest available observations of a deployment's current state.",
				Transform:   transform.FromField("Description.Deployment.Status.Conditions"),
			},
			{
				Name:        "collision_count",
				Type:        proto.ColumnType_INT,
				Description: "Count of hash collisions for the Deployment. The Deployment controller uses this field as a collision avoidance mechanism when it needs to create the name for the newest ReplicaSet.",
				Transform:   transform.FromField("Description.Deployment.Status.CollisionCount"),
			},
			{
				Name:        "title",
				Type:        proto.ColumnType_STRING,
				Description: ColumnDescriptionTitle,
				Transform:   transform.FromField("Description.Deployment.Name"),
			},
			{
				Name:        "tags",
				Type:        proto.ColumnType_JSON,
				Description: ColumnDescriptionTags,
				Transform:   transform.From(transformDeploymentTags),
			},
		}),
	}
}

func transformDeploymentTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	obj := d.HydrateItem.(opengovernance.KubernetesDeployment).Description.Deployment
	return mergeTags(obj.Labels, obj.Annotations), nil
}
