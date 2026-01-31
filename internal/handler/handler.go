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
