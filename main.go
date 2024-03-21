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
	thanosImageTag := "v0.34.1"
	depCfg := config.DeploymentConfig{
		Namespace: namespace,
		Platform:  platform,
	}
	objs := minio.StoreWithData(depCfg)
	WriteObjectsInDir(objs, "manifests/minio")

	// Set up Thanos Store
	depCfg.ImageTag = thanosImageTag
	storeOutput := thanos.NewThanosStore(depCfg, minio.ObjectStoreConfig)
	WriteObjectsInDir(storeOutput.Objects, "manifests/thanos-store")

	// Set up Thanos Query
	objs = thanos.NewThanosQuery(depCfg, &thanos.QueryConfig{
		StoreServiceName: storeOutput.SvcName,
		StoreGRPCPort:    storeOutput.GrpcPort,
	})
	WriteObjectsInDir(objs, "manifests/thanos-query")

	// thanos.QueryManifests("manifests/thanos-query")
	// thanos.QueryFrontendManifests("manifests/thanos-query-frontend")
	// WriteObjectsInDir(redis.Objects(), "manifests/redis")
}

func WriteObjectsInDir(kubeObjects []runtime.Object, dir string) {
	os.RemoveAll(dir)
	kubeyaml.WriteObjectsInDir(kubeObjects, dir)
}
