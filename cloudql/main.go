package main

import (
	"github.com/opengovern/og-describer-kubernetes/cloudql/kubernetes"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{PluginFunc: kubernetes.Plugin})
}
