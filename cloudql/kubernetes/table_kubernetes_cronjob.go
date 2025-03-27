package kubernetes

import (
	"context"
	opengovernance "github.com/opengovern/og-describer-kubernetes/discovery/pkg/es"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableKubernetesCronJob(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "kubernetes_cronjob",
		Description: "Cron jobs are useful for creating periodic and recurring tasks, like running backups or sending emails.",
		List: &plugin.ListConfig{
			Hydrate: opengovernance.ListKubernetesCronJob,
		},
		Columns: commonColumns([]*plugin.Column{
			//// CronJobSpec columns
			{
				Name:        "failed_jobs_history_limit",
				Type:        proto.ColumnType_INT,
				Description: "The number of failed finished jobs to retain. Value must be non-negative integer.",
				Transform:   transform.FromField("Description.CronJob.Spec.FailedJobsHistoryLimit"),
			},
			{
				Name:        "schedule",
				Type:        proto.ColumnType_STRING,
				Description: "The schedule in Cron format.",
				Transform:   transform.FromField("Description.CronJob.Spec.Schedule"),
			},
			{
				Name:        "starting_deadline_seconds",
				Type:        proto.ColumnType_INT,
				Description: "Optional deadline in seconds for starting the job if it misses scheduledtime for any reason.",
				Transform:   transform.FromField("Description.CronJob.Spec.StartingDeadlineSeconds"),
			},
			{
				Name:        "successful_jobs_history_limit",
				Type:        proto.ColumnType_INT,
				Description: "The number of successful finished jobs to retain. Value must be non-negative integer.",
				Transform:   transform.FromField("Description.CronJob.Spec.SuccessfulJobsHistoryLimit"),
			},
			{
				Name:        "suspend",
				Type:        proto.ColumnType_BOOL,
				Description: "This flag tells the controller to suspend subsequent executions, it does not apply to already started executions.  Defaults to false.",
				Transform:   transform.FromField("Description.CronJob.Spec.Suspend"),
			},
			{
				Name:        "concurrency_policy",
				Type:        proto.ColumnType_JSON,
				Description: "Specifies how to treat concurrent executions of a Job.",
				Transform:   transform.FromField("Description.CronJob.Spec.ConcurrencyPolicy"),
			},
			{
				Name:        "job_template",
				Type:        proto.ColumnType_JSON,
				Description: "Specifies the job that will be created when executing a CronJob.",
				Transform:   transform.FromField("Description.CronJob.Spec.JobTemplate"),
			},

			//// CronJobStatus columns
			{
				Name:        "last_schedule_time",
				Type:        proto.ColumnType_TIMESTAMP,
				Description: "Information when was the last time the job was successfully scheduled.",
				Transform:   transform.FromField("Description.CronJob.Status.LastScheduleTime").Transform(v1TimeToRFC3339),
			},
			{
				Name:        "last_successful_time",
				Type:        proto.ColumnType_TIMESTAMP,
				Description: "Information when was the last time the job successfully completed.",
				Transform:   transform.FromField("Description.CronJob.Status.LastSuccessfulTime").Transform(v1TimeToRFC3339),
			},
			{
				Name:        "active",
				Type:        proto.ColumnType_JSON,
				Description: "A list of pointers to currently running jobs.",
				Transform:   transform.FromField("Description.CronJob.Status.Active"),
			},
			{
				Name:        "title",
				Type:        proto.ColumnType_STRING,
				Description: ColumnDescriptionTitle,
				Transform:   transform.FromField("Description.CronJob.Name"),
			},
			{
				Name:        "tags",
				Type:        proto.ColumnType_JSON,
				Description: ColumnDescriptionTags,
				Transform:   transform.From(transformCronJobTags),
			},
		}),
	}
}

func transformCronJobTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	obj := d.HydrateItem.(opengovernance.KubernetesCronJob).Description.CronJob
	return mergeTags(obj.Labels, obj.Annotations), nil
}
