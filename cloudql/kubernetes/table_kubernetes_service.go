package kubernetes

import (
	"context"
	opengovernance "github.com/opengovern/og-describer-kubernetes/discovery/pkg/es"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableKubernetesService(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "kubernetes_service",
		Description: "A service provides an abstract way to expose an application running on a set of Pods as a network service.",
		List: &plugin.ListConfig{
			Hydrate: opengovernance.ListKubernetesService,
		},
		// Service is namespaced resource.
		Columns: commonColumns([]*plugin.Column{
			{
				Name:        "type",
				Type:        proto.ColumnType_STRING,
				Description: "Type determines how the Service is exposed.",
				Transform:   transform.FromField("Description.Service.Spec.Type").Transform(transform.ToString),
			},
			{
				Name:        "allocate_load_balancer_node_ports",
				Type:        proto.ColumnType_BOOL,
				Description: "Indicates whether NodePorts will be automatically allocated for services with type LoadBalancer, or not.",
				Transform:   transform.FromField("Description.Service.Spec.AllocateLoadBalancerNodePorts"),
			},
			{
				Name:        "cluster_ip",
				Type:        proto.ColumnType_STRING,
				Description: "IP address of the service and is usually assigned randomly.",
				Transform:   transform.FromField("Description.Service.Spec.ClusterIP"),
			},
			{
				Name:        "external_name",
				Type:        proto.ColumnType_STRING,
				Description: "The external reference that discovery mechanisms will return as an alias for this service (e.g. a DNS CNAME record).",
				Transform:   transform.FromField("Description.Service.Spec.ExternalName"),
			},
			{
				Name:        "external_traffic_policy",
				Type:        proto.ColumnType_STRING,
				Description: "Denotes whether the service desires to route external traffic to node-local or cluster-wide endpoints.",
				Transform:   transform.FromField("Description.Service.Spec.ExternalTrafficPolicy").Transform(transform.ToString),
			},
			{
				Name:        "health_check_node_port",
				Type:        proto.ColumnType_INT,
				Description: "Specifies the healthcheck nodePort for the service.",
				Transform:   transform.FromField("Description.Service.Spec.HealthCheckNodePort"),
			},
			{
				Name:        "ip_family_policy",
				Type:        proto.ColumnType_STRING,
				Description: "Specifies the dual-stack-ness requested or required by this service, and is gated by the 'IPv6DualStack' feature gate.",
				Transform:   transform.FromField("Description.Service.Spec.IPFamilyPolicy").Transform(transform.ToString),
			},
			{
				Name:        "load_balancer_ip",
				Type:        proto.ColumnType_IPADDR,
				Description: "The IP specified when the load balancer was created.",
				Transform:   transform.FromField("Description.Service.Spec.LoadBalancerIP"),
			},
			{
				Name:        "publish_not_ready_addresses",
				Type:        proto.ColumnType_BOOL,
				Description: "Indicates that any agent which deals with endpoints for this service should disregard any indications of ready/not-ready.",
				Transform:   transform.FromField("Description.Service.Spec.PublishNotReadyAddresses"),
			},
			{
				Name:        "session_affinity",
				Type:        proto.ColumnType_STRING,
				Description: "Supports 'ClientIP' and 'None'. Used to maintain session affinity.",
				Transform:   transform.FromField("Description.Service.Spec.SessionAffinity").Transform(transform.ToString),
			},
			{
				Name:        "session_affinity_client_ip_timeout",
				Type:        proto.ColumnType_INT,
				Description: "Specifies the ClientIP type session sticky time in seconds.",
				Transform:   transform.FromField("Description.Service.Spec.SessionAffinityConfig.ClientIP.TimeoutSeconds"),
			},
			{
				Name:        "cluster_ips",
				Type:        proto.ColumnType_JSON,
				Description: "A list of IP addresses assigned to this service, and are usually assigned randomly.",
				Transform:   transform.FromField("Description.Service.Spec.ClusterIPs"),
			},
			{
				Name:        "external_ips",
				Type:        proto.ColumnType_JSON,
				Description: "A list of IP addresses for which nodes in the cluster will also accept traffic for this service.",
				Transform:   transform.FromField("Description.Service.Spec.ExternalIPs"),
			},
			{
				Name:        "ip_families",
				Type:        proto.ColumnType_JSON,
				Description: "A list of IP families (e.g. IPv4, IPv6) assigned to this service, and is gated by the 'IPv6DualStack' feature gate.",
				Transform:   transform.FromField("Description.Service.Spec.IPFamilies"),
			},
			{
				Name:        "load_balancer_ingress",
				Type:        proto.ColumnType_JSON,
				Description: "A list containing ingress points for the load-balancer.",
				Transform:   transform.FromField("Description.Service.Status.LoadBalancer.Ingress"),
			},
			{
				Name:        "load_balancer_source_ranges",
				Type:        proto.ColumnType_JSON,
				Description: "A list of source ranges that will restrict traffic through the cloud-provider load-balancer will be restricted to the specified client IPs.",
				Transform:   transform.FromField("Description.Service.Spec.LoadBalancerSourceRanges"),
			},
			{
				Name:        "ports",
				Type:        proto.ColumnType_JSON,
				Description: "A list of ports that are exposed by this service.",
				Transform:   transform.FromField("Description.Service.Spec.Ports"),
			},
			{
				Name:        "selector_query",
				Type:        proto.ColumnType_STRING,
				Description: "A query string representation of the selector.",
				Transform:   transform.FromField("Description.Service.Spec.Selector").Transform(selectorMapToString),
			},
			{
				Name:        "selector",
				Type:        proto.ColumnType_JSON,
				Description: "Route service traffic to pods with label keys and values matching this selector.",
				Transform:   transform.FromField("Description.Service.Spec.Selector"),
			},
			{
				Name:        "topology_keys",
				Type:        proto.ColumnType_JSON,
				Description: "A preference-order list of topology keys which implementations of services should use to preferentially sort endpoints when accessing this Service, it can not be used at the same time as externalTrafficPolicy=Local.",
				Transform:   transform.FromField("Description.Service.Spec.TopologyKeys"),
			},
			{
				Name:        "title",
				Type:        proto.ColumnType_STRING,
				Description: ColumnDescriptionTitle,
				Transform:   transform.FromField("Description.Service.Name"),
			},
			{
				Name:        "tags",
				Type:        proto.ColumnType_JSON,
				Description: ColumnDescriptionTags,
				Transform:   transform.From(transformServiceTags),
			},
		}),
	}
}

func transformServiceTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	obj := d.HydrateItem.(opengovernance.KubernetesService).Description.Service
	return mergeTags(obj.Labels, obj.Annotations), nil
}
