package thanos

import (
	"github.com/observatorium/observatorium/configuration_go/abstr/kubernetes/thanos/store"
	kghelpers "github.com/observatorium/observatorium/configuration_go/kubegen/helpers"
	"github.com/observatorium/observatorium/configuration_go/schemas/log"
	"github.com/thibaultmg/thanos-testing-toolbox/internal/config"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func NewThanosStore(cfg config.DeploymentConfig, objstoreCfg string) []runtime.Object {
	opts := &store.StoreOptions{
		LogLevel:       log.LevelDebug,
		LogFormat:      log.FormatLogfmt,
		DataDir:        "/var/thanos/store",
		ObjstoreConfig: objstoreCfg,
	}

	if cfg.ImageTag == "" {
		cfg.ImageTag = "latest"
	}

	storeDepl := store.NewStore(opts, cfg.Namespace, cfg.ImageTag)
	storeDepl.VolumeSize = "1Gi"
	storeDepl.Replicas = 1
	storeDepl.ContainerResources = kghelpers.NewResourcesRequirements("100m", "", "200Mi", "400Mi")

	if cfg.Image != "" {
		storeDepl.Image = cfg.Image
	}

	// remove env var
	storeDepl.Env = storeDepl.Env[1:]

	ret := storeDepl.Objects()

	if cfg.Platform == config.PlatformKind {
		// persistent volume claim is changed to emptyDir volume
		// fmt.Println("changing persistent volume claim to emptyDir volume")
		sts := kghelpers.GetObject[*appsv1.StatefulSet](ret, "")
		oldVolume := sts.Spec.VolumeClaimTemplates[0]
		// fmt.Println(oldVolume)
		sts.Spec.VolumeClaimTemplates = nil
		sts.Spec.Template.Spec.Volumes = append(sts.Spec.Template.Spec.Volumes, corev1.Volume{
			Name: oldVolume.Name,
			VolumeSource: corev1.VolumeSource{
				EmptyDir: &corev1.EmptyDirVolumeSource{},
			},
		})
		// storeDepl.
	}

	// TODO: add route

	return ret
}
