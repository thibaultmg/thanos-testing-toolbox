package thanos

import (
	"fmt"
	"net"
	"net/netip"

	"github.com/observatorium/observatorium/configuration_go/abstr/kubernetes/thanos/store"
	kghelpers "github.com/observatorium/observatorium/configuration_go/kubegen/helpers"
	"github.com/observatorium/observatorium/configuration_go/schemas/log"
	"github.com/thibaultmg/thanos-testing-toolbox/internal/config"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

type StoreOutput struct {
	Objects  []runtime.Object
	GrpcPort int
	HttpPort int
	SvcName  string
}

func NewThanosStore(cfg config.DeploymentConfig, objstoreCfg string) StoreOutput {
	grpcPort := 10901
	httpPort := 10902
	name := "thanos-store"
	opts := &store.StoreOptions{
		LogLevel:       log.LevelDebug,
		LogFormat:      log.FormatLogfmt,
		DataDir:        "/var/thanos/store",
		ObjstoreConfig: objstoreCfg,
		GrpcAddress:    net.TCPAddrFromAddrPort(netip.MustParseAddrPort(fmt.Sprintf("0.0.0.0:%d", grpcPort))),
		HttpAddress:    net.TCPAddrFromAddrPort(netip.MustParseAddrPort(fmt.Sprintf("0.0.0.0:%d", httpPort))),
	}

	if cfg.ImageTag == "" {
		cfg.ImageTag = "latest"
	}

	storeDepl := store.NewStore(opts, cfg.Namespace, cfg.ImageTag)
	storeDepl.VolumeSize = "1Gi"
	storeDepl.Replicas = 1
	storeDepl.ContainerResources = kghelpers.NewResourcesRequirements("100m", "", "100Mi", "200Mi")
	storeDepl.Name = name

	if cfg.Image != "" {
		storeDepl.Image = cfg.Image
	}

	// remove env var
	storeDepl.Env = storeDepl.Env[1:]

	objects := storeDepl.Objects()

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

	// TODO: add route

	return StoreOutput{
		Objects:  objects,
		GrpcPort: grpcPort,
		HttpPort: httpPort,
		SvcName:  name,
	}
}
