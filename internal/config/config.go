package config

type Platform string

const (
	PlatformKubernetes Platform = "kubernetes"
	PlatformOpenShift  Platform = "openshift"
	PlatformKind       Platform = "kind"
)

type DeploymentConfig struct {
	Namespace string
	ImageTag  string
	Image     string
	Platform  Platform
}
