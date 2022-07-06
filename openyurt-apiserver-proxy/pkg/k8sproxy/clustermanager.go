package k8sproxy

import (
	"k8s.io/client-go/rest"
	ctrl "sigs.k8s.io/controller-runtime"

	"context"
)

type ClusterManager interface {
	KubeConfig(ctx context.Context, clusterName string) (*rest.Config, error)
}

var kubeConfig = ctrl.GetConfigOrDie()

func NewClusterManager() ClusterManager {

	return &defaultClusterManager{}
}

type defaultClusterManager struct {
}

func (cm *defaultClusterManager) KubeConfig(ctx context.Context, clusterName string) (*rest.Config, error) {
	return kubeConfig, nil
}
