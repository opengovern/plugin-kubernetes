// artifact_dockerfile.go
package describers

import (
	"context"
	"fmt"
	"github.com/opengovern/og-describer-kubernetes/discovery/pkg/models"
	model "github.com/opengovern/og-describer-kubernetes/discovery/provider"
	"github.com/opengovern/og-describer-kubernetes/discovery/provider/helpers"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func KubernetesResources(ctx context.Context, client model.Client, extra string, stream *models.StreamSender) ([]models.Resource, error) {
	var allValues []models.Resource

	resources, err := GetKubernetesResources(client.KubeConfig)
	if err != nil {
		return nil, err
	}

	for _, r := range resources {
		resource := models.Resource{
			ID:          r.UID,
			Name:        r.ObjectName,
			Description: r,
		}

		if stream != nil {
			if err := (*stream)(resource); err != nil {
				return allValues, fmt.Errorf("error streaming resource: %w", err)
			}
		} else {
			allValues = append(allValues, resource)
		}
	}

	return allValues, nil
}

func KubernetesCluster(ctx context.Context, client model.Client, extra string, stream *models.StreamSender) ([]models.Resource, error) {
	var allValues []models.Resource

	cluster, err := DoDiscovery(client.KubeConfig)
	if err != nil {
		return nil, err
	}
	resource := models.Resource{
		ID:          fmt.Sprintf("cluster/%s", cluster.ContextName),
		Name:        cluster.ContextName,
		Description: cluster,
	}

	if stream != nil {
		if err := (*stream)(resource); err != nil {
			return allValues, fmt.Errorf("error streaming resource: %w", err)
		}
	} else {
		allValues = append(allValues, resource)
	}
	return allValues, nil
}

func KubernetesClusterRole(ctx context.Context, client model.Client, extra string, stream *models.StreamSender) ([]models.Resource, error) {
	var allValues []models.Resource

	clusterRoles, err := client.KubernetesClient.RbacV1().ClusterRoles().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	for _, clusterRole := range clusterRoles.Items {
		var resource models.Resource
		// We don't need to include the managed fields in the description, also it causes issues in elastic search mapping generation
		clusterRole.ManagedFields = nil
		resource = models.Resource{
			ID:   fmt.Sprintf("clusterrole/%s", clusterRole.Name),
			Name: clusterRole.Name,
			Description: model.KubernetesClusterRoleDescription{
				MetaObject: helpers.ConvertObjectMeta(&clusterRole.ObjectMeta),
				ClusterRole: helpers.ClusterRole{
					TypeMeta:        helpers.ConvertTypeMeta(clusterRole.TypeMeta),
					ObjectMeta:      helpers.ConvertObjectMeta(&clusterRole.ObjectMeta),
					Rules:           helpers.ConvertPolicyRules(clusterRole.Rules),
					AggregationRule: helpers.ConvertAggregationRule(clusterRole.AggregationRule),
				},
			},
		}

		if stream != nil {
			if err := (*stream)(resource); err != nil {
				return allValues, fmt.Errorf("error streaming resource: %w", err)
			}
		} else {
			allValues = append(allValues, resource)
		}
	}

	return allValues, nil
}

func KubernetesClusterRoleBinding(ctx context.Context, client model.Client, extra string, stream *models.StreamSender) ([]models.Resource, error) {
	var allValues []models.Resource

	clusterRoleBindings, err := client.KubernetesClient.RbacV1().ClusterRoleBindings().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	for _, clusterRoleBinding := range clusterRoleBindings.Items {
		var resource models.Resource
		// We don't need to include the managed fields in the description, also it causes issues in elastic search mapping generation
		clusterRoleBinding.ManagedFields = nil
		resource = models.Resource{
			ID:   fmt.Sprintf("clusterrolebinding/%s", clusterRoleBinding.Name),
			Name: clusterRoleBinding.Name,
			Description: model.KubernetesClusterRoleBindingDescription{
				MetaObject: helpers.ConvertObjectMeta(&clusterRoleBinding.ObjectMeta),
				ClusterRoleBinding: helpers.ClusterRoleBinding{
					TypeMeta:   helpers.ConvertTypeMeta(clusterRoleBinding.TypeMeta),
					ObjectMeta: helpers.ConvertObjectMeta(&clusterRoleBinding.ObjectMeta),
					Subjects:   helpers.ConvertSubjects(clusterRoleBinding.Subjects),
					RoleRef:    helpers.ConvertRoleRef(clusterRoleBinding.RoleRef),
				},
			},
		}

		if stream != nil {
			if err := (*stream)(resource); err != nil {
				return allValues, fmt.Errorf("error streaming resource: %w", err)
			}
		} else {
			allValues = append(allValues, resource)
		}
	}

	return allValues, nil
}

func KubernetesConfigMap(ctx context.Context, client model.Client, extra string, stream *models.StreamSender) ([]models.Resource, error) {
	var allValues []models.Resource

	configMaps, err := client.KubernetesClient.CoreV1().ConfigMaps("").List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	for _, configMap := range configMaps.Items {
		var resource models.Resource

		configMap.BinaryData = nil
		resource = models.Resource{
			ID:   fmt.Sprintf("configmap/%s/%s", configMap.Namespace, configMap.Name),
			Name: fmt.Sprintf("%s/%s", configMap.Namespace, configMap.Name),
			Description: model.KubernetesConfigMapDescription{
				MetaObject: helpers.ConvertObjectMeta(&configMap.ObjectMeta),
				ConfigMap: helpers.ConfigMap{
					TypeMeta:   helpers.ConvertTypeMeta(configMap.TypeMeta),
					ObjectMeta: helpers.ConvertObjectMeta(&configMap.ObjectMeta),
					Immutable:  configMap.Immutable,
					Data:       configMap.Data,
				},
			},
		}

		if stream != nil {
			if err := (*stream)(resource); err != nil {
				return allValues, fmt.Errorf("error streaming resource: %w", err)
			}
		} else {
			allValues = append(allValues, resource)
		}
	}

	return allValues, nil
}

func KubernetesCronJob(ctx context.Context, client model.Client, extra string, stream *models.StreamSender) ([]models.Resource, error) {
	var allValues []models.Resource

	cronJobs, err := client.KubernetesClient.BatchV1().CronJobs("").List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	for _, cronJob := range cronJobs.Items {
		var resource models.Resource
		// We don't need to include the managed fields in the description, also it causes issues in elastic search mapping generation
		cronJob.ManagedFields = nil
		resource = models.Resource{
			ID:   fmt.Sprintf("cronjob/%s/%s", cronJob.Namespace, cronJob.Name),
			Name: fmt.Sprintf("%s/%s", cronJob.Namespace, cronJob.Name),
			Description: model.KubernetesCronJobDescription{
				MetaObject: helpers.ConvertObjectMeta(&cronJob.ObjectMeta),
				CronJob:    helpers.ConvertCronJob(&cronJob),
			},
		}

		if stream != nil {
			if err := (*stream)(resource); err != nil {
				return allValues, fmt.Errorf("error streaming resource: %w", err)
			}
		} else {
			allValues = append(allValues, resource)
		}
	}

	return allValues, nil
}

func KubernetesCustomResource(ctx context.Context, client model.Client, extra string, stream *models.StreamSender) ([]models.Resource, error) {
	var allValues []models.Resource

	// list crds
	customResourceDefinitions, err := client.CrdsClient.ApiextensionsV1().CustomResourceDefinitions().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	for _, customResourceDefinition := range customResourceDefinitions.Items {
		// use dynamic client to list resources of this crd
		for _, version := range customResourceDefinition.Spec.Versions {
			dynamicClient := client.DynamicClient.Resource(schema.GroupVersionResource{
				Group:    customResourceDefinition.Spec.Group,
				Version:  version.Name,
				Resource: customResourceDefinition.Spec.Names.Plural,
			})

			resources, err := dynamicClient.List(ctx, metav1.ListOptions{})
			if err != nil {
				return nil, err
			}

			for _, item := range resources.Items {
				var resource models.Resource
				// We don't need to include the managed fields in the description, also it causes issues in elastic search mapping generation
				item.SetManagedFields(nil)
				resource = models.Resource{
					ID:   fmt.Sprintf("customresource/%s.%s/%s/%s", item.GetKind(), item.GetAPIVersion(), item.GetNamespace(), item.GetName()),
					Name: fmt.Sprintf("%s/%s/%s", customResourceDefinition.Name, item.GetNamespace(), item.GetName()),
					Description: model.KubernetesCustomResourceDescription{
						MetaObject: helpers.ObjectMeta{
							Name:                       item.GetName(),
							GenerateName:               item.GetGenerateName(),
							Namespace:                  item.GetNamespace(),
							SelfLink:                   item.GetSelfLink(),
							UID:                        item.GetUID(),
							ResourceVersion:            item.GetResourceVersion(),
							Generation:                 item.GetGeneration(),
							CreationTimestamp:          helpers.ConvertTime(item.GetCreationTimestamp()),
							DeletionTimestamp:          helpers.ConvertTimePtr(item.GetDeletionTimestamp()),
							DeletionGracePeriodSeconds: item.GetDeletionGracePeriodSeconds(),
							Labels:                     item.GetLabels(),
							Annotations:                item.GetAnnotations(),
							OwnerReferences:            helpers.ConvertOwnerReferences(item.GetOwnerReferences()),
							Finalizers:                 item.GetFinalizers(),
						},
						CustomResource:     item,
						FullyQualifiedName: fmt.Sprintf("%s.%s", item.GetKind(), item.GetAPIVersion()),
					},
				}

				if stream != nil {
					if err := (*stream)(resource); err != nil {
						return allValues, fmt.Errorf("error streaming resource: %w", err)
					}
				} else {
					allValues = append(allValues, resource)
				}
			}
		}
	}

	return allValues, nil
}

func KubernetesCustomResourceDefinition(ctx context.Context, client model.Client, extra string, stream *models.StreamSender) ([]models.Resource, error) {
	var allValues []models.Resource

	customResourceDefinitions, err := client.CrdsClient.ApiextensionsV1().CustomResourceDefinitions().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	for _, customResourceDefinition := range customResourceDefinitions.Items {
		var resource models.Resource
		// We don't need to include the managed fields in the description, also it causes issues in elastic search mapping generation
		customResourceDefinition.ManagedFields = nil
		resource = models.Resource{
			ID:   fmt.Sprintf("customresourcedefinition/%s", customResourceDefinition.Name),
			Name: customResourceDefinition.Name,
			Description: model.KubernetesCustomResourceDefinitionDescription{
				MetaObject:               helpers.ConvertObjectMeta(&customResourceDefinition.ObjectMeta),
				CustomResourceDefinition: helpers.ConvertCustomResourceDefinition(&customResourceDefinition),
			},
		}

		if stream != nil {
			if err := (*stream)(resource); err != nil {
				return allValues, fmt.Errorf("error streaming resource: %w", err)
			}
		} else {
			allValues = append(allValues, resource)
		}
	}

	return allValues, nil
}

func KubernetesDaemonSet(ctx context.Context, client model.Client, extra string, stream *models.StreamSender) ([]models.Resource, error) {
	var allValues []models.Resource

	daemonSets, err := client.KubernetesClient.AppsV1().DaemonSets("").List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	for _, daemonSet := range daemonSets.Items {
		var resource models.Resource
		// We don't need to include the managed fields in the description, also it causes issues in elastic search mapping generation
		daemonSet.ManagedFields = nil
		labelSelectorString := ""
		ss, err := metav1.LabelSelectorAsSelector(daemonSet.Spec.Selector)
		if err == nil {
			labelSelectorString = ss.String()
		}
		resource = models.Resource{
			ID:   fmt.Sprintf("daemonset/%s/%s", daemonSet.Namespace, daemonSet.Name),
			Name: fmt.Sprintf("%s/%s", daemonSet.Namespace, daemonSet.Name),
			Description: model.KubernetesDaemonSetDescription{
				MetaObject:          helpers.ConvertObjectMeta(&daemonSet.ObjectMeta),
				DaemonSet:           helpers.ConvertDaemonSet(&daemonSet),
				LabelSelectorString: labelSelectorString,
			},
		}

		if stream != nil {
			if err := (*stream)(resource); err != nil {
				return allValues, fmt.Errorf("error streaming resource: %w", err)
			}
		} else {
			allValues = append(allValues, resource)
		}
	}

	return allValues, nil
}

func KubernetesDeployment(ctx context.Context, client model.Client, extra string, stream *models.StreamSender) ([]models.Resource, error) {
	var allValues []models.Resource

	deployments, err := client.KubernetesClient.AppsV1().Deployments("").List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	for _, deployment := range deployments.Items {
		var resource models.Resource
		// We don't need to include the managed fields in the description, also it causes issues in elastic search mapping generation
		deployment.ManagedFields = nil
		labelSelectorString := ""
		ss, err := metav1.LabelSelectorAsSelector(deployment.Spec.Selector)
		if err == nil {
			labelSelectorString = ss.String()
		}
		resource = models.Resource{
			ID:   fmt.Sprintf("deployment/%s/%s", deployment.Namespace, deployment.Name),
			Name: fmt.Sprintf("%s/%s", deployment.Namespace, deployment.Name),
			Description: model.KubernetesDeploymentDescription{
				MetaObject:          helpers.ConvertObjectMeta(&deployment.ObjectMeta),
				Deployment:          helpers.ConvertDeployment(&deployment),
				LabelSelectorString: labelSelectorString,
			},
		}

		if stream != nil {
			if err := (*stream)(resource); err != nil {
				return allValues, fmt.Errorf("error streaming resource: %w", err)
			}
		} else {
			allValues = append(allValues, resource)
		}
	}

	return allValues, nil
}

func KubernetesEndpointSlice(ctx context.Context, client model.Client, extra string, stream *models.StreamSender) ([]models.Resource, error) {
	var allValues []models.Resource

	endpointSlices, err := client.KubernetesClient.DiscoveryV1().EndpointSlices("").List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	for _, endpointSlice := range endpointSlices.Items {
		var resource models.Resource
		// We don't need to include the managed fields in the description, also it causes issues in elastic search mapping generation
		endpointSlice.ManagedFields = nil
		resource = models.Resource{
			ID:   fmt.Sprintf("endpointslice/%s/%s", endpointSlice.Namespace, endpointSlice.Name),
			Name: fmt.Sprintf("%s/%s", endpointSlice.Namespace, endpointSlice.Name),
			Description: model.KubernetesEndpointSliceDescription{
				MetaObject:    helpers.ConvertObjectMeta(&endpointSlice.ObjectMeta),
				EndpointSlice: helpers.ConvertEndpointSlice(&endpointSlice),
			},
		}
		if stream != nil {
			if err := (*stream)(resource); err != nil {
				return allValues, fmt.Errorf("error streaming resource: %w", err)
			}
		} else {
			allValues = append(allValues, resource)
		}
	}

	return allValues, nil
}

func KubernetesEndpoint(ctx context.Context, client model.Client, extra string, stream *models.StreamSender) ([]models.Resource, error) {
	var allValues []models.Resource

	endpoints, err := client.KubernetesClient.CoreV1().Endpoints("").List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	for _, endpoint := range endpoints.Items {
		var resource models.Resource
		// We don't need to include the managed fields in the description, also it causes issues in elastic search mapping generation
		endpoint.ManagedFields = nil
		resource = models.Resource{
			ID:   fmt.Sprintf("endpoint/%s/%s", endpoint.Namespace, endpoint.Name),
			Name: fmt.Sprintf("%s/%s", endpoint.Namespace, endpoint.Name),
			Description: model.KubernetesEndpointDescription{
				MetaObject: helpers.ConvertObjectMeta(&endpoint.ObjectMeta),
				Endpoint:   helpers.ConvertEndpoints(&endpoint),
			},
		}
		if stream != nil {
			if err := (*stream)(resource); err != nil {
				return allValues, fmt.Errorf("error streaming resource: %w", err)
			}
		} else {
			allValues = append(allValues, resource)
		}
	}

	return allValues, nil
}

func KubernetesEvent(ctx context.Context, client model.Client, extra string, stream *models.StreamSender) ([]models.Resource, error) {
	var allValues []models.Resource

	events, err := client.KubernetesClient.CoreV1().Events("").List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	for _, event := range events.Items {
		var resource models.Resource
		// We don't need to include the managed fields in the description, also it causes issues in elastic search mapping generation
		event.ManagedFields = nil
		resource = models.Resource{
			ID:   fmt.Sprintf("event/%s/%s", event.Namespace, event.Name),
			Name: fmt.Sprintf("%s/%s", event.Namespace, event.Name),
			Description: model.KubernetesEventDescription{
				MetaObject: helpers.ConvertObjectMeta(&event.ObjectMeta),
				Event:      helpers.ConvertEvent(&event),
			},
		}
		if stream != nil {
			if err := (*stream)(resource); err != nil {
				return allValues, fmt.Errorf("error streaming resource: %w", err)
			}
		} else {
			allValues = append(allValues, resource)
		}
	}

	return allValues, nil
}

func KubernetesHorizontalPodAutoscaler(ctx context.Context, client model.Client, extra string, stream *models.StreamSender) ([]models.Resource, error) {
	var allValues []models.Resource

	horizontalPodAutoscalers, err := client.KubernetesClient.AutoscalingV1().HorizontalPodAutoscalers("").List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	for _, horizontalPodAutoscaler := range horizontalPodAutoscalers.Items {
		var resource models.Resource
		// We don't need to include the managed fields in the description, also it causes issues in elastic search mapping generation
		horizontalPodAutoscaler.ManagedFields = nil
		resource = models.Resource{
			ID:   fmt.Sprintf("horizontalpodautoscaler/%s/%s", horizontalPodAutoscaler.Namespace, horizontalPodAutoscaler.Name),
			Name: fmt.Sprintf("%s/%s", horizontalPodAutoscaler.Namespace, horizontalPodAutoscaler.Name),
			Description: model.KubernetesHorizontalPodAutoscalerDescription{
				MetaObject:              helpers.ConvertObjectMeta(&horizontalPodAutoscaler.ObjectMeta),
				HorizontalPodAutoscaler: helpers.ConvertHorizontalPodAutoscaler(&horizontalPodAutoscaler),
			},
		}
		if stream != nil {
			if err := (*stream)(resource); err != nil {
				return allValues, fmt.Errorf("error streaming resource: %w", err)
			}
		} else {
			allValues = append(allValues, resource)
		}
	}

	return allValues, nil
}

func KubernetesIngress(ctx context.Context, client model.Client, extra string, stream *models.StreamSender) ([]models.Resource, error) {
	var allValues []models.Resource

	ingresses, err := client.KubernetesClient.NetworkingV1().Ingresses("").List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	for _, ingress := range ingresses.Items {
		var resource models.Resource
		// We don't need to include the managed fields in the description, also it causes issues in elastic search mapping generation
		ingress.ManagedFields = nil
		resource = models.Resource{
			ID:   fmt.Sprintf("ingress/%s/%s", ingress.Namespace, ingress.Name),
			Name: fmt.Sprintf("%s/%s", ingress.Namespace, ingress.Name),
			Description: model.KubernetesIngressDescription{
				MetaObject: helpers.ConvertObjectMeta(&ingress.ObjectMeta),
				Ingress:    helpers.ConvertIngress(&ingress),
			},
		}
		if stream != nil {
			if err := (*stream)(resource); err != nil {
				return allValues, fmt.Errorf("error streaming resource: %w", err)
			}
		} else {
			allValues = append(allValues, resource)
		}
	}

	return allValues, nil
}

func KubernetesJob(ctx context.Context, client model.Client, extra string, stream *models.StreamSender) ([]models.Resource, error) {
	var allValues []models.Resource

	jobs, err := client.KubernetesClient.BatchV1().Jobs("").List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	for _, job := range jobs.Items {
		var resource models.Resource
		// We don't need to include the managed fields in the description, also it causes issues in elastic search mapping
		job.ManagedFields = nil
		labelSelectorString := ""
		ss, err := metav1.LabelSelectorAsSelector(job.Spec.Selector)
		if err == nil {
			labelSelectorString = ss.String()
		}
		resource = models.Resource{
			ID:   fmt.Sprintf("job/%s/%s", job.Namespace, job.Name),
			Name: fmt.Sprintf("%s/%s", job.Namespace, job.Name),
			Description: model.KubernetesJobDescription{
				MetaObject:          helpers.ConvertObjectMeta(&job.ObjectMeta),
				Job:                 job,
				LabelSelectorString: labelSelectorString,
			},
		}
		if stream != nil {
			if err := (*stream)(resource); err != nil {
				return allValues, fmt.Errorf("error streaming resource: %w", err)
			}
		} else {
			allValues = append(allValues, resource)
		}
	}

	return allValues, nil
}

func KubernetesLimitRange(ctx context.Context, client model.Client, extra string, stream *models.StreamSender) ([]models.Resource, error) {
	var allValues []models.Resource

	limitRanges, err := client.KubernetesClient.CoreV1().LimitRanges("").List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	for _, limitRange := range limitRanges.Items {
		var resource models.Resource
		// We don't need to include the managed fields in the description, also it causes issues in elastic search mapping generation
		limitRange.ManagedFields = nil
		resource = models.Resource{
			ID:   fmt.Sprintf("limitrange/%s/%s", limitRange.Namespace, limitRange.Name),
			Name: fmt.Sprintf("%s/%s", limitRange.Namespace, limitRange.Name),
			Description: model.KubernetesLimitRangeDescription{
				MetaObject: helpers.ConvertObjectMeta(&limitRange.ObjectMeta),
				LimitRange: helpers.ConvertLimitRange(&limitRange),
			},
		}
		if stream != nil {
			if err := (*stream)(resource); err != nil {
				return allValues, fmt.Errorf("error streaming resource: %w", err)
			}
		} else {
			allValues = append(allValues, resource)
		}
	}

	return allValues, nil
}

func KubernetesNamespace(ctx context.Context, client model.Client, extra string, stream *models.StreamSender) ([]models.Resource, error) {
	var allValues []models.Resource

	namespaces, err := client.KubernetesClient.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	for _, namespace := range namespaces.Items {
		var resource models.Resource
		// We don't need to include the managed fields in the description, also it causes issues in elastic search mapping generation
		namespace.ManagedFields = nil
		resource = models.Resource{
			ID:   fmt.Sprintf("namespace/%s", namespace.Name),
			Name: namespace.Name,
			Description: model.KubernetesNamespaceDescription{
				MetaObject: helpers.ConvertObjectMeta(&namespace.ObjectMeta),
				Namespace:  helpers.ConvertNamespace(&namespace),
			},
		}
		if stream != nil {
			if err := (*stream)(resource); err != nil {
				return allValues, fmt.Errorf("error streaming resource: %w", err)
			}
		} else {
			allValues = append(allValues, resource)
		}
	}

	return allValues, nil
}

func KubernetesNetworkPolicy(ctx context.Context, client model.Client, extra string, stream *models.StreamSender) ([]models.Resource, error) {
	var allValues []models.Resource

	networkPolicies, err := client.KubernetesClient.NetworkingV1().NetworkPolicies("").List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	for _, networkPolicy := range networkPolicies.Items {
		var resource models.Resource
		// We don't need to include the managed fields in the description, also it causes issues in elastic search mapping generation
		networkPolicy.ManagedFields = nil
		resource = models.Resource{
			ID:   fmt.Sprintf("networkpolicy/%s/%s", networkPolicy.Namespace, networkPolicy.Name),
			Name: fmt.Sprintf("%s/%s", networkPolicy.Namespace, networkPolicy.Name),
			Description: model.KubernetesNetworkPolicyDescription{
				MetaObject:    helpers.ConvertObjectMeta(&networkPolicy.ObjectMeta),
				NetworkPolicy: helpers.ConvertNetworkPolicy(&networkPolicy),
			},
		}
		if stream != nil {
			if err := (*stream)(resource); err != nil {
				return allValues, fmt.Errorf("error streaming resource: %w", err)
			}
		} else {
			allValues = append(allValues, resource)
		}
	}

	return allValues, nil
}

func KubernetesNode(ctx context.Context, client model.Client, extra string, stream *models.StreamSender) ([]models.Resource, error) {
	var allValues []models.Resource

	// list nodes
	nodes, err := client.KubernetesClient.CoreV1().Nodes().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	for _, node := range nodes.Items {
		var resource models.Resource
		// We don't need to include the managed fields in the description, also it causes issues in elastic search mapping generation
		node.ManagedFields = nil
		resource = models.Resource{
			ID:   fmt.Sprintf("node/%s", node.Name),
			Name: node.Name,
			Description: model.KubernetesNodeDescription{
				MetaObject: helpers.ConvertObjectMeta(&node.ObjectMeta),
				Node:       helpers.ConvertNode(&node),
			},
		}

		if stream != nil {
			if err := (*stream)(resource); err != nil {
				return allValues, fmt.Errorf("error streaming resource: %w", err)
			}
		} else {
			allValues = append(allValues, resource)
		}
	}

	return allValues, nil
}

func KubernetesPersistentVolume(ctx context.Context, client model.Client, extra string, stream *models.StreamSender) ([]models.Resource, error) {
	var allValues []models.Resource

	// list persistent volumes
	pvs, err := client.KubernetesClient.CoreV1().PersistentVolumes().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	for _, pv := range pvs.Items {
		var resource models.Resource
		// We don't need to include the managed fields in the description, also it causes issues in elastic search mapping generation
		pv.ManagedFields = nil
		resource = models.Resource{
			ID:   fmt.Sprintf("persistentvolume/%s", pv.Name),
			Name: pv.Name,
			Description: model.KubernetesPersistentVolumeDescription{
				MetaObject: helpers.ConvertObjectMeta(&pv.ObjectMeta),
				PV:         helpers.ConvertPersistentVolume(&pv),
			},
		}

		if stream != nil {
			if err := (*stream)(resource); err != nil {
				return allValues, fmt.Errorf("error streaming resource: %w", err)
			}
		} else {
			allValues = append(allValues, resource)
		}
	}

	return allValues, nil
}

func KubernetesPersistentVolumeClaim(ctx context.Context, client model.Client, extra string, stream *models.StreamSender) ([]models.Resource, error) {
	var allValues []models.Resource

	pvcs, err := client.KubernetesClient.CoreV1().PersistentVolumeClaims("").List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	for _, pvc := range pvcs.Items {
		var resource models.Resource
		// We don't need to include the managed fields in the description, also it causes issues in elastic search mapping generation
		pvc.ManagedFields = nil
		resource = models.Resource{
			ID:   fmt.Sprintf("persistentvolumeclaim/%s/%s", pvc.Namespace, pvc.Name),
			Name: fmt.Sprintf("%s/%s", pvc.Namespace, pvc.Name),
			Description: model.KubernetesPersistentVolumeClaimDescription{
				MetaObject: helpers.ConvertObjectMeta(&pvc.ObjectMeta),
				PVC:        helpers.ConvertPersistentVolumeClaim(&pvc),
			},
		}

		if stream != nil {
			if err := (*stream)(resource); err != nil {
				return allValues, fmt.Errorf("error streaming resource: %w", err)
			}
		} else {
			allValues = append(allValues, resource)
		}
	}

	return allValues, nil
}

func KubernetesPod(ctx context.Context, client model.Client, extra string, stream *models.StreamSender) ([]models.Resource, error) {
	var allValues []models.Resource

	pods, err := client.KubernetesClient.CoreV1().Pods("").List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	for _, pod := range pods.Items {
		var resource models.Resource
		// We don't need to include the managed fields in the description, also it causes issues in elastic search mapping generation
		pod.ManagedFields = nil
		resource = models.Resource{
			ID:   fmt.Sprintf("pod/%s/%s", pod.Namespace, pod.Name),
			Name: fmt.Sprintf("%s/%s", pod.Namespace, pod.Name),
			Description: model.KubernetesPodDescription{
				MetaObject: helpers.ConvertObjectMeta(&pod.ObjectMeta),
				Pod:        helpers.ConvertPod(&pod),
			},
		}

		if stream != nil {
			if err := (*stream)(resource); err != nil {
				return allValues, fmt.Errorf("error streaming resource: %w", err)
			}
		} else {
			allValues = append(allValues, resource)
		}
	}

	return allValues, nil
}

func KubernetesPodDisruptionBudget(ctx context.Context, client model.Client, extra string, stream *models.StreamSender) ([]models.Resource, error) {
	var allValues []models.Resource

	podDisruptionBudgets, err := client.KubernetesClient.PolicyV1().PodDisruptionBudgets("").List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	for _, podDisruptionBudget := range podDisruptionBudgets.Items {
		var resource models.Resource
		// We don't need to include the managed fields in the description, also it causes issues in elastic search mapping
		podDisruptionBudget.ManagedFields = nil
		resource = models.Resource{
			ID:   fmt.Sprintf("poddisruptionbudget/%s/%s", podDisruptionBudget.Namespace, podDisruptionBudget.Name),
			Name: fmt.Sprintf("%s/%s", podDisruptionBudget.Namespace, podDisruptionBudget.Name),
			Description: model.KubernetesPodDisruptionBudgetDescription{
				MetaObject:          helpers.ConvertObjectMeta(&podDisruptionBudget.ObjectMeta),
				PodDisruptionBudget: helpers.ConvertPodDisruptionBudget(&podDisruptionBudget),
			},
		}

		if stream != nil {
			if err := (*stream)(resource); err != nil {
				return allValues, fmt.Errorf("error streaming resource: %w", err)
			}
		} else {
			allValues = append(allValues, resource)
		}
	}

	return allValues, nil
}

func KubernetesPodTemplate(ctx context.Context, client model.Client, extra string, stream *models.StreamSender) ([]models.Resource, error) {
	var allValues []models.Resource

	podTemplates, err := client.KubernetesClient.CoreV1().PodTemplates("").List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	for _, podTemplate := range podTemplates.Items {
		var resource models.Resource
		// We don't need to include the managed fields in the description, also it causes issues in elastic search mapping
		podTemplate.ManagedFields = nil
		resource = models.Resource{
			ID:   fmt.Sprintf("podtemplate/%s/%s", podTemplate.Namespace, podTemplate.Name),
			Name: fmt.Sprintf("%s/%s", podTemplate.Namespace, podTemplate.Name),
			Description: model.KubernetesPodTemplateDescription{
				MetaObject:  helpers.ConvertObjectMeta(&podTemplate.ObjectMeta),
				PodTemplate: helpers.ConvertPodTemplate(&podTemplate),
			},
		}

		if stream != nil {
			if err := (*stream)(resource); err != nil {
				return allValues, fmt.Errorf("error streaming resource: %w", err)
			}
		} else {
			allValues = append(allValues, resource)
		}
	}

	return allValues, nil
}

func KubernetesReplicaSet(ctx context.Context, client model.Client, extra string, stream *models.StreamSender) ([]models.Resource, error) {
	var allValues []models.Resource

	replicaSets, err := client.KubernetesClient.AppsV1().ReplicaSets("").List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	for _, replicaSet := range replicaSets.Items {
		var resource models.Resource
		// We don't need to include the managed fields in the description, also it causes issues in elastic search mapping
		replicaSet.ManagedFields = nil
		labelSelectorString := ""
		ss, err := metav1.LabelSelectorAsSelector(replicaSet.Spec.Selector)
		if err == nil {
			labelSelectorString = ss.String()
		}
		resource = models.Resource{
			ID:   fmt.Sprintf("replicaset/%s/%s", replicaSet.Namespace, replicaSet.Name),
			Name: fmt.Sprintf("%s/%s", replicaSet.Namespace, replicaSet.Name),
			Description: model.KubernetesReplicaSetDescription{
				MetaObject:          helpers.ConvertObjectMeta(&replicaSet.ObjectMeta),
				ReplicaSet:          helpers.ConvertReplicaSet(&replicaSet),
				LabelSelectorString: labelSelectorString,
			},
		}

		if stream != nil {
			if err := (*stream)(resource); err != nil {
				return allValues, fmt.Errorf("error streaming resource: %w", err)
			}
		} else {
			allValues = append(allValues, resource)
		}
	}

	return allValues, nil
}

func KubernetesReplicationController(ctx context.Context, client model.Client, extra string, stream *models.StreamSender) ([]models.Resource, error) {
	var allValues []models.Resource

	replicationControllers, err := client.KubernetesClient.CoreV1().ReplicationControllers("").List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	for _, replicationController := range replicationControllers.Items {
		var resource models.Resource
		// We don't need to include the managed fields in the description, also it causes issues in elastic search mapping
		replicationController.ManagedFields = nil
		resource = models.Resource{
			ID:   fmt.Sprintf("replicationcontroller/%s/%s", replicationController.Namespace, replicationController.Name),
			Name: fmt.Sprintf("%s/%s", replicationController.Namespace, replicationController.Name),
			Description: model.KubernetesReplicationControllerDescription{
				MetaObject:            helpers.ConvertObjectMeta(&replicationController.ObjectMeta),
				ReplicationController: helpers.ConvertReplicationController(&replicationController),
			},
		}

		if stream != nil {
			if err := (*stream)(resource); err != nil {
				return allValues, fmt.Errorf("error streaming resource: %w", err)
			}
		} else {
			allValues = append(allValues, resource)
		}
	}

	return allValues, nil
}

func KubernetesResourceQuota(ctx context.Context, client model.Client, extra string, stream *models.StreamSender) ([]models.Resource, error) {
	var allValues []models.Resource

	resourceQuotas, err := client.KubernetesClient.CoreV1().ResourceQuotas("").List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	for _, resourceQuota := range resourceQuotas.Items {
		var resource models.Resource
		// We don't need to include the managed fields in the description, also it causes issues in elastic search mapping
		resourceQuota.ManagedFields = nil
		resource = models.Resource{
			ID:   fmt.Sprintf("resourcequota/%s/%s", resourceQuota.Namespace, resourceQuota.Name),
			Name: fmt.Sprintf("%s/%s", resourceQuota.Namespace, resourceQuota.Name),
			Description: model.KubernetesResourceQuotaDescription{
				MetaObject:    helpers.ConvertObjectMeta(&resourceQuota.ObjectMeta),
				ResourceQuota: helpers.ConvertResourceQuota(&resourceQuota),
			},
		}

		if stream != nil {
			if err := (*stream)(resource); err != nil {
				return allValues, fmt.Errorf("error streaming resource: %w", err)
			}
		} else {
			allValues = append(allValues, resource)
		}
	}

	return allValues, nil
}

func KubernetesRole(ctx context.Context, client model.Client, extra string, stream *models.StreamSender) ([]models.Resource, error) {
	var allValues []models.Resource

	roles, err := client.KubernetesClient.RbacV1().Roles("").List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	for _, role := range roles.Items {
		var resource models.Resource
		// We don't need to include the managed fields in the description, also it causes issues in elastic search
		role.ManagedFields = nil
		resource = models.Resource{
			ID:   fmt.Sprintf("role/%s/%s", role.Namespace, role.Name),
			Name: fmt.Sprintf("%s/%s", role.Namespace, role.Name),
			Description: model.KubernetesRoleDescription{
				MetaObject: helpers.ConvertObjectMeta(&role.ObjectMeta),
				Role:       helpers.ConvertRole(&role),
			},
		}

		if stream != nil {
			if err := (*stream)(resource); err != nil {
				return allValues, fmt.Errorf("error streaming resource: %w", err)
			}
		} else {
			allValues = append(allValues, resource)
		}
	}

	return allValues, nil
}

func KubernetesRoleBinding(ctx context.Context, client model.Client, extra string, stream *models.StreamSender) ([]models.Resource, error) {
	var allValues []models.Resource

	roleBindings, err := client.KubernetesClient.RbacV1().RoleBindings("").List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	for _, roleBinding := range roleBindings.Items {
		var resource models.Resource
		// We don't need to include the managed fields in the description, also it causes issues in elastic search
		roleBinding.ManagedFields = nil
		resource = models.Resource{
			ID:   fmt.Sprintf("rolebinding/%s/%s", roleBinding.Namespace, roleBinding.Name),
			Name: fmt.Sprintf("%s/%s", roleBinding.Namespace, roleBinding.Name),
			Description: model.KubernetesRoleBindingDescription{
				MetaObject:  helpers.ConvertObjectMeta(&roleBinding.ObjectMeta),
				RoleBinding: helpers.ConvertRoleBinding(&roleBinding),
			},
		}

		if stream != nil {
			if err := (*stream)(resource); err != nil {
				return allValues, fmt.Errorf("error streaming resource: %w", err)
			}
		} else {
			allValues = append(allValues, resource)
		}
	}

	return allValues, nil
}

func KubernetesSecret(ctx context.Context, client model.Client, extra string, stream *models.StreamSender) ([]models.Resource, error) {
	var allValues []models.Resource

	secrets, err := client.KubernetesClient.CoreV1().Secrets("").List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	for _, secret := range secrets.Items {
		var resource models.Resource
		resource = models.Resource{
			ID:   fmt.Sprintf("secret/%s/%s", secret.Namespace, secret.Name),
			Name: fmt.Sprintf("%s/%s", secret.Namespace, secret.Name),
			Description: model.KubernetesSecretDescription{
				MetaObject: helpers.ConvertObjectMeta(&secret.ObjectMeta),
				Secret:     helpers.ConvertSecret(&secret),
			},
		}

		if stream != nil {
			if err := (*stream)(resource); err != nil {
				return allValues, fmt.Errorf("error streaming resource: %w", err)
			}
		} else {
			allValues = append(allValues, resource)
		}
	}

	return allValues, nil
}

func KubernetesService(ctx context.Context, client model.Client, extra string, stream *models.StreamSender) ([]models.Resource, error) {
	var allValues []models.Resource

	services, err := client.KubernetesClient.CoreV1().Services("").List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	for _, service := range services.Items {
		var resource models.Resource
		// We don't need to include the managed fields in the description, also it causes issues in elastic search mapping generation
		service.ManagedFields = nil
		resource = models.Resource{
			ID:   fmt.Sprintf("service/%s/%s", service.Namespace, service.Name),
			Name: fmt.Sprintf("%s/%s", service.Namespace, service.Name),
			Description: model.KubernetesServiceDescription{
				MetaObject: helpers.ConvertObjectMeta(&service.ObjectMeta),
				Service:    helpers.ConvertService(&service),
			},
		}

		if stream != nil {
			if err := (*stream)(resource); err != nil {
				return allValues, fmt.Errorf("error streaming resource: %w", err)
			}
		} else {
			allValues = append(allValues, resource)
		}
	}

	return allValues, nil
}

func KubernetesServiceAccount(ctx context.Context, client model.Client, extra string, stream *models.StreamSender) ([]models.Resource, error) {
	var allValues []models.Resource

	serviceAccounts, err := client.KubernetesClient.CoreV1().ServiceAccounts("").List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	for _, serviceAccount := range serviceAccounts.Items {
		var resource models.Resource
		// We don't need to include the managed fields in the description, also it causes issues in elastic search mapping generation
		serviceAccount.ManagedFields = nil
		resource = models.Resource{
			ID:   fmt.Sprintf("serviceaccount/%s/%s", serviceAccount.Namespace, serviceAccount.Name),
			Name: fmt.Sprintf("%s/%s", serviceAccount.Namespace, serviceAccount.Name),
			Description: model.KubernetesServiceAccountDescription{
				MetaObject:     helpers.ConvertObjectMeta(&serviceAccount.ObjectMeta),
				ServiceAccount: helpers.ConvertServiceAccount(&serviceAccount),
			},
		}

		if stream != nil {
			if err := (*stream)(resource); err != nil {
				return allValues, fmt.Errorf("error streaming resource: %w", err)
			}
		} else {
			allValues = append(allValues, resource)
		}
	}

	return allValues, nil
}

func KubernetesStatefulSet(ctx context.Context, client model.Client, extra string, stream *models.StreamSender) ([]models.Resource, error) {
	var allValues []models.Resource

	statefulSets, err := client.KubernetesClient.AppsV1().StatefulSets("").List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	for _, statefulSet := range statefulSets.Items {
		var resource models.Resource
		// We don't need to include the managed fields in the description, also it causes issues in elastic search mapping generation
		statefulSet.ManagedFields = nil
		resource = models.Resource{
			ID:   fmt.Sprintf("statefulset/%s/%s", statefulSet.Namespace, statefulSet.Name),
			Name: fmt.Sprintf("%s/%s", statefulSet.Namespace, statefulSet.Name),
			Description: model.KubernetesStatefulSetDescription{
				MetaObject:  helpers.ConvertObjectMeta(&statefulSet.ObjectMeta),
				StatefulSet: helpers.ConvertStatefulSet(&statefulSet),
			},
		}

		if stream != nil {
			if err := (*stream)(resource); err != nil {
				return allValues, fmt.Errorf("error streaming resource: %w", err)
			}
		} else {
			allValues = append(allValues, resource)
		}
	}

	return allValues, nil
}

func KubernetesStorageClass(ctx context.Context, client model.Client, extra string, stream *models.StreamSender) ([]models.Resource, error) {
	var allValues []models.Resource

	storageClasses, err := client.KubernetesClient.StorageV1().StorageClasses().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	for _, storageClass := range storageClasses.Items {
		var resource models.Resource
		// We don't need to include the managed fields in the description, also it causes issues in elastic search
		storageClass.ManagedFields = nil
		resource = models.Resource{
			ID:   fmt.Sprintf("storageclass/%s", storageClass.Name),
			Name: storageClass.Name,
			Description: model.KubernetesStorageClassDescription{
				MetaObject:   helpers.ConvertObjectMeta(&storageClass.ObjectMeta),
				StorageClass: helpers.ConvertStorageClass(&storageClass),
			},
		}

		if stream != nil {
			if err := (*stream)(resource); err != nil {
				return allValues, fmt.Errorf("error streaming resource: %w", err)
			}
		} else {
			allValues = append(allValues, resource)
		}
	}

	return allValues, nil
}
