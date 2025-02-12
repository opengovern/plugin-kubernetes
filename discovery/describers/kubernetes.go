// artifact_dockerfile.go
package describers

import (
	"context"
	"fmt"
	"github.com/opengovern/og-describer-kubernetes/discovery/pkg/models"
	model "github.com/opengovern/og-describer-kubernetes/discovery/provider"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func KubernetesNode(ctx context.Context, client model.Client, extra string, stream *models.StreamSender) ([]models.Resource, error) {
	var allValues []models.Resource

	// list nodes
	nodes, err := client.KubernetesClient.CoreV1().Nodes().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	for _, node := range nodes.Items {
		var resource models.Resource

		resource = models.Resource{
			ID:   fmt.Sprintf("node/%s", node.Name),
			Name: node.Name,
			Description: model.KubernetesNodeDescription{
				MetaObject: node.ObjectMeta,
				Node:       node,
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

		resource = models.Resource{
			ID:   fmt.Sprintf("persistentvolume/%s", pv.Name),
			Name: pv.Name,
			Description: model.KubernetesPersistentVolumeDescription{
				MetaObject: pv.ObjectMeta,
				PV:         pv,
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

		resource = models.Resource{
			ID:   fmt.Sprintf("persistentvolumeclaim/%s/%s", pvc.Namespace, pvc.Name),
			Name: fmt.Sprintf("%s/%s", pvc.Namespace, pvc.Name),
			Description: model.KubernetesPersistentVolumeClaimDescription{
				MetaObject: pvc.ObjectMeta,
				PVC:        pvc,
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

		resource = models.Resource{
			ID:   fmt.Sprintf("pod/%s/%s", pod.Namespace, pod.Name),
			Name: fmt.Sprintf("%s/%s", pod.Namespace, pod.Name),
			Description: model.KubernetesPodDescription{
				MetaObject: pod.ObjectMeta,
				Pod:        pod,
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

		resource = models.Resource{
			ID:   fmt.Sprintf("service/%s/%s", service.Namespace, service.Name),
			Name: fmt.Sprintf("%s/%s", service.Namespace, service.Name),
			Description: model.KubernetesServiceDescription{
				MetaObject: service.ObjectMeta,
				Service:    service,
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
		// Do not include the data in the secret
		secret.Data = nil
		secret.StringData = nil
		resource = models.Resource{
			ID:   fmt.Sprintf("secret/%s/%s", secret.Namespace, secret.Name),
			Name: fmt.Sprintf("%s/%s", secret.Namespace, secret.Name),
			Description: model.KubernetesSecretDescription{
				MetaObject: secret.ObjectMeta,
				Secret:     secret,
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
