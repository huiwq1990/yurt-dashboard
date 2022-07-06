package k8sproxy

import (
	"context"
	"fmt"
	"kubeapi/pkg/k8sproxy/proxy"
	"net/http"
	"strings"

	"github.com/pkg/errors"
	"k8s.io/klog/v2"

	"k8s.io/apiserver/pkg/authentication/user"
	"k8s.io/apiserver/pkg/authorization/authorizer"
	"k8s.io/apiserver/pkg/endpoints/request"
)

type Handler struct {
	authorizer     authorizer.Authorizer
	clustermanager ClusterManager
}

func NewApiServerHandler() *Handler {
	return &Handler{clustermanager: NewClusterManager()}
}

func GetClusterID(req *http.Request) (string, string) {
	parts := strings.Split(req.URL.Path, "/")
	if len(parts) > 4 && parts[0] == "" &&
		(parts[1] == "yurt") && parts[2] == "k8s" && parts[3] == "clusters" {
		return parts[1], parts[4]
	}
	return "", ""
}

const Module = "cluster proxy"

func (h *Handler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {

	//TODO 删除这块逻辑
	un := req.URL.Query().Get("user")
	if un != "" {
		klog.Errorf("should remove this backdoor")
		group := req.URL.Query().Get("group")
		ctx := request.WithUser(req.Context(), &user.DefaultInfo{
			Name:   un,
			Groups: []string{group},
		})
		req = req.WithContext(ctx)
	}

	//TODO 埋入默认user
	ctx := request.WithUser(req.Context(), &user.DefaultInfo{
		Name:   "admin",                    //tpaasUserInfo.UserName,
		Groups: []string{"system:masters"}, // groups,
	})
	req = req.WithContext(ctx)

	user, ok := request.UserFrom(req.Context())
	if !ok {
		klogAndHTTPError(rw, http.StatusUnauthorized, "%s could not found user context in request", Module)
		return
	}

	prefix, clusterID := GetClusterID(req)
	if clusterID == "" {
		klogAndHTTPError(rw, http.StatusUnauthorized, "%s could not resolve cluster info", Module)
		return
	}

	if !h.canAccess(req.Context(), user, clusterID) {
		klogAndHTTPError(rw, http.StatusUnauthorized, "%s not allow access", Module)
		return
	}

	prefix = fmt.Sprintf("/%s/k8s/clusters/%s", prefix, clusterID)

	handler, err := h.next(clusterID, prefix)
	if err != nil {
		klogAndHTTPError(rw, http.StatusUnauthorized, err.Error())
		return
	}

	handler.ServeHTTP(rw, req)
}

func (h *Handler) next(clusterID, prefix string) (http.Handler, error) {
	cd, err := h.clustermanager.KubeConfig(context.TODO(), clusterID)
	if err != nil {
		return nil, err
	}

	if cd == nil {
		return nil, errors.New(Module + "cluster not ready for service")
	}

	next := proxy.ImpersonatingHandler(prefix, cd)
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		req.Header.Set("X-API-URL-Prefix", prefix)
		next.ServeHTTP(rw, req)
	}), nil
}

func (h *Handler) canAccess(ctx context.Context, user user.Info, clusterID string) bool {
	return true
	//extra := map[string]authzv1.ExtraValue{}
	//for k, v := range user.GetExtra() {
	//	extra[k] = v
	//}
	//
	//resp, _, err := h.authorizer.Authorize(ctx, authorizer.AttributesRecord{
	//	ResourceRequest: true,
	//	User:            user,
	//	Verb:            "get",
	//	APIGroup:        managementv3.GroupName,
	//	APIVersion:      managementv3.Version,
	//	Resource:        "clusters",
	//	Name:            clusterID,
	//})
	//
	//return err == nil && resp == authorizer.DecisionAllow
}

// klogAndHTTPError logs the error message and write back to client
func klogAndHTTPError(w http.ResponseWriter, errCode int, format string, i ...interface{}) {
	errMsg := fmt.Sprintf(format, i...)
	klog.Error(errMsg)
	http.Error(w, errMsg, errCode)
}
