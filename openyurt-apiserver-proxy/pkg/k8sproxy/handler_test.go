package k8sproxy

import (
	"k8s.io/apiserver/pkg/authentication/user"
	"k8s.io/apiserver/pkg/endpoints/request"
	"k8s.io/klog/v2"
	"net/http"
	"testing"
)

func TestUser(t *testing.T) {
	req := &http.Request{}

	ctx := request.WithUser(req.Context(), &user.DefaultInfo{Name: "abc"})
	req = req.WithContext(ctx)

	info, ok := request.UserFrom(req.Context())
	if !ok {
		t.Fatal("user not found")
	}
	klog.Infof(info.GetName())
}
