// Package router sets up HTTP routes.
package router

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/mgcis-cn/ibookfs/internal/config"
	"github.com/mgcis-cn/ibookfs/internal/handler"
	"github.com/mgcis-cn/ibookfs/internal/middleware"
	"github.com/mgcis-cn/ibookfs/internal/model"
	"github.com/mgcis-cn/ibookfs/internal/pkg/util"
	"github.com/mgcis-cn/ibookfs/internal/repository"
	"github.com/mgcis-cn/ibookfs/internal/service"
	"github.com/mgcis-cn/ibookfs/internal/service/processor"
	"github.com/mgcis-cn/ibookfs/internal/service/storage"
	"github.com/mgcis-cn/ibookfs/internal/worker"
)

// App holds the application components including the worker for lifecycle management.
type App struct {
	Router  *gin.Engine
	Worker  *worker.ImageWorker
}

// Setup initializes the router with all routes and returns the App with worker.
func Setup(cfg *config.Config) *App {
	r := gin.Default()

	// CORS configuration
	r.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.Server.AllowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Static files (assets)
	r.Static("/assets", "./web/assets")

	// Serve frontend static files for all other routes (SPA)
	r.StaticFile("/", "./web/index.html")
	r.StaticFile("/favicon.ico", "./web/favicon.ico")
	r.NoRoute(func(c *gin.Context) {
		c.File("./web/index.html")
	})

	// Initialize services
	jwtManager := util.NewJWTManager(cfg.JWT.Secret, time.Duration(cfg.JWT.Expiration)*time.Hour)
	emailSvc := service.NewEmailService(&cfg.Email, cfg.Server.BaseURL)
	authSvc := service.NewAuthService(jwtManager, emailSvc)

	// OAuth configuration
	oauthConfigs := map[model.OAuthProvider]model.OAuthConfig{
		model.OAuthProviderGitHub: service.GetOAuthConfig(
			model.OAuthProviderGitHub,
			cfg.OAuth.GitHub.ClientID,
			cfg.OAuth.GitHub.ClientSecret,
			cfg.OAuth.GitHub.RedirectURL,
		),
		model.OAuthProviderGitee: service.GetOAuthConfig(
			model.OAuthProviderGitee,
			cfg.OAuth.Gitee.ClientID,
			cfg.OAuth.Gitee.ClientSecret,
			cfg.OAuth.Gitee.RedirectURL,
		),
	}
	oauthSvc := service.NewOAuthService(authSvc, oauthConfigs)

	// Initialize handlers
	authHandler := handler.NewAuthHandler(authSvc, oauthSvc, oauthConfigs, &cfg.Email)

	// Auth middleware
	authMiddleware := middleware.NewAuth(&middleware.AuthConfig{
		JWTManager:    jwtManager,
		SkipAuthPaths: cfg.Server.SkipAuthPaths,
	})

	// Initialize image storage
	imgStorage, err := storage.NewLocalStorage(
		cfg.Storage.Local.BasePath,
		cfg.Storage.Local.BaseURL,
	)
	if err != nil {
		panic("Failed to initialize image storage: " + err.Error())
	}

	// Initialize image processor
	processorCfg := processor.Config{
		BlurHashEnabled: cfg.Image.Processing.BlurHashEnabled,
		AllowedTypes:    cfg.Image.Upload.AllowedTypes,
		MaxFileSize:     cfg.Image.Upload.MaxFileSize,
	}
	for _, v := range cfg.Image.Processing.Variants {
		processorCfg.Variants = append(processorCfg.Variants, processor.VariantConfig{
			Name:      v.Name,
			MaxWidth:  v.MaxWidth,
			MaxHeight: v.MaxHeight,
		})
	}
	imgProcessor := processor.NewProcessor(processorCfg)

	// Initialize image repository
	imageRepo := repository.NewImageRepository()

	// Initialize image worker (without service initially)
	imageWorker := worker.NewImageWorker(nil, worker.DefaultConfig())

	// Initialize image service
	imageSvc := service.NewImageService(imgStorage, imgProcessor, imageRepo, imageWorker)
	handler.SetImageService(imageSvc)

	// Set the service for worker and start it
	imageWorker.SetService(imageSvc)
	imageWorker.Start()
	// Note: defer removed - worker will be stopped in main.go on app shutdown

	// Serve static files (generic storage, can be used for images, documents, etc.)
	r.Static(cfg.Storage.Local.BaseURL, cfg.Storage.Local.BasePath)

	// API v1
	v1 := r.Group("/api/v1")
	{
		// Auth routes
		handler.RegisterAuthRoutes(v1, authHandler, authMiddleware)

		// Book routes (protected)
		handler.RegisterBookRoutes(v1, authMiddleware)

		// Image routes (protected)
		handler.RegisterImageRoutes(v1, authMiddleware)
	}

	return &App{
		Router: r,
		Worker: imageWorker,
	}
}
