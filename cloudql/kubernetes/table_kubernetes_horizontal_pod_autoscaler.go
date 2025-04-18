package kubernetes

import (
	"context"
	opengovernance "github.com/opengovern/og-describer-kubernetes/discovery/pkg/es"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableKubernetesHorizontalPodAutoscaler(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "k8_horizontal_pod_autoscaler",
		Description: "Kubernetes HorizontalPodAutoscaler is the configuration for a horizontal pod autoscaler, which automatically manages the replica count of any resource implementing the scale subresource based on the metrics specified.",
		List: &plugin.ListConfig{
			Hydrate: opengovernance.ListKubernetesHorizontalPodAutoscaler,
		},
		Columns: commonColumns([]*plugin.Column{
			{
				Name:        "scale_target_ref",
				Type:        proto.ColumnType_JSON,
				Description: "ScaleTargetRef points to the target resource to scale, and is used to the pods for which metrics should be collected, as well as to actually change the replica count.",
				Transform:   transform.FromField("Description.HorizontalPodAutoscaler.Spec.ScaleTargetRef"),
			},
			{
				Name:        "min_replicas",
				Type:        proto.ColumnType_INT,
				Description: "MinReplicas is the lower limit for the number of replicas to which the autoscaler can scale down. It defaults to 1 pod. MinReplicas is allowed to be 0 if the alpha feature gate HPAScaleToZero is enabled and at least one Object or External metric is configured.",
				Transform:   transform.FromField("Description.HorizontalPodAutoscaler.Spec.MinReplicas"),
			},
			{
				Name:        "max_replicas",
				Type:        proto.ColumnType_INT,
				Description: "The Upper limit for the number of pods that can be set by the autoscaler. It cannot be smaller than MinReplicas.",
				Transform:   transform.FromField("Description.HorizontalPodAutoscaler.Spec.MaxReplicas"),
			},
			{
				Name:        "metrics",
				Type:        proto.ColumnType_JSON,
				Description: "Metrics contains the specifications for which to use to calculate the desired replica count (the maximum replica count across all metrics will be used). The desired replica count is calculated multiplying the ratio between the target value and the current value by the current number of pods.",
				Transform:   transform.FromField("Description.HorizontalPodAutoscaler.Spec.Metrics"),
			},
			{
				Name:        "scale_up_behavior",
				Type:        proto.ColumnType_JSON,
				Description: "Behavior configures the scaling behavior of the target in both Up and Down directions (scaleUp and scaleDown fields respectively). If not set, the default value is the higher of: * increase no more than 4 pods per 60 seconds * double the number of pods per 60 seconds.",
				Transform:   transform.FromField("Description.HorizontalPodAutoscaler.Spec.Behavior.ScaleUp"),
			},
			{
				Name:        "scale_down_behavior",
				Type:        proto.ColumnType_JSON,
				Description: "Behavior configures the scaling behavior of the target in both Up and Down directions (scaleUp and scaleDown fields respectively). If not set, the default value is to allow to scale down to minReplicas pods, with a 300 second stabilization window (i.e., the highest recommendation for the last 300sec is used).",
				Transform:   transform.FromField("Description.HorizontalPodAutoscaler.Spec.Behavior.ScaleDown"),
			},

			//// HpaStatus Columns
			{
				Name:        "observed_generation",
				Type:        proto.ColumnType_INT,
				Description: "The most recent generation observed by this autoscaler.",
				Transform:   transform.FromField("Description.HorizontalPodAutoscaler.Status.ObservedGeneration"),
			},
			{
				Name:        "last_scale_time",
				Type:        proto.ColumnType_TIMESTAMP,
				Description: "The last time the HorizontalPodAutoscaler scaled the number of pods used by the autoscaler to control how often the number of pods is changed.",
				Transform:   transform.FromField("Description.HorizontalPodAutoscaler.Status.LastScaleTime").Transform(v1TimeToRFC3339),
			},
			{
				Name:        "current_replicas",
				Type:        proto.ColumnType_INT,
				Description: "The current number of replicas of pods managed by this autoscaler.",
				Transform:   transform.FromField("Description.HorizontalPodAutoscaler.Status.CurrentReplicas"),
			},
			{
				Name:        "desired_replicas",
				Type:        proto.ColumnType_INT,
				Description: "The desired number of replicas of pods managed by this autoscaler.",
				Transform:   transform.FromField("Description.HorizontalPodAutoscaler.Status.DesiredReplicas"),
			},
			{
				Name:        "current_metrics",
				Type:        proto.ColumnType_JSON,
				Description: "CurrentMetrics is the last read state of the metrics used by this autoscaler.",
				Transform:   transform.FromField("Description.HorizontalPodAutoscaler.Status.CurrentMetrics"),
			},
			{
				Name:        "conditions",
				Type:        proto.ColumnType_JSON,
				Description: "Conditions is the set of conditions required for this autoscaler to scale its target and indicates whether or not those conditions are met.",
				Transform:   transform.FromField("Description.HorizontalPodAutoscaler.Status.Conditions"),
			},
			{
				Name:        "title",
				Type:        proto.ColumnType_STRING,
				Description: ColumnDescriptionTitle,
				Transform:   transform.FromField("Description.HorizontalPodAutoscaler.Name"),
			},
			{
				Name:        "tags",
				Type:        proto.ColumnType_JSON,
				Description: ColumnDescriptionTags,
				Transform:   transform.From(transformHpaTags),
			},
		}),
	}
}

func transformHpaTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	obj := d.HydrateItem.(opengovernance.KubernetesHorizontalPodAutoscaler).Description.HorizontalPodAutoscaler
	return mergeTags(obj.Labels, obj.Annotations), nil
}
