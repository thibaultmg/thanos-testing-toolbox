package redis

import (
	_ "embed"

	kghelpers "github.com/observatorium/observatorium/configuration_go/kubegen/helpers"
	"github.com/observatorium/observatorium/configuration_go/kubegen/workload"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

//go:embed resources/redis.conf
var redisConfig string

const (
	Port        = 6379
	ServiceName = "redis"
	Username    = "thanosuser"
	Password    = "thanospassword"
)

func Objects() []runtime.Object {
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
			ConfigMaps: map[string]map[string]string{
				"redis-config": {
					"redis.conf": redisConfig,
				},
			},
		},
	}

	container := deployment.ToContainer()
	container.Command = []string{"redis-server", "/etc/redis/redis.conf"}
	container.VolumeMounts = []corev1.VolumeMount{
		{
			Name:      "redis-config",
			MountPath: "/etc/redis",
		},
	}
	container.Volumes = []corev1.Volume{kghelpers.NewPodVolumeFromConfigMap("redis-config", "redis-config")}
	return deployment.Objects(container)
}
