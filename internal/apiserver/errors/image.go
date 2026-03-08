package errors

import (
	pkgerr "github.com/mgcis-cn/ibookfs/pkg/errors"
)

// Image module error codes: 30000-30999
const (
	CodeImageNotFound       = 30001
	CodeImageAccessDenied   = 30002
	CodeImageFileRequired   = 30003
	CodeImageUploadFailed   = 30004
	CodeImageDeleteFailed   = 30005
	CodeImageDownloadFailed = 30006
	CodeImageListFailed     = 30007
	CodeGroupCreateFailed   = 30008
	CodeGroupAddFailed      = 30009
	CodeImageInvalidType    = 30010
	CodeImageInvalidExt     = 30011
	CodeImageSizeTooLarge   = 30012
	CodeImageMIMENotAllowed = 30013
	CodeImageDecodeFailed   = 30014
	CodeImageProcessFailed  = 30015
	CodeImageWorkerStopped  = 30016
	CodeImageQueueFull      = 30017
)

// Image module sentinel errors.
var (
	ErrImageNotFound      = pkgerr.NotFound(CodeImageNotFound, "图片不存在")
	ErrImageAccessDenied  = pkgerr.Forbidden(CodeImageAccessDenied, "无权访问该图片")
	ErrImageFileRequired  = pkgerr.BadRequest(CodeImageFileRequired, "请上传文件")
	ErrImageInvalidExt    = pkgerr.BadRequest(CodeImageInvalidExt, "不支持的文件扩展名")
	ErrImageWorkerStopped = pkgerr.InternalServer(CodeImageWorkerStopped, "图片处理服务未运行")
	ErrImageQueueFull     = pkgerr.InternalServer(CodeImageQueueFull, "图片处理队列已满")
)

// ImageSizeTooLarge creates an error for file size exceeding the limit.
func ImageSizeTooLarge(size, maxSize int64) pkgerr.Error {
	return pkgerr.Newf(CodeImageSizeTooLarge, 400, "文件大小(%d)超过最大限制(%d)", size, maxSize)
}

// ImageMIMENotAllowed creates an error for disallowed MIME type.
func ImageMIMENotAllowed(mimeType string) pkgerr.Error {
	return pkgerr.BadRequest(CodeImageMIMENotAllowed, "不支持的文件类型: "+mimeType)
}

// ImageUploadFailed wraps an image upload failure.
func ImageUploadFailed(cause error) pkgerr.Error {
	return pkgerr.WithCause(
		pkgerr.InternalServer(CodeImageUploadFailed, "图片上传失败"),
		cause,
	)
}

// ImageDeleteFailed wraps an image deletion failure.
func ImageDeleteFailed(cause error) pkgerr.Error {
	return pkgerr.WithCause(
		pkgerr.InternalServer(CodeImageDeleteFailed, "图片删除失败"),
		cause,
	)
}

// ImageDownloadFailed wraps an image download failure.
func ImageDownloadFailed(cause error) pkgerr.Error {
	return pkgerr.WithCause(
		pkgerr.InternalServer(CodeImageDownloadFailed, "图片下载失败"),
		cause,
	)
}

// ImageListFailed wraps an image list failure.
func ImageListFailed(cause error) pkgerr.Error {
	return pkgerr.WithCause(
		pkgerr.InternalServer(CodeImageListFailed, "获取图片列表失败"),
		cause,
	)
}

// GroupCreateFailed wraps a group creation failure.
func GroupCreateFailed(cause error) pkgerr.Error {
	return pkgerr.WithCause(
		pkgerr.InternalServer(CodeGroupCreateFailed, "创建图片组失败"),
		cause,
	)
}

// GroupAddFailed wraps an add-to-group failure.
func GroupAddFailed(cause error) pkgerr.Error {
	return pkgerr.WithCause(
		pkgerr.InternalServer(CodeGroupAddFailed, "添加图片到组失败"),
		cause,
	)
}
