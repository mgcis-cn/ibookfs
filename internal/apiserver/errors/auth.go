package errors

import (
	"net/http"

	pkgerr "github.com/mgcis-cn/ibookfs/pkg/errors"
)

// Auth module error codes: 10000-10999
const (
	CodeRateLimited       = 10001
	CodeCodeInvalid       = 10002
	CodeCodeExpired       = 10003
	CodeCodeWrong         = 10004
	CodeUserNotFound      = 10005
	CodeUserDisabled      = 10006
	CodeEmailExists       = 10007
	CodeEmailInUse        = 10008
	CodeInvalidToken      = 10009
	CodeSessionExpired    = 10010
	CodeCannotUnlinkLast  = 10011
	CodeInvalidCodeType   = 10012
	CodeProviderRequired  = 10013
	CodeStateGenFailed    = 10014
	CodeInvalidUserID     = 10015
)

// Auth module sentinel errors.
var (
	ErrRateLimited      = pkgerr.TooManyRequests(CodeRateLimited, "请等待60秒后重新发送")
	ErrCodeInvalid      = pkgerr.BadRequest(CodeCodeInvalid, "验证码无效或已过期")
	ErrCodeWrong        = pkgerr.BadRequest(CodeCodeWrong, "验证码错误")
	ErrUserNotFound     = pkgerr.NotFound(CodeUserNotFound, "用户不存在，请先注册")
	ErrUserDisabled     = pkgerr.Forbidden(CodeUserDisabled, "账号已被禁用")
	ErrEmailExists      = pkgerr.Conflict(CodeEmailExists, "该邮箱已被注册")
	ErrEmailInUse       = pkgerr.Conflict(CodeEmailInUse, "该邮箱已被其他账号使用")
	ErrInvalidToken     = pkgerr.Unauthorized(CodeInvalidToken, "无效的令牌")
	ErrSessionExpired   = pkgerr.Unauthorized(CodeSessionExpired, "会话不存在或已过期")
	ErrCannotUnlinkLast = pkgerr.BadRequest(CodeCannotUnlinkLast, "无法解绑最后一个登录方式")
	ErrInvalidCodeType  = pkgerr.BadRequest(CodeInvalidCodeType, "无效的验证码类型")
	ErrProviderRequired = pkgerr.BadRequest(CodeProviderRequired, "OAuth提供商不能为空")
	ErrStateGenFailed   = pkgerr.New(CodeStateGenFailed, http.StatusInternalServerError, "生成安全状态失败")
	ErrInvalidUserID    = pkgerr.Unauthorized(CodeInvalidUserID, "无效的用户ID")
)
