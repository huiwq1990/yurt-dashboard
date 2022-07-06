package app

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"kubeapi/pkg/k8sproxy"
	"net/http"
)

var Router *gin.Engine

func prometheusHandler() gin.HandlerFunc {
	h := promhttp.Handler()

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

// InitRouter init global web service http router
func InitRouter() error {
	Router = gin.Default()
	Router.Use(Cors())

	// global api
	// 监控路由
	Router.GET("/metrics", prometheusHandler())

	// apiserver的http代理
	k8sproxy := k8sproxy.NewApiServerHandler()
	clusterRouter(k8sproxy)

	return nil
}

func clusterRouter(handler http.Handler) {
	Router.Any("/yurt/k8s/clusters/*target", gin.WrapH(handler))
}
