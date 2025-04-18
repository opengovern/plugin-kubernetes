package kubernetes

import (
	"context"

	essdk "github.com/opengovern/og-util/pkg/opengovernance-es-sdk"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

// Plugin returns this plugin
func Plugin(ctx context.Context) *plugin.Plugin {
	p := &plugin.Plugin{
		Name: "steampipe-plugin-kubernetes",
		ConnectionConfigSchema: &plugin.ConnectionConfigSchema{
			NewInstance: essdk.ConfigInstance,
			Schema:      essdk.ConfigSchema(),
		},
		DefaultTransform: transform.FromCamel(),
		TableMap: map[string]*plugin.Table{
			"k8_cluster":                            tableKubernetesCluster(ctx),
			"k8_cluster_role":                       tableKubernetesClusterRole(ctx),
			"k8_cluster_role_binding":               tableKubernetesClusterRoleBinding(ctx),
			"k8_config_map":                         tableKubernetesConfigMap(ctx),
			"k8_cronjob":                            tableKubernetesCronJob(ctx),
			"k8_custom_resource":                    tableKubernetesCustomResource(ctx),
			"k8_custom_resource_definition":         tableKubernetesCustomResourceDefinition(ctx),
			"k8_daemonset":                          tableKubernetesDaemonset(ctx),
			"k8_deployment":                         tableKubernetesDeployment(ctx),
			"k8_endpoint_slice":                     tableKubernetesEndpointSlice(ctx),
			"k8_endpoints":                          tableKubernetesEndpoints(ctx),
			"k8_event":                              tableKubernetesEvent(ctx),
			"k8_horizontal_pod_autoscaler":          tableKubernetesHorizontalPodAutoscaler(ctx),
			"k8_ingress":                            tableKubernetesIngress(ctx),
			"k8_job":                                tableKubernetesJob(ctx),
			"k8_limit_range":                        tableKubernetesLimitRange(ctx),
			"k8_namespace":                          tableKubernetesNamespace(ctx),
			"k8_network_policy":                     tableKubernetesNetworkPolicy(ctx),
			"k8_node":                               tableKubernetesNode(ctx),
			"k8_persistent_volume_claim":            tableKubernetesPersistentVolumeClaim(ctx),
			"k8_persistent_volume":                  tableKubernetesPersistentVolume(ctx),
			"k8_pod":                                tableKubernetesPod(ctx),
			"k8_pod_disruption_budget":              tableKubernetesPDB(ctx),
			"k8_pod_template":                       tableKubernetesPodTemplate(ctx),
			"k8_replicaset":                         tableKubernetesReplicaSet(ctx),
			"k8_replication_controller":             tableKubernetesReplicaController(ctx),
			"k8_resource_quota":                     tableKubernetesResourceQuota(ctx),
			"k8_role":                               tableKubernetesRole(ctx),
			"k8_role_binding":                       tableKubernetesRoleBinding(ctx),
			"k8_secret":                             tableKubernetesSecret(ctx),
			"k8_service":                            tableKubernetesService(ctx),
			"k8_service_account":                    tableKubernetesServiceAccount(ctx),
			"k8_stateful_set":                       tableKubernetesStatefulSet(ctx),
			"k8_storage_class":                      tableKubernetesStorageClass(ctx),
			"kubernetes_cluster":                    tableKubernetesCluster(ctx),
			"kubernetes_cluster_role":               tableKubernetesClusterRole(ctx),
			"kubernetes_cluster_role_binding":       tableKubernetesClusterRoleBinding(ctx),
			"kubernetes_config_map":                 tableKubernetesConfigMap(ctx),
			"kubernetes_cronjob":                    tableKubernetesCronJob(ctx),
			"kubernetes_custom_resource":            tableKubernetesCustomResource(ctx),
			"kubernetes_custom_resource_definition": tableKubernetesCustomResourceDefinition(ctx),
			"kubernetes_daemonset":                  tableKubernetesDaemonset(ctx),
			"kubernetes_deployment":                 tableKubernetesDeployment(ctx),
			"kubernetes_endpoint_slice":             tableKubernetesEndpointSlice(ctx),
			"kubernetes_endpoints":                  tableKubernetesEndpoints(ctx),
			"kubernetes_event":                      tableKubernetesEvent(ctx),
			"kubernetes_horizontal_pod_autoscaler":  tableKubernetesHorizontalPodAutoscaler(ctx),
			"kubernetes_ingress":                    tableKubernetesIngress(ctx),
			"kubernetes_job":                        tableKubernetesJob(ctx),
			"kubernetes_limit_range":                tableKubernetesLimitRange(ctx),
			"kubernetes_namespace":                  tableKubernetesNamespace(ctx),
			"kubernetes_network_policy":             tableKubernetesNetworkPolicy(ctx),
			"kubernetes_node":                       tableKubernetesNode(ctx),
			"kubernetes_persistent_volume_claim":    tableKubernetesPersistentVolumeClaim(ctx),
			"kubernetes_persistent_volume":          tableKubernetesPersistentVolume(ctx),
			"kubernetes_pod":                        tableKubernetesPod(ctx),
			"kubernetes_pod_disruption_budget":      tableKubernetesPDB(ctx),
			"kubernetes_pod_template":               tableKubernetesPodTemplate(ctx),
			"kubernetes_replicaset":                 tableKubernetesReplicaSet(ctx),
			"kubernetes_replication_controller":     tableKubernetesReplicaController(ctx),
			"kubernetes_resource_quota":             tableKubernetesResourceQuota(ctx),
			"kubernetes_role":                       tableKubernetesRole(ctx),
			"kubernetes_role_binding":               tableKubernetesRoleBinding(ctx),
			"kubernetes_secret":                     tableKubernetesSecret(ctx),
			"kubernetes_service":                    tableKubernetesService(ctx),
			"kubernetes_service_account":            tableKubernetesServiceAccount(ctx),
			"kubernetes_stateful_set":               tableKubernetesStatefulSet(ctx),
			"kubernetes_storage_class":              tableKubernetesStorageClass(ctx),
		},
	}
	for key, table := range p.TableMap {
		if table == nil {
			continue
		}
		if table.Get != nil && table.Get.Hydrate == nil {
			delete(p.TableMap, key)
			continue
		}
		if table.List != nil && table.List.Hydrate == nil {
			delete(p.TableMap, key)
			continue
		}

		opengovernanceTable := false
		for _, col := range table.Columns {
			if col != nil && col.Name == "platform_integration_id" {
				opengovernanceTable = true
			}
		}

		if opengovernanceTable {
			if table.Get != nil {
				table.Get.KeyColumns = append(table.Get.KeyColumns, plugin.OptionalColumns([]string{"platform_integration_id", "platform_resource_id"})...)
			}

			if table.List != nil {
				table.List.KeyColumns = append(table.List.KeyColumns, plugin.OptionalColumns([]string{"platform_integration_id", "platform_resource_id"})...)
			}
		}
	}
	return p
}
