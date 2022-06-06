module github.com/KubeOperator/kubepi

go 1.16

require (
	github.com/asdine/storm/v3 v3.2.1
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/coreos/etcd v3.3.13+incompatible
	github.com/creack/pty v1.1.15 // indirect
	github.com/docker/distribution v2.7.1+incompatible
	github.com/docker/docker v20.10.9+incompatible // indirect
	github.com/fatih/color v1.13.0 // indirect
	github.com/fsnotify/fsnotify v1.5.1 // indirect
	github.com/go-ldap/ldap/v3 v3.4.1
	github.com/go-logr/logr v0.4.0
	github.com/gofrs/flock v0.8.1
	github.com/golang/snappy v0.0.4 // indirect
	github.com/google/go-cmp v0.5.6 // indirect
	github.com/google/uuid v1.2.0
	github.com/imdario/mergo v0.3.12 // indirect
	github.com/iris-contrib/swagger/v12 v12.0.1
	github.com/kataras/iris/v12 v12.2.0-alpha2.0.20210427211137-fa175eb84754
	github.com/klauspost/compress v1.13.5 // indirect
	github.com/onsi/gomega v1.15.0 // indirect
	github.com/pkg/errors v0.9.1
	github.com/prometheus/common v0.29.0 // indirect
	github.com/sirupsen/logrus v1.8.1
	github.com/skip2/go-qrcode v0.0.0-20200617195104-da1b6568686e
	github.com/spf13/cobra v1.2.1
	github.com/spf13/viper v1.8.1
	github.com/swaggo/swag v1.8.2
	github.com/xeipuuv/gojsonpointer v0.0.0-20190905194746-02993c407bfb // indirect
	github.com/xlzd/gotp v0.0.0-20220110052318-fab697c03c2c
	golang.org/x/crypto v0.0.0-20210921155107-089bfa567519
	golang.org/x/oauth2 v0.0.0-20211005180243-6b3c2da341f1 // indirect
	golang.org/x/text v0.3.7
	google.golang.org/protobuf v1.27.1 // indirect
	gopkg.in/igm/sockjs-go.v2 v2.1.0
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b
	helm.sh/helm/v3 v3.7.0
	k8s.io/api v0.22.2
	k8s.io/apiextensions-apiserver v0.22.2
	k8s.io/apimachinery v0.22.2
	k8s.io/cli-runtime v0.22.1
	k8s.io/client-go v0.22.2
	k8s.io/klog/v2 v2.9.0
	k8s.io/kubectl v0.22.1
)

replace github.com/KubeOperator/webkubectl/gotty v0.0.0-20210927072155-e9ce79172471 => ./thirdparty/gotty
