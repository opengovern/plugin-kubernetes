package provider

import (
	helmclient "github.com/mittwald/go-helm-client"
	model "github.com/opengovern/og-describer-kubernetes/discovery/pkg/models"
	"github.com/opengovern/og-util/pkg/describe/enums"
	"golang.org/x/net/context"
	apiextensionsclientset "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type Client struct {
	KubernetesClient *kubernetes.Clientset
	CrdsClient       *apiextensionsclientset.Clientset
	DynamicClient    *dynamic.DynamicClient
	HelmClient       helmclient.Client
	KubeConfig       string
}

func DescribeByIntegration(describe func(context.Context, Client, string, *model.StreamSender) ([]model.Resource, error)) model.ResourceDescriber {
	return func(ctx context.Context, cfg model.IntegrationCredentials, triggerType enums.DescribeTriggerType, additionalParameters map[string]string, stream *model.StreamSender) ([]model.Resource, error) {
		var values []model.Resource

		config, err := clientcmd.RESTConfigFromKubeConfig([]byte(cfg.KubeConfig))
		if err != nil {
			return nil, err
		}

		kubernetesClient, err := kubernetes.NewForConfig(config)
		if err != nil {
			return nil, err
		}

		helmClient, err := helmclient.NewClientFromRestConf(&helmclient.RestConfClientOptions{
			Options:    &helmclient.Options{},
			RestConfig: config,
		})
		if err != nil {
			return nil, err
		}

		crdClient, err := apiextensionsclientset.NewForConfig(config)
		if err != nil {
			return nil, err
		}

		dynmicClient, err := dynamic.NewForConfig(config)
		if err != nil {
			return nil, err
		}

		client := Client{
			KubernetesClient: kubernetesClient,
			CrdsClient:       crdClient,
			DynamicClient:    dynmicClient,
			HelmClient:       helmClient,
			KubeConfig:       cfg.KubeConfig,
		}
		values, err = describe(ctx, client, "", stream)
		if err != nil {
			return nil, err
		}

		return values, nil
	}
}
