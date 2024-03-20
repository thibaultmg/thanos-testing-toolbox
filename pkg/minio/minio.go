package minio

import (
	_ "embed"
	"fmt"
	"log/slog"
	"time"

	kghelpers "github.com/observatorium/observatorium/configuration_go/kubegen/helpers"
	"github.com/observatorium/observatorium/configuration_go/kubegen/workload"
	routev1 "github.com/openshift/api/route/v1"
	"github.com/thibaultmg/thanos-testing-toolbox/internal/config"
	"gopkg.in/yaml.v3"
	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
)

//go:embed resources/thanosbench.yaml
var thanosbenchConfig string

const (
	RootUser     = "minioadmin"
	RootPassword = "minioadmin"
	APIPort      = 9000
	Bucket       = "thanos"
)

var ObjectStoreConfig = fmt.Sprintf(`
type: S3 
config:
  bucket: %s
  endpoint: minio:%d
  insecure: true
  access_key: %s
  secret_key: %s
`, Bucket, APIPort, RootUser, RootPassword)

// StoreWithData is a store with data
func StoreWithData(cfg config.DeploymentConfig) []runtime.Object {
	slog.Info("generating minio object store manifests")

	if cfg.ImageTag == "" {
		cfg.ImageTag = "latest"
	}

	if cfg.Image == "" {
		cfg.Image = "quay.io/minio/minio"
	}

	minioDeploy := workload.StatefulSetWorkload{
		Replicas:   1,
		VolumeSize: "1Gi",
		PodConfig: workload.PodConfig{
			Namespace:          cfg.Namespace,
			Name:               "minio",
			CommonLabels:       map[string]string{"app": "minio"},
			Image:              cfg.Image,
			ImageTag:           cfg.ImageTag,
			ContainerResources: kghelpers.NewResourcesRequirements("100m", "", "200Mi", "400Mi"),
			LivenessProbe:      kghelpers.NewProbe("/minio/health/live", APIPort, kghelpers.ProbeConfig{InitialDelaySeconds: 10, PeriodSeconds: 20}),
			ReadinessProbe:     kghelpers.NewProbe("/minio/health/ready", APIPort, kghelpers.ProbeConfig{InitialDelaySeconds: 10, PeriodSeconds: 20}),
			Env: []corev1.EnvVar{
				{
					Name:  "MINIO_ROOT_USER",
					Value: RootUser,
				},
				{
					Name:  "MINIO_ROOT_PASSWORD",
					Value: RootPassword,
				},
			},
		},
	}

	minioContainer := minioDeploy.ToContainer()
	minioContainer.Args = []string{
		"server",
		"/data",
		"--console-address",
		":9090",
	}

	minioContainer.Ports = []corev1.ContainerPort{
		{
			Name:          "api",
			ContainerPort: APIPort,
			Protocol:      corev1.ProtocolTCP,
		},
		{
			Name:          "console",
			ContainerPort: 9090,
			Protocol:      corev1.ProtocolTCP,
		},
	}

	minioContainer.ServicePorts = []corev1.ServicePort{
		{
			Name:       "api",
			Port:       APIPort,
			TargetPort: intstr.FromInt(APIPort),
		},
		{
			Name:       "console",
			Port:       9090,
			TargetPort: intstr.FromInt(9090),
		},
	}

	minioContainer.VolumeMounts = []corev1.VolumeMount{
		{
			Name:      "data",
			MountPath: "/data",
		},
	}

	minioContainer.Volumes = []corev1.Volume{
		{
			Name: "data",
			VolumeSource: corev1.VolumeSource{
				EmptyDir: &corev1.EmptyDirVolumeSource{},
			},
		},
	}

	// make init container the creates the thanos bucket
	bucketCreator := &workload.Container{
		Name:     "bucket-creator",
		Image:    "alpine",
		ImageTag: "latest",
		Command: []string{
			"sh",
			"-c",
			fmt.Sprintf("mkdir -p /data/%s", Bucket),
		},
		VolumeMounts: []corev1.VolumeMount{
			{
				Name:      "data",
				MountPath: "/data",
			},
		},
	}

	minioDeploy.InitContainers = append(minioDeploy.InitContainers, bucketCreator)

	objects := minioDeploy.Objects(minioContainer)
	depl := kghelpers.GetObject[*appsv1.StatefulSet](objects, "minio")
	depl.ObjectMeta.Annotations = map[string]string{
		"openshift.io/required-scc": "privileged",
	}

	if cfg.Platform == config.PlatformOpenShift {
		objects = append(objects, makeOpenshiftResources(cfg)...)
	}

	objects = append(objects, newDataGenerator(cfg)...)

	return objects
}

// newDataGenerator generates data using a thanosbench init container
// and exports it to the minio object store using the minio client
func newDataGenerator(cfg config.DeploymentConfig) []runtime.Object {
	// Create a thanosbench init container
	thanosBench := corev1.Container{
		Name:  "thanosbench",
		Image: "quay.io/thanos/thanosbench:v0.3.0-rc.0",
		Args: []string{
			"block",
			"gen",
			"--output.dir=/data",
			"--config-file=/etc/thanosbench/thanosbench.yaml",
		},
		VolumeMounts: []corev1.VolumeMount{
			{
				Name:      "data",
				MountPath: "/data",
			},
			{
				Name:      "config",
				MountPath: "/etc/thanosbench",
			},
		},
	}

	// Create the job
	job := batchv1.Job{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Job",
			APIVersion: batchv1.SchemeGroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "thanosbench",
			Namespace: cfg.Namespace,
		},
		Spec: batchv1.JobSpec{
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					RestartPolicy:  corev1.RestartPolicyNever,
					InitContainers: []corev1.Container{thanosBench},
					Containers: []corev1.Container{
						{
							Name:  "mc",
							Image: cfg.Image + ":" + cfg.ImageTag,
							Command: []string{
								// add alias to minio client
								"sh",
								"-c",
								fmt.Sprintf("sleep 15 && mc alias set local http://minio.thanos.svc.cluster.local:%d %s %s && mc cp --recursive /data/ local/%s/", APIPort, RootUser, RootPassword, Bucket),
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "data",
									MountPath: "/data",
								},
							},
						},
					},

					Volumes: []corev1.Volume{
						{
							Name: "data",
							VolumeSource: corev1.VolumeSource{
								EmptyDir: &corev1.EmptyDirVolumeSource{},
							},
						},
						{
							Name: "config",
							VolumeSource: corev1.VolumeSource{
								ConfigMap: &corev1.ConfigMapVolumeSource{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: "thanosbench",
									},
								},
							},
						},
					},
				},
			},
		},
	}

	parsedCfg := []interface{}{}
	err := yaml.Unmarshal([]byte(thanosbenchConfig), &parsedCfg)
	if err != nil {
		panic(err)
	}

	// set timestamps
	for _, v := range parsedCfg {
		for _, serie := range v.(map[string]interface{})["series"].([]interface{}) {
			serie.(map[string]interface{})["mintime"] = time.Now().Add(-30 * 24 * time.Hour).UTC().UnixMilli()
			serie.(map[string]interface{})["maxtime"] = time.Now().UTC().UnixMilli()
		}
	}

	thanosbenchConfigR, err := yaml.Marshal(parsedCfg)
	if err != nil {
		panic(err)
	}

	// add thanosbench configmap
	configMap := corev1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ConfigMap",
			APIVersion: corev1.SchemeGroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "thanosbench",
			Namespace: cfg.Namespace,
		},
		Data: map[string]string{
			"thanosbench.yaml": string(thanosbenchConfigR),
		},
	}

	return []runtime.Object{&job, &configMap}
}

func makeOpenshiftResources(cfg config.DeploymentConfig) []runtime.Object {
	ret := []runtime.Object{}

	// Add rbac providing scc with hostaccess to minio
	rolebinding := rbacv1.RoleBinding{
		TypeMeta: metav1.TypeMeta{
			Kind:       "RoleBinding",
			APIVersion: "rbac.authorization.k8s.io/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "minio-scc",
			Namespace: cfg.Namespace,
		},
		Subjects: []rbacv1.Subject{
			{
				Kind:      "ServiceAccount",
				Name:      "minio",
				Namespace: cfg.Namespace,
			},
		},
		RoleRef: rbacv1.RoleRef{
			Kind:     "ClusterRole",
			Name:     "system:openshift:scc:privileged",
			APIGroup: "rbac.authorization.k8s.io",
		},
	}
	ret = append(ret, &rolebinding)

	// add route to minio
	route := &routev1.Route{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Route",
			APIVersion: routev1.SchemeGroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "minio",
			Namespace: cfg.Namespace,
			// Labels:    maps.Clone(kghelpers.GetObject[*appsv1.Deployment](manifests, "").ObjectMeta.Labels),
			// Annotations: map[string]string{
			// 	"cert-manager.io/issuer-kind": "ClusterIssuer",
			// 	"cert-manager.io/issuer-name": "letsencrypt-prod-http",
			// },
		},
		Spec: routev1.RouteSpec{
			Port: &routev1.RoutePort{
				TargetPort: intstr.FromString("console"),
			},
			// TLS: &routev1.TLSConfig{
			// 	Termination:                   routev1.TLSTerminationReencrypt,
			// 	InsecureEdgeTerminationPolicy: routev1.InsecureEdgeTerminationPolicyRedirect,
			// },
			To: routev1.RouteTargetReference{
				Kind: "Service",
				Name: "minio",
			},
		},
	}

	ret = append(ret, route)

	return ret
}
