package errors

import (
	pkgerr "github.com/mgcis-cn/ibookfs/pkg/errors"
)

// Email module error codes: 50000-50999
const (
	CodeEmailDisabled = 50001
)

// Email module sentinel errors.
var (
	ErrEmailDisabled = pkgerr.InternalServer(CodeEmailDisabled, "邮件服务未启用")
)
