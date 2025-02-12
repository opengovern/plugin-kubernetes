package models

type IntegrationCredentials struct {
	KubeConfig []byte `json:"kubeconfig"`
}
