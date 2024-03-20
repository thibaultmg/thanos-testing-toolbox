package thanos

import (
	"log/slog"
	"os"

	"github.com/observatorium/observatorium/configuration_go/abstr/kubernetes/thanos/query"
	kghelpers "github.com/observatorium/observatorium/configuration_go/kubegen/helpers"
	"github.com/observatorium/observatorium/configuration_go/kubegen/kubeyaml"
	"github.com/observatorium/observatorium/configuration_go/schemas/log"
)

func QueryManifests(dir string) {
	slog.Info("generating thanos-store manifests")
	// generate manifests
	kubeObjects := makeQuery().Objects()
	os.RemoveAll(dir)
	kubeyaml.WriteObjectsInDir(kubeObjects, dir)
}

func makeQuery() *query.QueryDeployment {
	opts := &query.QueryOptions{
		LogLevel:  log.LevelDebug,
		LogFormat: log.FormatLogfmt,
	}

	queryDepl := query.NewQuery(opts, "thanos-query", "latest")
	queryDepl.Replicas = 1
	queryDepl.ContainerResources = kghelpers.NewResourcesRequirements("100m", "", "200Mi", "400Mi")

	return queryDepl
}
