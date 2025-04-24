package pkg

import (
	esSinkClient "github.com/opengovern/og-util/pkg/es/ingest/client"
	"github.com/opengovern/og-util/pkg/jq"
	"github.com/opengovern/og-util/pkg/opengovernance-es-sdk"
	"go.uber.org/zap"
	"os"
)

type Worker struct {
	logger   *zap.Logger
	esClient opengovernance.Client
	jq       *jq.JobQueue

	esSinkClient esSinkClient.EsSinkServiceClient
}

var (
	ManualTriggers = os.Getenv("MANUAL_TRIGGERS")
)
