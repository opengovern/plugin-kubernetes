package global

import "github.com/opengovern/og-util/pkg/integration"

const (
	IntegrationTypeLower = "kubernetes"                                    // example: aws, azure
	IntegrationName      = integration.Type("kubernetes")                  // example: aws_account, github_account
	OGPluginRepoURL      = "github.com/opengovern/og-describer-kubernetes" // example: github.com/opengovern/og-describer-aws
)

type IntegrationCredentials struct {
	// TODO
}
