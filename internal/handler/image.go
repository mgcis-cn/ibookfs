// Package handler provides HTTP handlers for image management.
package handler

import (
	"io"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mgcis-cn/ibookfs/internal/middleware"
	"github.com/mgcis-cn/ibookfs/internal/model"
	"github.com/mgcis-cn/ibookfs/internal/service"
)

var imageService *service.ImageService

// SetImageService sets the image service instance.
func SetImageService(svc *service.ImageService) {
	imageService = svc
}

// UploadImage handles image upload requests.
func UploadImage(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// Get file from form
	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file is required"})
		return
	}

	// Open file
	file, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to open file"})
		return
	}
	defer file.Close()

	// Get optional group_id
	var groupID *uint
	if groupIDStr := c.PostForm("group_id"); groupIDStr != "" {
		if id, err := strconv.ParseUint(groupIDStr, 10, 32); err == nil {
			gid := uint(id)
			groupID = &gid
		}
	}

	req := &service.UploadRequest{
		FileName:    fileHeader.Filename,
		FileSize:    fileHeader.Size,
		ContentType: fileHeader.Header.Get("Content-Type"),
		Reader:      file,
		OwnerID:     userID,
		GroupID:     groupID,
	}

	result, err := imageService.Upload(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data":    result,
		"message": "uploaded",
	})
}

// GetImage returns image metadata by ID.
func GetImage(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	image, err := imageService.GetByID(c.Request.Context(), uint(id), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "image not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    image,
		"message": "success",
	})
}

// ListImages returns a paginated list of images for the current user.
func ListImages(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	images, total, err := imageService.List(c.Request.Context(), userID, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"items":     images,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
		"message": "success",
	})
}

// DeleteImage deletes an image by ID.
func DeleteImage(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := imageService.Delete(c.Request.Context(), uint(id), userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "deleted",
	})
}

// ServeImage serves an image file.
// This endpoint can be accessed with either JWT auth or access token.
func ServeImage(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	variant := c.DefaultQuery("variant", "original")
	token := c.Query("token")

	var image *model.Image
	ctx := c.Request.Context()

	// Try to get image by access token first (for shared/public access)
	if token != "" {
		image, err = imageService.GetByAccessToken(ctx, token)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "image not found"})
			return
		}

		// Verify ID matches
		if image.ID != uint(id) {
			c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
			return
		}
	} else {
		// Fall back to JWT authentication
		userID, exists := middleware.GetUserID(c)
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		image, err = imageService.GetByID(ctx, uint(id), userID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "image not found"})
			return
		}
	}

	// Download the file
	reader, mimeType, err := imageService.DownloadFile(ctx, image, variant)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to download image"})
		return
	}
	defer reader.Close()

	// Set headers
	c.Header("Content-Type", mimeType)
	c.Header("Cache-Control", "public, max-age=31536000") // 1 year cache

	// Stream the file
	_, _ = io.Copy(c.Writer, reader)
}

// CreateGroup creates a new image group.
func CreateGroup(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req struct {
		GroupType string  `json:"group_type" binding:"required"`
		GroupName string  `json:"group_name" binding:"required"`
		RefID     *uint   `json:"ref_id"`
		RefType   *string `json:"ref_type"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	group, err := imageService.CreateGroup(
		c.Request.Context(),
		userID,
		model.ImageGroupType(req.GroupType),
		req.GroupName,
		req.RefID,
		req.RefType,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data":    group,
		"message": "created",
	})
}

// AddImagesToGroup adds images to an existing group.
func AddImagesToGroup(c *gin.Context) {
	_, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req struct {
		ImageIDs []uint `json:"image_ids" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := imageService.AddImagesToGroup(c.Request.Context(), uint(id), req.ImageIDs); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "added",
	})
}
