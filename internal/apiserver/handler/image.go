// Package handler provides HTTP handlers for image management.
package handler

import (
	"context"
	"errors"
	"io"

	"github.com/mgcis-cn/ibookfs/internal/apiserver/biz/image"
	"github.com/mgcis-cn/ibookfs/internal/apiserver/model"
	v1 "github.com/mgcis-cn/ibookfs/pkg/api/apiserver/v1"
	contextx "github.com/mgcis-cn/ibookfs/pkg/context"
)

type ImageRouter interface {
	ListImages(ctx context.Context, req *v1.ListImagesRequest) (*v1.ListImagesResponse, error)
	GetImage(ctx context.Context, req *v1.GetImageRequest) (*v1.GetImageResponse, error)
	UploadImage(ctx context.Context, req *v1.UploadImageRequest) (*v1.UploadImageResponse, error)
	DeleteImage(ctx context.Context, req *v1.DeleteImageRequest) (*v1.DeleteImageResponse, error)
	CreateGroup(ctx context.Context, req *v1.CreateGroupRequest) (*v1.CreateGroupResponse, error)
	AddImagesToGroup(ctx context.Context, req *v1.AddImagesToGroupRequest) (*v1.AddImagesToGroupResponse, error)
	ServeImage(ctx context.Context, req *v1.ServeImageRequest) (*v1.ServeImageResponse, error)
}

// UploadImage handles image upload requests.
func (h *handler) UploadImage(ctx context.Context, req *v1.UploadImageRequest) (*v1.UploadImageResponse, error) {
	// Get file from form
	if req.File == nil {
		return &v1.UploadImageResponse{}, errors.New("file is required")
	}

	request := contextx.Request(ctx)
	file, handler, err := request.FormFile("file")
	if err != nil {
		return &v1.UploadImageResponse{}, err
	}
	defer func() {
		_ = file.Close()
	}()

	var groupId *uint
	if req.GroupId > 0 {
		gid := uint(req.GroupId)
		groupId = &gid
	}
	ownerId := uint(contextx.UserId(ctx))
	data := &image.UploadRequest{
		FileName:    handler.Filename,
		FileSize:    handler.Size,
		ContentType: handler.Header.Get("Content-Type"),
		Reader:      file,
		OwnerID:     ownerId,
		GroupID:     groupId,
	}

	result, err := h.biz.Image().Upload(ctx, data)
	if err != nil {
		return &v1.UploadImageResponse{}, err
	}

	return &v1.UploadImageResponse{
		Data:    result,
		Message: "success",
	}, nil
}

// GetImage returns image metadata by ID.
func (h *handler) GetImage(ctx context.Context, req *v1.GetImageRequest) (*v1.GetImageResponse, error) {
	userID := contextx.UserId(ctx)
	img, err := h.biz.Image().GetByID(ctx, uint(req.Id), uint(userID))
	if err != nil {
		return &v1.GetImageResponse{}, errors.New("image not found")
	}

	return &v1.GetImageResponse{
		Data:    img,
		Message: "success",
	}, nil
}

// ListImages returns a paginated list of images for the current user.
func (h *handler) ListImages(ctx context.Context, req *v1.ListImagesRequest) (*v1.ListImagesResponse, error) {
	userID := contextx.UserId(ctx)
	if req.Page < 1 {
		req.Page = 1
	}
	if req.PageSize < 1 || req.PageSize > 100 {
		req.PageSize = 20
	}

	images, total, err := h.biz.Image().List(ctx, uint(userID), req.Page, req.PageSize)
	if err != nil {
		return &v1.ListImagesResponse{}, err
	}

	data := map[string]any{
		"items": images,
		"total": total,
	}

	return &v1.ListImagesResponse{
		Data:    data,
		Message: "success",
	}, nil
}

// DeleteImage deletes an image by ID.
func (h *handler) DeleteImage(ctx context.Context, req *v1.DeleteImageRequest) (*v1.DeleteImageResponse, error) {
	if err := h.biz.Image().Delete(ctx, uint(req.Id), uint(contextx.UserId(ctx))); err != nil {
		return &v1.DeleteImageResponse{}, err
	}
	return &v1.DeleteImageResponse{
		Message: "success",
	}, nil
}

// CreateGroup creates a new image group.
func (h *handler) CreateGroup(ctx context.Context, req *v1.CreateGroupRequest) (*v1.CreateGroupResponse, error) {
	userID := uint(contextx.UserId(ctx))

	var refId *uint
	if req.RefID != nil {
		refId = req.RefID
	}

	group, err := h.biz.Image().CreateGroup(
		ctx,
		userID,
		model.ImageGroupType(req.GroupType),
		req.GroupName,
		refId,
		req.RefType,
	)
	if err != nil {
		return &v1.CreateGroupResponse{}, err
	}

	return &v1.CreateGroupResponse{
		Data:    group,
		Message: "success",
	}, nil
}

// AddImagesToGroup adds images to an existing group.
func (h *handler) AddImagesToGroup(ctx context.Context, req *v1.AddImagesToGroupRequest) (*v1.AddImagesToGroupResponse, error) {
	imageIds := make([]uint, len(req.ImageIds))
	for i, id := range req.ImageIds {
		imageIds[i] = uint(id)
	}

	if err := h.biz.Image().AddImagesToGroup(ctx, uint(req.Id), imageIds); err != nil {
		return &v1.AddImagesToGroupResponse{}, err
	}

	return &v1.AddImagesToGroupResponse{
		Message: "success",
	}, nil
}

// ServeImage serves an image file.
// This endpoint can be accessed with either JWTOptions auth or access token.
// Note: This method is kept as a gin.HandlerFunc because it streams files directly.
func (h *handler) ServeImage(ctx context.Context, req *v1.ServeImageRequest) (*v1.ServeImageResponse, error) {
	var res *model.Image

	// Try to get image by access token first (for shared/public access)
	if req.Token != "" {
		img, err := h.biz.Image().GetByAccessToken(ctx, req.Token)
		if err != nil {
			//c.JSON(http.StatusNotFound, gin.H{"error": "image not found"})
			return &v1.ServeImageResponse{}, err
		}

		// Verify ID matches
		if img.ID != uint(req.Id) {
			//c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
			return &v1.ServeImageResponse{}, errors.New("access denied")
		}
		res = img
	} else {
		// Fall back to JWTOptions authentication
		userId := contextx.UserId(ctx)

		img, err := h.biz.Image().GetByID(ctx, uint(req.Id), uint(userId))
		if err != nil {
			//c.JSON(http.StatusNotFound, gin.H{"error": "image not found"})
			return &v1.ServeImageResponse{}, err
		}
		res = img
	}

	// Download the file
	reader, mimeType, err := h.biz.Image().DownloadFile(ctx, res, req.Variant)
	if err != nil {
		//c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to download image"})
		return &v1.ServeImageResponse{}, err
	}
	defer reader.Close()

	resp := contextx.Response(ctx)

	// Set headers
	resp.Header().Set("Content-Type", mimeType)
	resp.Header().Set("Cache-Control", "public, max-age=31536000") // 1 year cache

	// Stream the file
	_, err = io.Copy(resp, reader)
	return &v1.ServeImageResponse{}, err
}
