package thanos

import (
	"fmt"
	"net"
	"net/netip"

	"github.com/observatorium/observatorium/configuration_go/abstr/kubernetes/thanos/query"
	kghelpers "github.com/observatorium/observatorium/configuration_go/kubegen/helpers"
	"github.com/observatorium/observatorium/configuration_go/schemas/log"
	"github.com/thibaultmg/thanos-testing-toolbox/internal/config"
	"k8s.io/apimachinery/pkg/runtime"
)

type StoreEndpoint struct {
	GrpcPort    int
	ServiceName string
}

type QueryConfig struct {
	StoreEndpoints []StoreEndpoint
}

type QueryOutput struct {
	Objects     []runtime.Object
	ServiceName string
	HttpPort    int
}

func NewThanosQuery(cfg config.DeploymentConfig, queryCfg *QueryConfig) *QueryOutput {
	serviceName := "thanos-query"
	httpPort := 10902
	opts := &query.QueryOptions{
		LogLevel:    log.LevelDebug,
		LogFormat:   log.FormatLogfmt,
		HttpAddress: net.TCPAddrFromAddrPort(netip.MustParseAddrPort(fmt.Sprintf("0.0.0.0:%d", httpPort))),
	}

	for _, storeEndpoint := range queryCfg.StoreEndpoints {
		opts.Endpoint = append(opts.Endpoint, fmt.Sprintf("dnssrv+_grpc._tcp.%s.%s.svc.cluster.local", storeEndpoint.ServiceName, cfg.Namespace))
	}

	if cfg.ImageTag == "" {
		cfg.ImageTag = "latest"
	}

	queryDepl := query.NewQuery(opts, cfg.Namespace, cfg.ImageTag)
	queryDepl.Replicas = 1
	queryDepl.ContainerResources = kghelpers.NewResourcesRequirements("100m", "", "100Mi", "200Mi")
	queryDepl.Name = serviceName

	if cfg.Image != "" {
		queryDepl.Image = cfg.Image
	}

	return &QueryOutput{
		Objects:     queryDepl.Objects(),
		ServiceName: serviceName,
		HttpPort:    httpPort,
	}
}
