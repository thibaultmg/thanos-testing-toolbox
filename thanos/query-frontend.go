package thanos

import (
	"log/slog"
	"os"

	"github.com/observatorium/observatorium/configuration_go/abstr/kubernetes/thanos/queryfrontend"
	kghelpers "github.com/observatorium/observatorium/configuration_go/kubegen/helpers"
	"github.com/observatorium/observatorium/configuration_go/kubegen/kubeyaml"
	"github.com/observatorium/observatorium/configuration_go/schemas/log"
)

func QueryFrontendManifests(dir string) {
	slog.Info("generating thanos-store manifests")
	// generate manifests
	kubeObjects := makeQuery().Objects()
	os.RemoveAll(dir)
	kubeyaml.WriteObjectsInDir(kubeObjects, dir)
}

func makeQueryFrontend() *queryfrontend.QueryFrontendDeployment {
	opts := &queryfrontend.QueryFrontendOptions{
		LogLevel:  log.LevelDebug,
		LogFormat: log.FormatLogfmt,
	}

	queryFrontendDepl := queryfrontend.NewQueryFrontend(opts, "thanos-query-frontend", "latest")
	queryFrontendDepl.Replicas = 1
	queryFrontendDepl.ContainerResources = kghelpers.NewResourcesRequirements("100m", "", "200Mi", "400Mi")

	return queryFrontendDepl
}
