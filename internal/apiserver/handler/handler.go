// Package handler contains HTTP request handlers.
package handler

import (
	"context"

	"github.com/mgcis-cn/ibookfs/internal/apiserver/biz"
	"github.com/mgcis-cn/ibookfs/internal/apiserver/model"
	v1 "github.com/mgcis-cn/ibookfs/pkg/api/apiserver/v1"
)

type handler struct {
	oauthConfig map[model.OAuthProvider]model.OAuthConfig
	biz         biz.Biz
}

var _ v1.ApiServerHTTPServer = (*handler)(nil)

func New(
	oauthConfig map[model.OAuthProvider]model.OAuthConfig,
	biz biz.Biz,
) *handler {
	return &handler{oauthConfig: oauthConfig, biz: biz}
}

func (h *handler) Health(ctx context.Context, req *v1.HealthRequest) (*v1.HealthResponse, error) {
	return &v1.HealthResponse{
		Status: "ok",
	}, nil
}
