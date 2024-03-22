package thanos

import (
	"github.com/observatorium/observatorium/configuration_go/abstr/kubernetes/thanos/receive"
	kghelpers "github.com/observatorium/observatorium/configuration_go/kubegen/helpers"
	"github.com/observatorium/observatorium/configuration_go/schemas/log"
	"github.com/thibaultmg/thanos-testing-toolbox/internal/config"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

type ReceiveOutput struct {
	Objects  []runtime.Object
	GrpcPort int
	SvcName  string
}

func NewThanosReceive(cfg config.DeploymentConfig, objstoreCfg string) *ReceiveOutput {
	opts := receive.NewDefaultIngestorOptions()
	opts.LogLevel = log.LevelDebug
	opts.ObjstoreConfig = objstoreCfg

	receiveDepl := receive.NewIngestor(opts, cfg.Namespace, cfg.ImageTag)
	receiveDepl.Replicas = 1
	receiveDepl.ContainerResources = kghelpers.NewResourcesRequirements("100m", "", "100Mi", "200Mi")
	receiveDepl.Name = "thanos-receive"

	if cfg.ImageTag == "" {
		receiveDepl.ImageTag = "latest"
	}

	if cfg.Image != "" {
		receiveDepl.Image = cfg.Image
	}

	// remove env var
	receiveDepl.Env = receiveDepl.Env[1:]

	objects := receiveDepl.Objects()

	if cfg.Platform == config.PlatformKind {
		// persistent volume claim is changed to emptyDir volume
		sts := kghelpers.GetObject[*appsv1.StatefulSet](objects, "")
		oldVolume := sts.Spec.VolumeClaimTemplates[0]
		sts.Spec.VolumeClaimTemplates = nil
		sts.Spec.Template.Spec.Volumes = append(sts.Spec.Template.Spec.Volumes, corev1.Volume{
			Name: oldVolume.Name,
			VolumeSource: corev1.VolumeSource{
				EmptyDir: &corev1.EmptyDirVolumeSource{},
			},
		})
	}

	return &ReceiveOutput{
		Objects:  objects,
		SvcName:  receiveDepl.Name,
		GrpcPort: opts.GrpcAddress.Port,
	}
}
