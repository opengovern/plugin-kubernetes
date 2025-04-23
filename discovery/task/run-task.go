package task

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/opengovern/og-describer-kubernetes/discovery/pkg/orchestrator"
	authApi "github.com/opengovern/og-util/pkg/api"
	"github.com/opengovern/og-util/pkg/describe"
	"github.com/opengovern/og-util/pkg/httpclient"
	"github.com/opengovern/og-util/pkg/integration"
	"github.com/opengovern/og-util/pkg/jq"
	"github.com/opengovern/og-util/pkg/opengovernance-es-sdk"
	"github.com/opengovern/og-util/pkg/tasks"
	"github.com/opengovern/og-util/pkg/vault"
	coreApi "github.com/opengovern/opensecurity/services/core/api"
	coreClient "github.com/opengovern/opensecurity/services/core/client"
	"github.com/opengovern/opensecurity/services/tasks/scheduler"
	"go.uber.org/zap"
	"time"
)

type TaskRunner struct {
	vaultSrc            vault.VaultSourceConfig
	jq                  *jq.JobQueue
	coreServiceEndpoint string
	describeToken       string
	esClient            opengovernance.Client
	logger              *zap.Logger
	request             tasks.TaskRequest
	response            *scheduler.TaskResponse
}

func NewTaskRunner(ctx context.Context, jq *jq.JobQueue, coreServiceEndpoint string, describeToken string, esClient opengovernance.Client,
	logger *zap.Logger, request tasks.TaskRequest, response *scheduler.TaskResponse) (*TaskRunner, error) {

	vaultSc, err := vault.NewHashiCorpVaultClient(ctx, logger, request.VaultConfig.HashiCorp, request.VaultConfig.KeyId)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize HashiCorp vault: %w", err)
	}

	logger.Info("Vault setup complete")

	return &TaskRunner{
		vaultSrc:            vaultSc,
		jq:                  jq,
		coreServiceEndpoint: coreServiceEndpoint,
		describeToken:       describeToken,
		esClient:            esClient,
		logger:              logger,
		request:             request,
		response:            response,
	}, nil
}

type TaskResult struct {
}

type ResourceType struct {
	Name string
}

type Integration struct {
	IntegrationID   string
	ProviderID      string
	IntegrationType string
	Secret          string
	Labels          map[string]string
	Annotations     map[string]string
}

func (tr *TaskRunner) RunTask(ctx context.Context) error {

	taskResult := &TaskResult{}
	var err error
	var integrations []Integration

	inventoryClient := coreClient.NewCoreServiceClient(tr.coreServiceEndpoint)
	if _, ok := tr.request.TaskDefinition.Params["integrations_query"]; ok {
		integrations, err = retryWithBackoff("GetIntegrations", func() ([]Integration, error) {
			return GetIntegrationsFromQuery(inventoryClient, tr.request.TaskDefinition.Params)
		})
		if err != nil {
			tr.logger.Error("Error fetching integrations", zap.Error(err))
			return err
		}
	}

	for _, i := range integrations {
		err = tr.describeIntegrationResourceTypes(ctx, i)
		if err != nil {
			tr.logger.Error("Error describing integrations", zap.Error(err))
			return err
		}
	}

	jsonBytes, err := json.Marshal(taskResult)
	if err != nil {
		err = fmt.Errorf("failed Marshaling task result: %s", err.Error())
		return err
	}
	tr.response.Result = jsonBytes

	return nil
}

func (tr *TaskRunner) describeIntegrationResourceTypes(ctx context.Context, i Integration) error {
	var resourceTypes []ResourceType
	var err error

	config, err := tr.vaultSrc.Decrypt(ctx, i.Secret)
	if err != nil {
		return fmt.Errorf("decrypt error: %w", err)
	}

	inventoryClient := coreClient.NewCoreServiceClient(tr.coreServiceEndpoint)
	if _, ok := tr.request.TaskDefinition.Params["resource_types_query"]; ok {
		resourceTypes, err = retryWithBackoff("GetResourceTypes", func() ([]ResourceType, error) {
			return GetResourceTypesFromQuery(inventoryClient, tr.request.TaskDefinition.Params)
		})
		if err != nil {
			tr.logger.Error("Error fetching integrations", zap.Error(err))
			return err
		}
	}

	for _, rt := range resourceTypes {
		params := make(map[string]string)
		for key, value := range tr.request.TaskDefinition.Params {
			params[key] = fmt.Sprintf("%v", value)
		}
		for k, v := range params {
			ctx = context.WithValue(ctx, k, v)
		}

		job := describe.DescribeJob{
			JobID:                  tr.request.TaskDefinition.RunID,
			ResourceType:           rt.Name,
			IntegrationID:          i.IntegrationID,
			ProviderID:             i.ProviderID,
			DescribedAt:            time.Now().Unix(),
			IntegrationType:        integration.Type(i.IntegrationType),
			CipherText:             i.Secret,
			IntegrationLabels:      i.Labels,
			IntegrationAnnotations: i.Annotations,
		}
		_, err = orchestrator.Describe(ctx, tr.logger, job, params, config, tr.request.EsDeliverEndpoint,
			tr.request.IngestionPipelineEndpoint, tr.describeToken, tr.request.UseOpenSearch)
		if err != nil {
			tr.logger.Error("Error describing job", zap.Error(err))
			return err
		}
	}

	return nil
}

func GetIntegrationsFromQuery(coreServiceClient coreClient.CoreServiceClient, params map[string]any) ([]Integration, error) {
	if v, ok := params["integrations_query"]; ok {
		if vv, ok := v.(string); !ok {
			return nil, fmt.Errorf("query id should be a string")
		} else {
			queryResponse, err := coreServiceClient.RunQuery(&httpclient.Context{UserRole: authApi.AdminRole}, coreApi.RunQueryRequest{
				Query: &vv,
				Page: coreApi.Page{
					No:   1,
					Size: 1000,
				},
			})
			if err != nil {
				return nil, err
			}
			var integrations []Integration
			for _, r := range queryResponse.Result {
				integ := Integration{}
				for i, rc := range r {
					switch queryResponse.Headers[i] {
					case "integration_id":
						integ.IntegrationID = rc.(string)
					case "provider_id":
						integ.ProviderID = rc.(string)
					case "integration_type":
						integ.IntegrationType = rc.(string)
					case "secret":
						integ.Secret = rc.(string)
					case "annotations":
						if rc != nil {
							if jsonStr, ok := rc.(string); ok {
								var ann map[string]string
								if err := json.Unmarshal([]byte(jsonStr), &ann); err == nil {
									integ.Annotations = ann
								}
							}
						}
					case "labels":
						if rc != nil {
							if jsonStr, ok := rc.(string); ok {
								var lbl map[string]string
								if err := json.Unmarshal([]byte(jsonStr), &lbl); err == nil {
									integ.Labels = lbl
								}
							}
						}
					}
				}
				integrations = append(integrations, integ)
			}
			return integrations, nil
		}
	} else {
		return nil, fmt.Errorf("query id should be a string")
	}
}

func GetResourceTypesFromQuery(coreServiceClient coreClient.CoreServiceClient, params map[string]any) ([]ResourceType, error) {
	if v, ok := params["resource_types_query"]; ok {
		if vv, ok := v.(string); !ok {
			return nil, fmt.Errorf("query id should be a string")
		} else {
			queryResponse, err := coreServiceClient.RunQuery(&httpclient.Context{UserRole: authApi.AdminRole}, coreApi.RunQueryRequest{
				Query: &vv,
				Page: coreApi.Page{
					No:   1,
					Size: 1000,
				},
			})
			if err != nil {
				return nil, err
			}
			var resourceTypes []ResourceType
			for _, r := range queryResponse.Result {
				resourceType := ResourceType{}
				for i, rc := range r {
					if queryResponse.Headers[i] == "resource_type" {
						resourceType.Name = rc.(string)
					}
				}
				resourceTypes = append(resourceTypes, resourceType)
			}
			return resourceTypes, nil
		}
	} else {
		return nil, fmt.Errorf("query id should be a string")
	}
}
