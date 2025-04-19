package helpers

import (
	"time"

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
)

type JobTemplateSpec struct {
	ObjectMeta
	Spec JobSpec
}

func ConvertJobTemplateSpec(spec batchv1.JobTemplateSpec) JobTemplateSpec {
	return JobTemplateSpec{
		ObjectMeta: ConvertObjectMeta(&spec.ObjectMeta),
		Spec:       ConvertJobSpec(spec.Spec),
	}
}

type CronJobSpec struct {
	Schedule                   string
	TimeZone                   *string
	StartingDeadlineSeconds    *int64
	ConcurrencyPolicy          string // batchv1.ConcurrencyPolicy
	Suspend                    *bool
	JobTemplate                JobTemplateSpec
	SuccessfulJobsHistoryLimit *int32
	FailedJobsHistoryLimit     *int32
}

func ConvertConcurrencyPolicy(cp batchv1.ConcurrencyPolicy) string {
	return string(cp)
}

func ConvertCronJobSpec(spec batchv1.CronJobSpec) CronJobSpec {
	return CronJobSpec{
		Schedule:                   spec.Schedule,
		TimeZone:                   spec.TimeZone,
		StartingDeadlineSeconds:    spec.StartingDeadlineSeconds,
		ConcurrencyPolicy:          ConvertConcurrencyPolicy(spec.ConcurrencyPolicy),
		Suspend:                    spec.Suspend,
		JobTemplate:                ConvertJobTemplateSpec(spec.JobTemplate),
		SuccessfulJobsHistoryLimit: spec.SuccessfulJobsHistoryLimit,
		FailedJobsHistoryLimit:     spec.FailedJobsHistoryLimit,
	}
}

func ConvertObjectReferences(objReferences []corev1.ObjectReference) []ObjectReference {
	references := make([]ObjectReference, len(objReferences))
	for i, objReference := range objReferences {
		references[i] = ObjectReference{
			Kind:            objReference.Kind,
			Namespace:       objReference.Namespace,
			Name:            objReference.Name,
			UID:             objReference.UID,
			APIVersion:      objReference.APIVersion,
			ResourceVersion: objReference.ResourceVersion,
			FieldPath:       objReference.FieldPath,
		}
	}
	return references
}

type CronJobStatus struct {
	Active             []ObjectReference
	LastScheduleTime   *time.Time
	LastSuccessfulTime *time.Time
}

func ConvertCronJobStatus(status batchv1.CronJobStatus) CronJobStatus {
	return CronJobStatus{
		Active:             ConvertObjectReferences(status.Active),
		LastScheduleTime:   ConvertTimePtr(status.LastScheduleTime),
		LastSuccessfulTime: ConvertTimePtr(status.LastSuccessfulTime),
	}
}

type CronJob struct {
	TypeMeta
	ObjectMeta
	Spec   CronJobSpec
	Status CronJobStatus
}

func ConvertCronJob(cj *batchv1.CronJob) CronJob {
	return CronJob{
		TypeMeta:   ConvertTypeMeta(cj.TypeMeta),
		ObjectMeta: ConvertObjectMeta(&cj.ObjectMeta),
		Spec:       ConvertCronJobSpec(cj.Spec),
		Status:     ConvertCronJobStatus(cj.Status),
	}
}
