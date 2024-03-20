module github.com/thibaultmg/thanos-testing-toolbox

go 1.21.5

replace github.com/observatorium/observatorium => ../observatorium

require (
	github.com/observatorium/observatorium v0.0.0-00010101000000-000000000000
	github.com/openshift/api v3.9.0+incompatible
	gopkg.in/yaml.v3 v3.0.1
	k8s.io/api v0.29.1
	k8s.io/apimachinery v0.29.1
)

require (
	github.com/alecthomas/template v0.0.0-20190718012654-fb15b899a751 // indirect
	github.com/alecthomas/units v0.0.0-20211218093645-b94a6e3cc137 // indirect
	github.com/bwplotka/mimic v0.2.1-0.20230303101552-f705cca2f4a4 // indirect
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/ghodss/yaml v1.0.0 // indirect
	github.com/go-kit/log v0.2.1 // indirect
	github.com/go-logfmt/logfmt v0.6.0 // indirect
	github.com/go-logr/logr v1.3.0 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/google/gofuzz v1.2.0 // indirect
	github.com/grafana/regexp v0.0.0-20221122212121-6b5c0a4cb7fd // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring v0.70.0 // indirect
	github.com/prometheus/common v0.45.0 // indirect
	github.com/prometheus/prometheus v0.48.1 // indirect
	github.com/rodaine/hclencoder v0.0.1 // indirect
	golang.org/x/exp v0.0.0-20231006140011-7918f672742d // indirect
	golang.org/x/net v0.19.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	gopkg.in/alecthomas/kingpin.v2 v2.2.6 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	k8s.io/klog/v2 v2.110.1 // indirect
	k8s.io/utils v0.0.0-20231127182322-b307cd553661 // indirect
	sigs.k8s.io/json v0.0.0-20221116044647-bc3834ca7abd // indirect
	sigs.k8s.io/structured-merge-diff/v4 v4.4.1 // indirect
	sigs.k8s.io/yaml v1.4.0 // indirect
)
