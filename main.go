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

	// Set up Thanos Receive
	receiveOutput := thanos.NewThanosReceive(depCfg, minio.ObjectStoreConfig)
	WriteObjectsInDir(receiveOutput.Objects, "manifests/thanos-receive")

	// Set up Thanos Query
	queryOutput := thanos.NewThanosQuery(depCfg, &thanos.QueryConfig{
		StoreEndpoints: []thanos.StoreEndpoint{
			{
				GrpcPort:    receiveOutput.GrpcPort,
				ServiceName: receiveOutput.SvcName,
			},
			{
				GrpcPort:    storeOutput.GrpcPort,
				ServiceName: storeOutput.SvcName,
			},
		},
	})
	WriteObjectsInDir(queryOutput.Objects, "manifests/thanos-query")

	// Set up Thanos Query-frontend
	objs = thanos.NewThanosQueryFrontend(depCfg, &thanos.QueryFrontendConfig{
		QueryServiceName:       queryOutput.ServiceName,
		QueryPort:              queryOutput.HttpPort,
		WithRedisResponseCache: true,
	})
	WriteObjectsInDir(objs, "manifests/thanos-query-frontend")
}

func WriteObjectsInDir(kubeObjects []runtime.Object, dir string) {
	os.RemoveAll(dir)
	kubeyaml.WriteObjectsInDir(kubeObjects, dir)
}
