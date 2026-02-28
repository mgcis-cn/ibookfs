// Package handler provides HTTP handlers for account secret management.
package handler

import (
	"context"

	"github.com/mgcis-cn/ibookfs/internal/apiserver/biz/accountsecret"
	"github.com/mgcis-cn/ibookfs/internal/apiserver/model"

	v1 "github.com/mgcis-cn/ibookfs/pkg/api/apiserver/v1"
	contextx "github.com/mgcis-cn/ibookfs/pkg/context"
)

type AccountSecretRouter interface {
	ListAccountSecrets(ctx context.Context, req *v1.ListAccountSecretsRequest) (*v1.ListAccountSecretsResponse, error)
	GetAccountSecret(ctx context.Context, req *v1.GetAccountSecretRequest) (*v1.GetAccountSecretResponse, error)
	CreateAccountSecret(ctx context.Context, req *v1.CreateAccountSecretRequest) (*v1.CreateAccountSecretResponse, error)
	DeleteAccountSecret(ctx context.Context, req *v1.DeleteAccountSecretRequest) (*v1.DeleteAccountSecretResponse, error)
	DisableAccountSecret(ctx context.Context, req *v1.DisableAccountSecretRequest) (*v1.DisableAccountSecretResponse, error)
	EnableAccountSecret(ctx context.Context, req *v1.EnableAccountSecretRequest) (*v1.EnableAccountSecretResponse, error)
}

// ListAccountSecrets handles account secret list requests.
func (h *handler) ListAccountSecrets(ctx context.Context, req *v1.ListAccountSecretsRequest) (*v1.ListAccountSecretsResponse, error) {
	userID := contextx.UserId(ctx)
	if req.Page < 1 {
		req.Page = 1
	}
	if req.PageSize < 1 || req.PageSize > 100 {
		req.PageSize = 20
	}

	secrets, total, err := h.biz.AccountSecret().List(ctx, uint(userID), req.Page, req.PageSize)
	if err != nil {
		return &v1.ListAccountSecretsResponse{}, err
	}

	// Remove secret key from response (it should never be returned in list)
	cleanedSecrets := make([]map[string]interface{}, len(secrets))
	for i, s := range secrets {
		cleanedSecrets[i] = map[string]interface{}{
			"id":            s.ID,
			"name":          s.Name,
			"account_key":   s.AccountKey,
			"scope":         s.Scope,
			"resource_id":   s.ResourceID,
			"resource_type": s.ResourceType,
			"permissions":   s.Permissions,
			"is_enabled":    s.IsEnabled,
			"expires_at":    s.ExpiresAt,
			"created_at":    s.CreatedAt,
			"last_used_at":  s.LastUsedAt,
		}
	}

	data := map[string]interface{}{
		"items":     cleanedSecrets,
		"total":     total,
		"page":      req.Page,
		"page_size": req.PageSize,
	}

	return &v1.ListAccountSecretsResponse{
		Data:     data,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
		Message:  "success",
	}, nil
}

// GetAccountSecret handles account secret detail requests.
func (h *handler) GetAccountSecret(ctx context.Context, req *v1.GetAccountSecretRequest) (*v1.GetAccountSecretResponse, error) {
	userID := contextx.UserId(ctx)
	secret, err := h.biz.AccountSecret().GetByID(ctx, uint(req.Id), uint(userID))
	if err != nil {
		return &v1.GetAccountSecretResponse{}, err
	}

	// Don't return secret key in detail view
	data := map[string]interface{}{
		"id":            secret.ID,
		"name":          secret.Name,
		"account_key":   secret.AccountKey,
		"scope":         secret.Scope,
		"resource_id":   secret.ResourceID,
		"resource_type": secret.ResourceType,
		"permissions":   secret.Permissions,
		"is_enabled":    secret.IsEnabled,
		"expires_at":    secret.ExpiresAt,
		"created_at":    secret.CreatedAt,
		"last_used_at":  secret.LastUsedAt,
	}

	return &v1.GetAccountSecretResponse{
		Data:    data,
		Message: "success",
	}, nil
}

// CreateAccountSecret handles account secret creation requests.
func (h *handler) CreateAccountSecret(ctx context.Context, req *v1.CreateAccountSecretRequest) (*v1.CreateAccountSecretResponse, error) {
	userID := contextx.UserId(ctx)

	var expiresIn *int64
	if req.ExpiresInDays != nil {
		seconds := int64(*req.ExpiresInDays * 86400)
		expiresIn = &seconds
	}

	createReq := &accountsecret.CreateRequest{
		Name:         req.Name,
		Scope:        model.AccountSecretScope(req.Scope),
		ResourceID:   req.ResourceID,
		ResourceType: req.ResourceType,
		Permissions:  req.Permissions,
		ExpiresIn:    expiresIn,
	}

	result, err := h.biz.AccountSecret().Create(ctx, uint(userID), createReq)
	if err != nil {
		return &v1.CreateAccountSecretResponse{}, err
	}

	return &v1.CreateAccountSecretResponse{
		Data:    result,
		Message: "Account secret created. Please save the secret key securely, it will not be shown again.",
	}, nil
}

// DeleteAccountSecret handles account secret deletion requests.
func (h *handler) DeleteAccountSecret(ctx context.Context, req *v1.DeleteAccountSecretRequest) (*v1.DeleteAccountSecretResponse, error) {
	if err := h.biz.AccountSecret().Delete(ctx, uint(req.Id), uint(contextx.UserId(ctx))); err != nil {
		return &v1.DeleteAccountSecretResponse{}, err
	}

	return &v1.DeleteAccountSecretResponse{
		Message: "deleted",
	}, nil
}

// DisableAccountSecret disables an account secret.
func (h *handler) DisableAccountSecret(ctx context.Context, req *v1.DisableAccountSecretRequest) (*v1.DisableAccountSecretResponse, error) {
	if err := h.biz.AccountSecret().Disable(ctx, uint(req.Id), uint(contextx.UserId(ctx))); err != nil {
		return &v1.DisableAccountSecretResponse{}, err
	}

	return &v1.DisableAccountSecretResponse{
		Message: "disabled",
	}, nil
}

// EnableAccountSecret enables an account secret.
func (h *handler) EnableAccountSecret(ctx context.Context, req *v1.EnableAccountSecretRequest) (*v1.EnableAccountSecretResponse, error) {
	if err := h.biz.AccountSecret().Enable(ctx, uint(req.Id), uint(contextx.UserId(ctx))); err != nil {
		return &v1.EnableAccountSecretResponse{}, err
	}

	return &v1.EnableAccountSecretResponse{
		Message: "enabled",
	}, nil
}
