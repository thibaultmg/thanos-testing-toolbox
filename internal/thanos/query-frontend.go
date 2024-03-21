package thanos

import (
	"fmt"

	"github.com/observatorium/observatorium/configuration_go/abstr/kubernetes/thanos/queryfrontend"
	kghelpers "github.com/observatorium/observatorium/configuration_go/kubegen/helpers"
	"github.com/observatorium/observatorium/configuration_go/schemas/log"
	"github.com/observatorium/observatorium/configuration_go/schemas/thanos/cache"
	rediscfg "github.com/observatorium/observatorium/configuration_go/schemas/thanos/cache/redis"
	"github.com/thibaultmg/thanos-testing-toolbox/internal/config"
	"github.com/thibaultmg/thanos-testing-toolbox/internal/redis"
	"k8s.io/apimachinery/pkg/runtime"
)

type QueryFrontendConfig struct {
	QueryServiceName       string
	QueryPort              int
	WithRedisResponseCache bool
}

func NewThanosQueryFrontend(cfg config.DeploymentConfig, qfCfg *QueryFrontendConfig) []runtime.Object {
	ret := []runtime.Object{}
	opts := &queryfrontend.QueryFrontendOptions{
		LogLevel:                   log.LevelDebug,
		LogFormat:                  log.FormatLogfmt,
		QueryFrontendDownstreamURL: fmt.Sprintf("http://%s.%s.svc.cluster.local:%d", qfCfg.QueryServiceName, cfg.Namespace, qfCfg.QueryPort),
	}

	if qfCfg.WithRedisResponseCache {
		opts.QueryRangeResponseCacheConfig = cache.NewResponseCacheConfig(rediscfg.RedisClientConfig{
			Addr: fmt.Sprintf("%s:%d", redis.ServiceName, redis.Port),
		})
		ret = append(ret, redis.Objects()...)
	}

	queryFrontendDepl := queryfrontend.NewQueryFrontend(opts, cfg.Namespace, cfg.ImageTag)
	queryFrontendDepl.Replicas = 1
	queryFrontendDepl.ContainerResources = kghelpers.NewResourcesRequirements("100m", "", "200Mi", "400Mi")

	if cfg.ImageTag == "" {
		queryFrontendDepl.ImageTag = "latest"
	}

	if cfg.Image != "" {
		queryFrontendDepl.Image = cfg.Image
	}

	ret = append(ret, queryFrontendDepl.Objects()...)
	return ret
}
