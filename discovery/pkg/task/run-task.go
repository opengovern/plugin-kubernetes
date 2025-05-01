package task

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/opengovern/og-describer-kubernetes/discovery/envs"
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
	AllIntegrations             []string                      `json:"all_integrations"`
	AllIntegrationsCount        int                           `json:"all_integrations_count"`
	ProgressedIntegrations      map[string]*IntegrationResult `json:"progressed_integrations"`
	ProgressedIntegrationsCount int                           `json:"progressed_integrations_count"`
}

type IntegrationResult struct {
	IntegrationID              string               `json:"integration_id"`
	AllResourceTypes           []string             `json:"all_resource_types"`
	AllResourceTypesCount      int                  `json:"all_resource_types_count"`
	ResourceTypeResults        []ResourceTypeResult `json:"resource_type_results"`
	FinishedResourceTypesCount int                  `json:"finished_resource_types_count"`
}

type ResourceTypeResult struct {
	ResourceType  string `json:"resource_type"`
	Error         string `json:"error"`
	ResourceCount int    `json:"resource_count"`
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
	tr.logger.Info("Run task")

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
		taskResult.AllIntegrations = append(taskResult.AllIntegrations, i.IntegrationID)
	}
	taskResult.AllIntegrationsCount = len(integrations)
	taskResult.ProgressedIntegrations = make(map[string]*IntegrationResult)

	tr.logger.Info("Describing integrations", zap.Any("integrations", integrations))

	for _, i := range integrations {
		err = tr.describeIntegrationResourceTypes(ctx, i, taskResult)
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

func (tr *TaskRunner) describeIntegrationResourceTypes(ctx context.Context, i Integration, taskResult *TaskResult) error {
	taskResult.ProgressedIntegrations[i.IntegrationID] = &IntegrationResult{
		IntegrationID: i.IntegrationID,
	}
	taskResult.ProgressedIntegrationsCount++

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

	tr.logger.Info("Describing integration", zap.String("integration_id", i.IntegrationID), zap.Any("resource_types", resourceTypes))

	for _, rt := range resourceTypes {
		taskResult.ProgressedIntegrations[i.IntegrationID].AllResourceTypes = append(taskResult.ProgressedIntegrations[i.IntegrationID].AllResourceTypes, rt.Name)
	}
	taskResult.ProgressedIntegrations[i.IntegrationID].AllResourceTypesCount = len(resourceTypes)

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
		resources, err := orchestrator.Describe(ctx, tr.logger, job, params, config, tr.request.EsDeliverEndpoint,
			tr.request.IngestionPipelineEndpoint, tr.describeToken, tr.request.UseOpenSearch)
		errMsg := ""
		if err != nil {
			tr.logger.Error("Error describing job", zap.Error(err))
			errMsg = err.Error()
		}
		taskResult.ProgressedIntegrations[i.IntegrationID].ResourceTypeResults = append(
			taskResult.ProgressedIntegrations[i.IntegrationID].ResourceTypeResults,
			ResourceTypeResult{
				ResourceType:  rt.Name,
				Error:         errMsg,
				ResourceCount: len(resources),
			})

		taskResult.ProgressedIntegrations[i.IntegrationID].FinishedResourceTypesCount = len(taskResult.ProgressedIntegrations[i.IntegrationID].ResourceTypeResults)
		jsonBytes, err := json.Marshal(taskResult)
		if err != nil {
			err = fmt.Errorf("failed Marshaling task result: %s", err.Error())
			return err
		}
		tr.response.Result = jsonBytes
		responseJson, marshalErr := json.Marshal(tr.response)
		if marshalErr != nil {
			tr.logger.Error("failed to create final job result json", zap.Error(marshalErr))
			return marshalErr
		}
		msgId := fmt.Sprintf("task-run-result-%d", tr.request.TaskDefinition.RunID)
		if _, err = tr.jq.Produce(ctx, envs.ResultTopicName, responseJson, msgId); err != nil { // Use original ctx
			tr.logger.Error("failed to publish initial InProgress job status", zap.String("response", string(responseJson)), zap.Error(err))
			return err
		}
		tr.logger.Info("describing resource type finished", zap.String("integration_id", i.IntegrationID), zap.String("resource_type", rt.Name))
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
							if obj, ok := rc.(map[string]interface{}); ok {
								for k, v := range obj {
									if vStr, ok := v.(string); ok {
										integ.Annotations[k] = vStr
									}
								}
							}
						}
					case "labels":
						if rc != nil {
							if obj, ok := rc.(map[string]interface{}); ok {
								for k, v := range obj {
									if vStr, ok := v.(string); ok {
										integ.Labels[k] = vStr
									}
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
