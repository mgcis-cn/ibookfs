package errors

import (
	pkgerr "github.com/mgcis-cn/ibookfs/pkg/errors"
)

// AccountSecret module error codes: 40000-40999
const (
	CodeAccountAccessDenied      = 40001
	CodeAccountInvalidCredentials = 40002
	CodeAccountResourceRequired  = 40003
	CodeAccountPermissionRequired = 40004
	CodeAccountInvalidPermission = 40005
	CodeAccountCreateFailed      = 40006
	CodeAccountKeyGenFailed      = 40007
	CodeAccountSecretGenFailed   = 40008
	CodeAccountHashFailed        = 40009
)

// AccountSecret module sentinel errors.
var (
	ErrAccountAccessDenied       = pkgerr.Forbidden(CodeAccountAccessDenied, "无权操作")
	ErrAccountInvalidCredentials = pkgerr.Unauthorized(CodeAccountInvalidCredentials, "凭证无效")
	ErrAccountResourceRequired   = pkgerr.BadRequest(CodeAccountResourceRequired, "资源级密钥需要提供resource_id和resource_type")
	ErrAccountPermissionRequired = pkgerr.BadRequest(CodeAccountPermissionRequired, "至少需要一个权限")
)

// AccountInvalidPermission creates an error for invalid permission.
func AccountInvalidPermission(perm string) pkgerr.Error {
	return pkgerr.BadRequest(CodeAccountInvalidPermission, "无效的权限: "+perm)
}
