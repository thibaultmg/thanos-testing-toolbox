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

type QueryConfig struct {
	StoreServiceName string
	StoreGRPCPort    int
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
		LogLevel:  log.LevelDebug,
		LogFormat: log.FormatLogfmt,
		Endpoint: []string{
			fmt.Sprintf("dnssrv+_grpc._tcp.%s.%s.svc.cluster.local", queryCfg.StoreServiceName, cfg.Namespace),
		},
		HttpAddress: net.TCPAddrFromAddrPort(netip.MustParseAddrPort(fmt.Sprintf("127.0.0.1:%d", httpPort))),
	}

	if cfg.ImageTag == "" {
		cfg.ImageTag = "latest"
	}

	queryDepl := query.NewQuery(opts, cfg.Namespace, cfg.ImageTag)
	queryDepl.Replicas = 1
	queryDepl.ContainerResources = kghelpers.NewResourcesRequirements("100m", "", "200Mi", "400Mi")
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
