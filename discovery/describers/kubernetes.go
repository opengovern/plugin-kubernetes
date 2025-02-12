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
