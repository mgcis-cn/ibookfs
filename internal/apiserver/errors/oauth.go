package errors

import (
	pkgerr "github.com/mgcis-cn/ibookfs/pkg/errors"
)

// OAuth module error codes: 11000-11999
const (
	CodeOAuthUnsupported    = 11001
	CodeOAuthStateProcessed = 11002
	CodeOAuthExchangeFailed = 11003
	CodeOAuthUserInfoFailed = 11004
	CodeOAuthMissingUserID  = 11005
	CodeOAuthMissingEmail   = 11006
	CodeOAuthCreateFailed   = 11007
	CodeOAuthTokenFailed    = 11008
	CodeOAuthAlreadyLinked  = 11009
	CodeOAuthLinkedOther    = 11010
)

// OAuth module sentinel errors.
var (
	ErrOAuthStateProcessed = pkgerr.BadRequest(CodeOAuthStateProcessed, "该OAuth授权已被处理")
	ErrOAuthMissingUserID  = pkgerr.BadRequest(CodeOAuthMissingUserID, "OAuth用户信息缺失: 缺少用户ID")
	ErrOAuthMissingEmail   = pkgerr.BadRequest(CodeOAuthMissingEmail, "OAuth用户信息缺失: 请公开邮箱或授权邮箱权限")
	ErrOAuthAlreadyLinked  = pkgerr.Conflict(CodeOAuthAlreadyLinked, "该OAuth账号已绑定")
	ErrOAuthLinkedOther    = pkgerr.Conflict(CodeOAuthLinkedOther, "该OAuth账号已绑定到其他用户")
)

// OAuthUnsupported creates an error for unsupported OAuth provider.
func OAuthUnsupported(provider string) pkgerr.Error {
	return pkgerr.BadRequest(CodeOAuthUnsupported, "不支持的OAuth提供商: "+provider)
}

// OAuthExchangeFailed wraps a code exchange failure.
func OAuthExchangeFailed(cause error) pkgerr.Error {
	return pkgerr.WithCause(
		pkgerr.InternalServer(CodeOAuthExchangeFailed, "OAuth授权码交换失败"),
		cause,
	)
}

// OAuthUserInfoFailed wraps a user info retrieval failure.
func OAuthUserInfoFailed(cause error) pkgerr.Error {
	return pkgerr.WithCause(
		pkgerr.InternalServer(CodeOAuthUserInfoFailed, "获取OAuth用户信息失败"),
		cause,
	)
}

// OAuthCreateFailed wraps a user creation failure.
func OAuthCreateFailed(cause error) pkgerr.Error {
	return pkgerr.WithCause(
		pkgerr.InternalServer(CodeOAuthCreateFailed, "创建用户失败"),
		cause,
	)
}

// OAuthTokenFailed wraps a token generation failure.
func OAuthTokenFailed(cause error) pkgerr.Error {
	return pkgerr.WithCause(
		pkgerr.InternalServer(CodeOAuthTokenFailed, "生成令牌失败"),
		cause,
	)
}
