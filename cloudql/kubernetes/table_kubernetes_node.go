package kubernetes

import (
	"context"
	opengovernance "github.com/opengovern/og-describer-kubernetes/discovery/pkg/es"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableKubernetesNode(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "k8_node",
		Description: "Kubernetes Node is a worker node in Kubernetes.",
		List: &plugin.ListConfig{
			Hydrate: opengovernance.ListKubernetesNode,
		},
		Columns: commonColumns([]*plugin.Column{
			{
				Name:        "pod_cidr",
				Type:        proto.ColumnType_CIDR,
				Description: "Pod IP range assigned to the node.",
				Transform:   transform.FromField("Description.Node.Spec.PodCIDR"),
			},
			{
				Name:        "pod_cidrs",
				Type:        proto.ColumnType_JSON,
				Description: "List of the IP ranges assigned to the node for usage by Pods.",
				Transform:   transform.FromField("Description.Node.Spec.PodCIDRs"),
			},
			{
				Name:        "provider_id",
				Type:        proto.ColumnType_STRING,
				Description: "ID of the node assigned by the cloud provider in the format: <ProviderName>://<ProviderSpecificNodeID>.",
				Transform:   transform.FromField("Description.Node.Spec.ProviderID"),
			},
			{
				Name:        "unschedulable",
				Type:        proto.ColumnType_BOOL,
				Description: "Unschedulable controls node schedulability of new pods. By default, node is schedulable.",
				Transform:   transform.FromField("Description.Node.Spec.Unschedulable"),
			},
			{
				Name:        "taints",
				Type:        proto.ColumnType_JSON,
				Description: "List of the taints attached to the node to has the \"effect\" on pod that does not tolerate the Taint",
				Transform:   transform.FromField("Description.Node.Spec.Taints"),
			},
			{
				Name:        "config_source",
				Type:        proto.ColumnType_JSON,
				Description: "The source to get node configuration from.",
				Transform:   transform.FromField("Description.Node.Spec.ConfigSource"),
			},

			//// NodeStatus Columns
			{
				Name:        "capacity_cpu",
				Type:        proto.ColumnType_STRING,
				Description: "Raw capacity CPU value as provided by the system.",
				Transform:   transform.FromP(transformNodeCpuAndMemory, "Capacity.CPU"),
			},
			{
				Name:        "capacity_memory",
				Type:        proto.ColumnType_STRING,
				Description: "Raw capacity memory value as provided by the system.",
				Transform:   transform.FromP(transformNodeCpuAndMemory, "Capacity.Memory"),
			},
			{
				Name:        "allocatable_cpu",
				Type:        proto.ColumnType_STRING,
				Description: "Raw allocatable CPU value as provided by the system.",
				Transform:   transform.FromP(transformNodeCpuAndMemory, "Allocatable.CPU"),
			},
			{
				Name:        "allocatable_memory",
				Type:        proto.ColumnType_STRING,
				Description: "Raw allocatable memory value as provided by the system.",
				Transform:   transform.FromP(transformNodeCpuAndMemory, "Allocatable.Memory"),
			},
			{
				Name:        "capacity_cpu_std",
				Type:        proto.ColumnType_INT,
				Description: "Standardized capacity CPU value in millicores (m).",
				Transform:   transform.FromP(transformNodeCpuAndMemoryUnit, "Capacity.CPU"),
			},
			{
				Name:        "capacity_memory_std",
				Type:        proto.ColumnType_INT,
				Description: "Standardized capacity memory value in bytes.",
				Transform:   transform.FromP(transformNodeCpuAndMemoryUnit, "Capacity.Memory"),
			},
			{
				Name:        "allocatable_cpu_std",
				Type:        proto.ColumnType_INT,
				Description: "Standardized allocatable CPU value in millicores (m).",
				Transform:   transform.FromP(transformNodeCpuAndMemoryUnit, "Allocatable.CPU"),
			},
			{
				Name:        "allocatable_memory_std",
				Type:        proto.ColumnType_INT,
				Description: "Standardized allocatable memory value in bytes.",
				Transform:   transform.FromP(transformNodeCpuAndMemoryUnit, "Allocatable.Memory"),
			},
			{
				Name:        "capacity",
				Type:        proto.ColumnType_JSON,
				Description: "Capacity represents the total resources of a node.",
				Transform:   transform.FromField("Description.Node.Status.Capacity"),
			},
			{
				Name:        "allocatable",
				Type:        proto.ColumnType_JSON,
				Description: "Allocatable represents the resources of a node that are available for scheduling. Defaults to capacity.",
				Transform:   transform.FromField("Description.Node.Status.Allocatable"),
			},
			{
				Name:        "phase",
				Type:        proto.ColumnType_STRING,
				Description: "The recently observed lifecycle phase of the node.",
				Transform:   transform.FromField("Description.Node.Status.Phase"),
			},
			{
				Name:        "conditions",
				Type:        proto.ColumnType_JSON,
				Description: "List of current observed node conditions.",
				Transform:   transform.FromField("Description.Node.Status.Conditions"),
			},
			{
				Name:        "addresses",
				Type:        proto.ColumnType_JSON,
				Description: "Endpoints of daemons running on the Node.",
				Transform:   transform.FromField("Description.Node.Status.Addresses"),
			},
			{
				Name:        "daemon_endpoints",
				Type:        proto.ColumnType_JSON,
				Description: "Set of ids/uuids to uniquely identify the node.",
				Transform:   transform.FromField("Description.Node.Status.DaemonEndpoints"),
			},
			{
				Name:        "node_info",
				Type:        proto.ColumnType_JSON,
				Description: "List of container images on this node.",
				Transform:   transform.FromField("Description.Node.Status.NodeInfo"),
			},
			{
				Name:        "images",
				Type:        proto.ColumnType_JSON,
				Description: "List of container images on this node.",
				Transform:   transform.FromField("Description.Node.Status.Images"),
			},
			{
				Name:        "volumes_in_use",
				Type:        proto.ColumnType_JSON,
				Description: "List of attachable volumes in use (mounted) by the node.",
				Transform:   transform.FromField("Description.Node.Status.VolumesInUse"),
			},
			{
				Name:        "volumes_attached",
				Type:        proto.ColumnType_JSON,
				Description: "List of volumes that are attached to the node.",
				Transform:   transform.FromField("Description.Node.Status.VolumesAttached"),
			},
			{
				Name:        "config",
				Type:        proto.ColumnType_JSON,
				Description: "Status of the config assigned to the node via the dynamic Kubelet config feature.",
				Transform:   transform.FromField("Description.Node.Status.Config"),
			},
			{
				Name:        "title",
				Type:        proto.ColumnType_STRING,
				Description: ColumnDescriptionTitle,
				Transform:   transform.FromField("Description.Node.Name"),
			},
			{
				Name:        "tags",
				Type:        proto.ColumnType_JSON,
				Description: ColumnDescriptionTags,
				Transform:   transform.From(transformNodeTags),
			},
		}),
	}
}

func transformNodeTags(_ context.Context, d *transform.TransformData) (any, error) {
	obj := d.HydrateItem.(opengovernance.KubernetesNode).Description.Node
	return mergeTags(obj.Labels, obj.Annotations), nil
}

func transformNodeCpuAndMemoryUnit(_ context.Context, d *transform.TransformData) (any, error) {
	param := d.Param.(string)

	node := d.HydrateItem.(opengovernance.KubernetesNode).Description.Node

	switch param {
	case "Capacity.CPU":
		if v, ok := node.Status.Capacity["cpu"]; ok {
			if vv, k := v.AsInt64(); k {
				return vv, nil
			}
			return nil, nil
		}
	case "Capacity.Memory":
		if v, ok := node.Status.Capacity["memory"]; ok {
			if vv, k := v.AsInt64(); k {
				return vv, nil
			}
			return nil, nil
		}
	case "Allocatable.CPU":
		if v, ok := node.Status.Allocatable["cpu"]; ok {
			if vv, k := v.AsInt64(); k {
				return vv, nil
			}
			return nil, nil
		}
	case "Allocatable.Memory":
		if v, ok := node.Status.Allocatable["memory"]; ok {
			if vv, k := v.AsInt64(); k {
				return vv, nil
			}
			return nil, nil
		}
	}

	return nil, nil
}

func transformNodeCpuAndMemory(_ context.Context, d *transform.TransformData) (any, error) {
	param := d.Param.(string)

	node := d.HydrateItem.(opengovernance.KubernetesNode).Description.Node
	switch param {
	case "Capacity.CPU":
		if v, ok := node.Status.Capacity["cpu"]; ok {
			return v.String(), nil
		}
	case "Capacity.Memory":
		if v, ok := node.Status.Capacity["memory"]; ok {
			return v.String(), nil
		}
	case "Allocatable.CPU":
		if v, ok := node.Status.Allocatable["cpu"]; ok {
			return v.String(), nil
		}
	case "Allocatable.Memory":
		if v, ok := node.Status.Allocatable["memory"]; ok {
			return v.String(), nil
		}
	}

	return nil, nil
}
