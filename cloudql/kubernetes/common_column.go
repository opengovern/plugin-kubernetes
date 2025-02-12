package kubernetes

import (
	"context"
	"encoding/json"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

const (
	ColumnDescriptionTitle = "Title of the resource."
	ColumnDescriptionAkas  = "Array of globally unique identifier strings (also known as) for the resource."
	ColumnDescriptionTags  = "A map of tags for the resource. This includes both labels and annotations."
)

func objectMetadataColumns() []*plugin.Column {
	return []*plugin.Column{
		{Name: "name", Type: proto.ColumnType_STRING, Transform: transform.FromField("Description.MetaObject.Name"), Description: "Name of the object.  Name must be unique within a namespace."},
		{Name: "namespace", Type: proto.ColumnType_STRING, Transform: transform.FromField("Description.MetaObject.Namespace"), Description: "Namespace defines the space within which each name must be unique."},
		{Name: "uid", Type: proto.ColumnType_STRING, Description: "UID is the unique in time and space value for this object.", Transform: transform.FromField("Description.MetaObject.UID").Transform(transform.NullIfZeroValue)},
		{Name: "generate_name", Type: proto.ColumnType_STRING, Transform: transform.FromField("Description.MetaObject.GenerateName"), Description: "GenerateName is an optional prefix, used by the server, to generate a unique name ONLY IF the Name field has not been provided."},
		{Name: "resource_version", Type: proto.ColumnType_STRING, Transform: transform.FromField("Description.MetaObject.ResourceVersion"), Description: "An opaque value that represents the internal version of this object that can be used by clients to determine when objects have changed."},
		{Name: "generation", Type: proto.ColumnType_INT, Transform: transform.FromField("Description.MetaObject.Generation"), Description: "A sequence number representing a specific generation of the desired state."},
		{Name: "creation_timestamp", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("Description.MetaObject.CreationTimestamp").Transform(v1TimeToRFC3339), Description: "CreationTimestamp is a timestamp representing the server time when this object was created."},
		{Name: "deletion_timestamp", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("Description.MetaObject.DeletionTimestamp").Transform(v1TimeToRFC3339), Description: "DeletionTimestamp is RFC 3339 date and time at which this resource will be deleted."},
		{Name: "deletion_grace_period_seconds", Type: proto.ColumnType_INT, Transform: transform.FromField("Description.MetaObject.DeletionGracePeriodSeconds"), Description: "Number of seconds allowed for this object to gracefully terminate before it will be removed from the system.  Only set when deletionTimestamp is also set."},
		{Name: "labels", Type: proto.ColumnType_JSON, Transform: transform.FromField("Description.MetaObject.Labels"), Description: "Map of string keys and values that can be used to organize and categorize (scope and select) objects. May match selectors of replication controllers and services."},
		{Name: "annotations", Type: proto.ColumnType_JSON, Transform: transform.FromField("Description.MetaObject.Annotations"), Description: "Annotations is an unstructured key value map stored with a resource that may be set by external tools to store and retrieve arbitrary metadata."},
		{Name: "owner_references", Type: proto.ColumnType_JSON, Transform: transform.FromField("Description.MetaObject.OwnerReferences"), Description: "List of objects depended by this object. If ALL objects in the list have been deleted, this object will be garbage collected. If this object is managed by a controller, then an entry in this list will point to this controller, with the controller field set to true. There cannot be more than one managing controller."},
		{Name: "finalizers", Type: proto.ColumnType_JSON, Transform: transform.FromField("Description.MetaObject.Finalizers"), Description: "Must be empty before the object is deleted from the registry. Each entry is an identifier for the responsible component that will remove the entry from the list. If the deletionTimestamp of the object is non-nil, entries in this list can only be removed."},
	}
}

func commonColumns(c []*plugin.Column) []*plugin.Column {
	res := objectMetadataColumns()
	res = append(res, c...)
	res = append(res, []*plugin.Column{
		{
			Name:        "platform_integration_id",
			Type:        proto.ColumnType_STRING,
			Description: "The Platform Integration ID in which the resource is located.",
			Transform:   transform.FromField("IntegrationID"),
		},
		{
			Name:        "platform_resource_id",
			Type:        proto.ColumnType_STRING,
			Description: "The unique ID of the resource in opengovernance.",
			Transform:   transform.FromField("PlatformID"),
		},
		{
			Name:        "platform_metadata",
			Type:        proto.ColumnType_JSON,
			Description: "The metadata of the resource",
			Transform:   transform.FromField("Metadata").Transform(marshalJSON),
		},
		{
			Name:        "platform_resource_description",
			Type:        proto.ColumnType_JSON,
			Description: "The full model description of the resource",
			Transform:   transform.FromField("Description").Transform(marshalJSON),
		},
	}...)
	return res
}

func marshalJSON(_ context.Context, d *transform.TransformData) (interface{}, error) {
	b, err := json.Marshal(d.Value)
	if err != nil {
		return nil, err
	}
	return string(b), nil
}
