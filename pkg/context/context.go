package context

import (
	"context"

	"github.com/go-kratos/kratos/v2/transport/http"
)

type userIdCtx struct{}

func WithUserId(ctx context.Context, userId int) context.Context {
	return context.WithValue(ctx, userIdCtx{}, userId)
}

func UserId(ctx context.Context) int {
	userId, _ := ctx.Value(userIdCtx{}).(int)
	return userId
}

type userEmailCtx struct{}

func WithUserEmail(ctx context.Context, email string) context.Context {
	return context.WithValue(ctx, userEmailCtx{}, email)
}

func UserEmail(ctx context.Context) (string, bool) {
	v := ctx.Value(userEmailCtx{})
	if v == nil {
		return "", false
	}
	return v.(string), true
}

type userRoleCtx struct{}

func WithUserRole(ctx context.Context, role string) context.Context {
	return context.WithValue(ctx, userRoleCtx{}, role)
}

func UserRole(ctx context.Context) (string, bool) {
	v := ctx.Value(userRoleCtx{})
	if v == nil {
		return "", false
	}
	return v.(string), true
}

type authTypeCtx struct{}

func WithAuthType(ctx context.Context, authType string) context.Context {
	return context.WithValue(ctx, authTypeCtx{}, authType)
}

func AuthType(ctx context.Context) (string, bool) {
	v := ctx.Value(authTypeCtx{})
	if v == nil {
		return "", false
	}
	return v.(string), true
}

type signatureAuthCtx struct{}

type SignatureContext struct {
	AccountSecret interface{}
	AccountKey    string
}

func WithSignatureAuth(ctx context.Context, sigCtx *SignatureContext) context.Context {
	return context.WithValue(ctx, signatureAuthCtx{}, sigCtx)
}

func SignatureAuth(ctx context.Context) (*SignatureContext, bool) {
	v := ctx.Value(signatureAuthCtx{})
	if v == nil {
		return nil, false
	}
	return v.(*SignatureContext), true
}

type ownerIDCtx struct{}

func WithOwnerID(ctx context.Context, ownerID uint) context.Context {
	return context.WithValue(ctx, ownerIDCtx{}, ownerID)
}

func OwnerID(ctx context.Context) (uint, bool) {
	v := ctx.Value(ownerIDCtx{})
	if v == nil {
		return 0, false
	}
	return v.(uint), true
}

type requestCtx struct{}

func WithRequest(ctx context.Context, req *http.Request) context.Context {
	return context.WithValue(ctx, requestCtx{}, req)
}

func Request(ctx context.Context) *http.Request {
	req, _ := ctx.Value(requestCtx{}).(*http.Request)
	return req
}

type responseCtx struct{}

func WithResponse(ctx context.Context, resp http.ResponseWriter) context.Context {
	return context.WithValue(ctx, responseCtx{}, resp)
}

func Response(ctx context.Context) http.ResponseWriter {
	resp, _ := ctx.Value(responseCtx{}).(http.ResponseWriter)
	return resp
}
