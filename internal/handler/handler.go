// Package handler contains HTTP request handlers.
package handler

import "github.com/gin-gonic/gin"

// RegisterBookRoutes registers all book-related routes.
func RegisterBookRoutes(r *gin.RouterGroup, authMiddleware gin.HandlerFunc) {
	books := r.Group("/books")
	books.Use(authMiddleware)
	{
		books.GET("", ListBooks)
		books.GET("/:id", GetBook)
		books.POST("", CreateBook)
		books.PUT("/:id", UpdateBook)
		books.DELETE("/:id", DeleteBook)
	}
}

// RegisterImageRoutes registers all image-related routes.
func RegisterImageRoutes(r *gin.RouterGroup, authMiddleware gin.HandlerFunc) {
	images := r.Group("/images")
	images.Use(authMiddleware)
	{
		images.GET("", ListImages)
		images.GET("/:id", GetImage)
		images.POST("", UploadImage)
		images.DELETE("/:id", DeleteImage)
	}

	// Image file serving - supports both JWT auth and access token
	images.GET("/:id/file", ServeImage)

	// Group management
	groups := r.Group("/images/groups")
	groups.Use(authMiddleware)
	{
		groups.POST("", CreateGroup)
		groups.POST("/:id/images", AddImagesToGroup)
	}
}
