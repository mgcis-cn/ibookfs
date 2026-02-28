// Package handler contains HTTP request handlers.
package handler

import (
	"context"

	"github.com/mgcis-cn/ibookfs/internal/apiserver/biz"
	v1 "github.com/mgcis-cn/ibookfs/pkg/api/apiserver/v1"
	oauthPkg "github.com/mgcis-cn/ibookfs/pkg/authn/oauth"
)

type handler struct {
	oauthFactory *oauthPkg.Factory
	biz          biz.Biz
}

var _ v1.ApiServerHTTPServer = (*handler)(nil)

func New(
	oauthFactory *oauthPkg.Factory,
	biz biz.Biz,
) *handler {
	return &handler{oauthFactory: oauthFactory, biz: biz}
}

func (h *handler) Health(ctx context.Context, req *v1.HealthRequest) (*v1.HealthResponse, error) {
	return &v1.HealthResponse{
		Status: "ok",
	}, nil
}
