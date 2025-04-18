package kubernetes

import (
	"context"
	opengovernance "github.com/opengovern/og-describer-kubernetes/discovery/pkg/es"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableKubernetesEvent(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "k8_event",
		Description: "Kubernetes Event is a report of an event somewhere in the cluster.",
		List: &plugin.ListConfig{
			Hydrate: opengovernance.ListKubernetesEvent,
		},
		Columns: commonColumns([]*plugin.Column{
			{
				Name:        "last_timestamp",
				Type:        proto.ColumnType_TIMESTAMP,
				Description: "Time when this event was last observed.",
				Transform:   transform.FromField("Description.Event.LastTimestamp").Transform(v1TimeToRFC3339),
			},
			{
				Name:        "type",
				Type:        proto.ColumnType_STRING,
				Description: "Type of this event (Normal, Warning), new types could be added in the future.",
				Transform:   transform.FromField("Description.Event.Type"),
			},
			{
				Name:        "reason",
				Type:        proto.ColumnType_STRING,
				Description: "The reason the transition into the object's current status.",
				Transform:   transform.FromField("Description.Event.Reason"),
			},
			{
				Name:        "message",
				Type:        proto.ColumnType_STRING,
				Description: "A description of the status of this operation.",
				Transform:   transform.FromField("Description.Event.Message"),
			},
			{
				Name:        "action",
				Type:        proto.ColumnType_STRING,
				Description: "What action was taken/failed with the regarding object.",
				Transform:   transform.FromField("Description.Event.Action"),
			},
			{
				Name:        "count",
				Type:        proto.ColumnType_INT,
				Description: "The number of times this event has occurred.",
				Transform:   transform.FromField("Description.Event.Count"),
			},
			{
				Name:        "event_time",
				Type:        proto.ColumnType_TIMESTAMP,
				Description: "Time when this event was first observed.",
				Transform:   transform.FromField("Description.Event.EventTime").Transform(v1MicroTimeToRFC3339),
			},
			{
				Name:        "first_timestamp",
				Type:        proto.ColumnType_TIMESTAMP,
				Description: "The time at which the event was first recorded.",
				Transform:   transform.FromField("Description.Event.FirstTimestamp").Transform(v1TimeToRFC3339),
			},
			{
				Name:        "reporting_component",
				Type:        proto.ColumnType_STRING,
				Description: "Name of the controller that emitted this event.",
				Transform:   transform.FromField("Description.Event.ReportingComponent"),
			},
			{
				Name:        "reporting_instance",
				Type:        proto.ColumnType_STRING,
				Description: "ID of the controller instance.",
				Transform:   transform.FromField("Description.Event.ReportingInstance"),
			},
			{
				Name:        "involved_object",
				Type:        proto.ColumnType_JSON,
				Description: "The object that this event is about.",
				Transform:   transform.FromField("Description.Event.InvolvedObject"),
			},
			{
				Name:        "related",
				Type:        proto.ColumnType_JSON,
				Description: "Optional secondary object for more complex actions.",
				Transform:   transform.FromField("Description.Event.Related"),
			},
			{
				Name:        "series",
				Type:        proto.ColumnType_JSON,
				Description: "Data about the event series this event represents.",
				Transform:   transform.FromField("Description.Event.Series"),
			},
			{
				Name:        "source",
				Type:        proto.ColumnType_JSON,
				Description: "The component reporting this event.",
				Transform:   transform.FromField("Description.Event.Source"),
			},
		}),
	}
}
