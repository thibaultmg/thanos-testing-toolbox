package thanos

import (
	"fmt"

	"github.com/observatorium/observatorium/configuration_go/abstr/kubernetes/thanos/query"
	kghelpers "github.com/observatorium/observatorium/configuration_go/kubegen/helpers"
	"github.com/observatorium/observatorium/configuration_go/schemas/log"
	"github.com/thibaultmg/thanos-testing-toolbox/internal/config"
	"k8s.io/apimachinery/pkg/runtime"
)

type QueryConfig struct {
	StoreServiceName string
	StoreGRPCPort    int
}

func NewThanosQuery(cfg config.DeploymentConfig, queryCfg *QueryConfig) []runtime.Object {
	opts := &query.QueryOptions{
		LogLevel:  log.LevelDebug,
		LogFormat: log.FormatLogfmt,
		Endpoint: []string{
			fmt.Sprintf("dnssrv+_grpc._tcp.%s.%s.svc.cluster.local", queryCfg.StoreServiceName, cfg.Namespace),
		},
	}

	if cfg.ImageTag == "" {
		cfg.ImageTag = "latest"
	}

	queryDepl := query.NewQuery(opts, cfg.Namespace, cfg.ImageTag)
	queryDepl.Replicas = 1
	queryDepl.ContainerResources = kghelpers.NewResourcesRequirements("100m", "", "200Mi", "400Mi")

	if cfg.Image != "" {
		queryDepl.Image = cfg.Image
	}

	return queryDepl.Objects()
}
