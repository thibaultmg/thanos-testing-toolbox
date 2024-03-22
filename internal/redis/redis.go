package redis

import (
	_ "embed"

	kghelpers "github.com/observatorium/observatorium/configuration_go/kubegen/helpers"
	"github.com/observatorium/observatorium/configuration_go/kubegen/workload"
	"github.com/thibaultmg/thanos-testing-toolbox/internal/config"
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

func Objects(cfg config.DeploymentConfig) []runtime.Object {
	deployment := workload.DeploymentWorkload{
		Replicas: 1,
		PodConfig: workload.PodConfig{
			Namespace:          cfg.Namespace,
			Name:               "redis",
			CommonLabels:       map[string]string{"app": "redis"},
			Image:              cfg.Image,
			ImageTag:           cfg.ImageTag,
			ContainerResources: kghelpers.NewResourcesRequirements("100m", "", "100Mi", "200Mi"),
			LivenessProbe: &corev1.Probe{
				InitialDelaySeconds: 15,
				TimeoutSeconds:      5,
				ProbeHandler: corev1.ProbeHandler{
					Exec: &corev1.ExecAction{
						Command: []string{"sh", "-c", "redis-cli ping"},
					},
				},
			},
			// LivenessProbe:      kghelpers.NewProbe("", Port, kghelpers.ProbeConfig{InitialDelaySeconds: 15, PeriodSeconds: 20}),
			// ReadinessProbe:     kghelpers.NewProbe("", Port, kghelpers.ProbeConfig{InitialDelaySeconds: 5, PeriodSeconds: 10}),
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
	container.Ports = []corev1.ContainerPort{
		{
			Name:          "redis",
			ContainerPort: Port,
		},
	}
	container.ServicePorts = []corev1.ServicePort{
		kghelpers.NewServicePort("http", Port, Port),
	}
	return deployment.Objects(container)
}
