package kubernetes

import (
	"context"
	opengovernance "github.com/opengovern/og-describer-kubernetes/discovery/pkg/es"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableKubernetesJob(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "kubernetes_job",
		Description: "A Job creates one or more Pods and will continue to retry execution of the Pods until a specified number of them successfully terminate.",
		List: &plugin.ListConfig{
			Hydrate: opengovernance.ListKubernetesIngress,
		},
		Columns: commonColumns([]*plugin.Column{
			{
				Name:        "parallelism",
				Type:        proto.ColumnType_INT,
				Description: "The maximum desired number of pods the job should run at any given time. The actual number of pods running in steady state will be less than this number when ((.spec.completions - .status.successful) < .spec.parallelism), i.e. when the work left to do is less than max parallelism.",
				Transform:   transform.FromField("Description.Job.Spec.Parallelism"),
			},
			{
				Name:        "completions",
				Type:        proto.ColumnType_INT,
				Description: "The desired number of successfully finished pods the job should be run with.",
				Transform:   transform.FromField("Description.Job.Spec.Completions"),
			},
			{
				Name:        "active_deadline_seconds",
				Type:        proto.ColumnType_INT,
				Description: "The duration in seconds relative to the startTime that the job may be active before the system tries to terminate it.",
				Transform:   transform.FromField("Description.Job.Spec.ActiveDeadlineSeconds"),
			},
			{
				Name:        "backoff_limit",
				Type:        proto.ColumnType_INT,
				Description: "The number of retries before marking this job failed. Defaults to 6.",
				Transform:   transform.FromField("Description.Job.Spec.BackoffLimit"),
			},
			{
				Name:        "manual_selector",
				Type:        proto.ColumnType_BOOL,
				Description: "ManualSelector controls generation of pod labels and pod selectors. When false or unset, the system pick labels unique to this job and appends those labels to the pod template.  When true, the user is responsible for picking unique labels and specifying the selector.",
				Transform:   transform.FromField("Description.Job.Spec.ManualSelector"),
			},
			{
				Name:        "ttl_seconds_after_finished",
				Type:        proto.ColumnType_INT,
				Description: "limits the lifetime of a Job that has finished execution (either Complete or Failed). If this field is set, ttlSecondsAfterFinished after the Job finishes, it is eligible to be automatically deleted.",
				Transform:   transform.FromField("Description.Job.Spec.TTLSecondsAfterFinished"),
			},
			{
				Name:        "selector_query",
				Type:        proto.ColumnType_STRING,
				Description: "A query string representation of the selector.",
				Transform:   transform.FromField("Description.Job.Spec.Selector").Transform(labelSelectorToString),
			},
			{
				Name:        "selector",
				Type:        proto.ColumnType_JSON,
				Description: "A label query over pods that should match the pod count.",
				Transform:   transform.FromField("Description.Job.Spec.Selector"),
			},
			{
				Name:        "template",
				Type:        proto.ColumnType_JSON,
				Description: "Describes the pod that will be created when executing a job.",
				Transform:   transform.FromField("Description.Job.Spec.Template"),
			},

			//// JobStatus columns
			{
				Name:        "start_time",
				Type:        proto.ColumnType_TIMESTAMP,
				Description: "Time when the job was acknowledged by the job controller.",
				Transform:   transform.FromField("Description.Job.Status.StartTime").Transform(v1TimeToRFC3339),
			},
			{
				Name:        "completion_time",
				Type:        proto.ColumnType_TIMESTAMP,
				Description: "Time when the job was completed.",
				Transform:   transform.FromField("Description.Job.Status.CompletionTime").Transform(v1TimeToRFC3339),
			},
			{
				Name:        "active",
				Type:        proto.ColumnType_INT,
				Description: "The number of actively running pods.",
				Transform:   transform.FromField("Description.Job.Status.Active"),
			},
			{
				Name:        "succeeded",
				Type:        proto.ColumnType_INT,
				Description: "The number of pods which reached phase Succeeded.",
				Transform:   transform.FromField("Description.Job.Status.Succeeded"),
			},
			{
				Name:        "failed",
				Type:        proto.ColumnType_INT,
				Description: "The number of pods which reached phase Failed.",
				Transform:   transform.FromField("Description.Job.Status.Failed"),
			},
			{
				Name:        "conditions",
				Type:        proto.ColumnType_JSON,
				Description: "The latest available observations of an object's current state.",
				Transform:   transform.FromField("Description.Job.Status.Conditions"),
			},
			{
				Name:        "title",
				Type:        proto.ColumnType_STRING,
				Description: ColumnDescriptionTitle,
				Transform:   transform.FromField("Description.Job.Name"),
			},
			{
				Name:        "tags",
				Type:        proto.ColumnType_JSON,
				Description: ColumnDescriptionTags,
				Transform:   transform.From(transformJobTags),
			},
		}),
	}
}

func transformJobTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	obj := d.HydrateItem.(opengovernance.KubernetesJob).Description.Job
	return mergeTags(obj.Labels, obj.Annotations), nil
}
