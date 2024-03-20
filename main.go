package main

import (
	"os"

	"github.com/observatorium/observatorium/configuration_go/kubegen/kubeyaml"
	"github.com/thibaultmg/thanos-testing-toolbox/internal/config"
	"github.com/thibaultmg/thanos-testing-toolbox/internal/thanos"
	"github.com/thibaultmg/thanos-testing-toolbox/pkg/minio"
	"k8s.io/apimachinery/pkg/runtime"
)

func main() {
	namespace := "thanos"
	platform := config.PlatformKind
	objs := minio.StoreWithData(config.DeploymentConfig{
		Namespace: namespace,
		Platform:  platform,
	})
	WriteObjectsInDir(objs, "manifests/minio")

	// Set up Thanos Store
	objs = thanos.NewThanosStore(config.DeploymentConfig{
		Namespace: namespace,
		Platform:  platform,
		ImageTag:  "v0.34.1",
	}, minio.ObjectStoreConfig)
	WriteObjectsInDir(objs, "manifests/thanos-store")

	// thanos.QueryManifests("manifests/thanos-query")
	// thanos.QueryFrontendManifests("manifests/thanos-query-frontend")
	// WriteObjectsInDir(redis.Objects(), "manifests/redis")
}

func WriteObjectsInDir(kubeObjects []runtime.Object, dir string) {
	os.RemoveAll(dir)
	kubeyaml.WriteObjectsInDir(kubeObjects, dir)
}
