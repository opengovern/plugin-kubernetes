package main

import (
	"encoding/json"
	"github.com/jackc/pgtype"
	"github.com/opengovern/og-describer-kubernetes/global"
	constants2 "github.com/opengovern/og-describer-kubernetes/global/constants"
	"github.com/opengovern/og-describer-kubernetes/global/maps"
	"github.com/opengovern/og-describer-kubernetes/platform/constants"
	"github.com/opengovern/og-util/pkg/integration"
	"github.com/opengovern/og-util/pkg/integration/interfaces"
	"k8s.io/client-go/tools/clientcmd"
)

type Integration struct{}

func (i *Integration) GetConfiguration() (interfaces.IntegrationConfiguration, error) {
	return interfaces.IntegrationConfiguration{
		NatsScheduledJobsTopic:   global.JobQueueTopic,
		NatsManualJobsTopic:      global.JobQueueTopicManuals,
		NatsStreamName:           global.StreamName,
		NatsConsumerGroup:        global.ConsumerGroup,
		NatsConsumerGroupManuals: global.ConsumerGroupManuals,

		SteampipePluginName: "kubernetes",

		UISpec:   constants.UISpec,
		Manifest: constants.Manifest,
		SetupMD:  constants.SetupMd,

		DescriberDeploymentName: constants.DescriberDeploymentName,
		DescriberRunCommand:     constants.DescriberRunCommand,
	}, nil
}

func (i *Integration) HealthCheck(jsonData []byte, providerId string, labels map[string]string, annotations map[string]string) (bool, error) {
	var credentials constants2.IntegrationCredentials
	err := json.Unmarshal(jsonData, &credentials)
	if err != nil {
		return false, err
	}
	isHealthy, err := IntegrationHealthcheck(credentials, Config{})

	return isHealthy, err
}

func (i *Integration) DiscoverIntegrations(jsonData []byte) ([]integration.Integration, error) {
	var credentials constants2.IntegrationCredentials
	err := json.Unmarshal(jsonData, &credentials)
	if err != nil {
		return nil, err
	}

	config, err := clientcmd.RESTConfigFromKubeConfig([]byte(credentials.KubeConfig))
	if err != nil {
		return nil, err
	}

	info, err := DoDiscovery(credentials.KubeConfig)
	if err != nil {
		return nil, err
	}

	var labels pgtype.JSONB
	if err := labels.Set(info); err != nil {
		return nil, err
	}

	return []integration.Integration{
		{
			ProviderID: config.Host,
			Name:       config.Host,
			Labels:     labels,
		},
	}, nil
}

func (i *Integration) GetResourceTypesByLabels(labels map[string]string) ([]interfaces.ResourceTypeConfiguration, error) {
	var resourceTypesMap []interfaces.ResourceTypeConfiguration
	for _, resourceType := range maps.ResourceTypesList {
		var resource interfaces.ResourceTypeConfiguration
		if v, ok := maps.ResourceTypeConfigs[resourceType]; ok {
			resource.Description = v.Description
			resource.Params = v.Params
			resource.Name = v.Name
			resource.IntegrationType = v.IntegrationType
			resource.Table = maps.ResourceTypesToTables[v.Name]
			resourceTypesMap = append(resourceTypesMap, resource)

		}
	}
	return resourceTypesMap, nil
}
func (i *Integration) GetResourceTypeFromTableName(tableName string) (string, error) {
	if v, ok := maps.TablesToResourceTypes[tableName]; ok {
		return v, nil
	}

	return "", nil
}

func (i *Integration) GetIntegrationType() (integration.Type, error) {
	return constants.IntegrationName, nil
}

func (i *Integration) ListAllTables() (map[string][]interfaces.CloudQLColumn, error) {
	plugin := global.Plugin()
	tables := make(map[string][]interfaces.CloudQLColumn)
	for tableKey, table := range plugin.TableMap {
		columns := make([]interfaces.CloudQLColumn, 0, len(table.Columns))
		for _, column := range table.Columns {
			columns = append(columns, interfaces.CloudQLColumn{Name: column.Name, Type: column.Type.String()})
		}
		tables[tableKey] = columns
	}

	return tables, nil
}

func (i *Integration) Ping() error {
	return nil
}
