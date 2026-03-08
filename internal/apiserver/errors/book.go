package errors

import (
	pkgerr "github.com/mgcis-cn/ibookfs/pkg/errors"
)

// Book module error codes: 20000-20999
const (
	CodeBookNotFound     = 20001
	CodeBookCreateFailed = 20002
	CodeBookUpdateFailed = 20003
	CodeBookDeleteFailed = 20004
	CodeBookListFailed   = 20005
)

// Book module sentinel errors.
var (
	ErrBookNotFound = pkgerr.NotFound(CodeBookNotFound, "图书不存在")
)

// BookCreateFailed wraps a book creation failure.
func BookCreateFailed(cause error) pkgerr.Error {
	return pkgerr.WithCause(
		pkgerr.InternalServer(CodeBookCreateFailed, "创建图书失败"),
		cause,
	)
}

// BookUpdateFailed wraps a book update failure.
func BookUpdateFailed(cause error) pkgerr.Error {
	return pkgerr.WithCause(
		pkgerr.InternalServer(CodeBookUpdateFailed, "更新图书失败"),
		cause,
	)
}

// BookDeleteFailed wraps a book deletion failure.
func BookDeleteFailed(cause error) pkgerr.Error {
	return pkgerr.WithCause(
		pkgerr.InternalServer(CodeBookDeleteFailed, "删除图书失败"),
		cause,
	)
}

// BookListFailed wraps a book list failure.
func BookListFailed(cause error) pkgerr.Error {
	return pkgerr.WithCause(
		pkgerr.InternalServer(CodeBookListFailed, "获取图书列表失败"),
		cause,
	)
}
