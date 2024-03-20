package redis

import (
	"log/slog"
	"os"

	kghelpers "github.com/observatorium/observatorium/configuration_go/kubegen/helpers"
	"github.com/observatorium/observatorium/configuration_go/kubegen/kubeyaml"
	"github.com/observatorium/observatorium/configuration_go/kubegen/workload"
	"k8s.io/apimachinery/pkg/runtime"
)

func Manifests(dir string) {
	slog.Info("generating thanos-store manifests")
	// generate manifests
	os.RemoveAll(dir)
	kubeyaml.WriteObjectsInDir(makeRedis(), dir)
}

func makeRedis() []runtime.Object {
	deployment := workload.DeploymentWorkload{
		Replicas: 1,
		PodConfig: workload.PodConfig{
			Name:               "redis",
			CommonLabels:       map[string]string{"app": "redis"},
			Image:              "redis",
			ImageTag:           "latest",
			ContainerResources: kghelpers.NewResourcesRequirements("100m", "", "200Mi", "400Mi"),
			LivenessProbe:      kghelpers.NewProbe("", 6379, kghelpers.ProbeConfig{InitialDelaySeconds: 15, PeriodSeconds: 20}),
			ReadinessProbe:     kghelpers.NewProbe("", 6379, kghelpers.ProbeConfig{InitialDelaySeconds: 5, PeriodSeconds: 10}),
		},
	}

	container := deployment.ToContainer()
	return deployment.Objects(container)
}
