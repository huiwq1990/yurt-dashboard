module kubeapi

go 1.16

require (
	github.com/gin-gonic/gin v1.6.3
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.11.0
	github.com/sirupsen/logrus v1.8.1
	github.com/spf13/cobra v1.2.1
	github.com/spf13/pflag v1.0.5
	k8s.io/apimachinery v0.18.8
	k8s.io/apiserver v0.18.8
	k8s.io/client-go v0.18.8
	k8s.io/klog/v2 v2.0.0
	sigs.k8s.io/controller-runtime v0.6.2
)
