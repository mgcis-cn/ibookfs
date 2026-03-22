package apiserver

import (
	"context"
	"os"

	"github.com/mgcis-cn/ibookfs/cmd/apiserver/app/options"
	"github.com/mgcis-cn/ibookfs/internal/apiserver/biz"
	"github.com/mgcis-cn/ibookfs/internal/apiserver/biz/processor"
	"github.com/mgcis-cn/ibookfs/internal/apiserver/handler"
	mw "github.com/mgcis-cn/ibookfs/internal/apiserver/middleware"
	"github.com/mgcis-cn/ibookfs/internal/apiserver/store"
	"github.com/mgcis-cn/ibookfs/internal/apiserver/worker"
	"github.com/mgcis-cn/ibookfs/internal/pkg/bootstrap"
	v1 "github.com/mgcis-cn/ibookfs/pkg/api/apiserver/v1"
	"github.com/mgcis-cn/ibookfs/pkg/authn/jwt"
	"github.com/mgcis-cn/ibookfs/pkg/authn/oauth"
	"github.com/mgcis-cn/ibookfs/pkg/database"
	"github.com/mgcis-cn/ibookfs/pkg/email"
	"github.com/mgcis-cn/ibookfs/pkg/storage"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
)

var (
	ID, _ = os.Hostname()
)

// ServerConfig holds server configuration.
type ServerConfig struct {
	cfg         *Config
	handler     v1.ApiServerHTTPServer
	middlewares []middleware.Middleware
	storagePath string // Base path for static file serving
}

// NewServerConfig creates a new ServerConfig from options.
func NewServerConfig(opts *options.ServerRunOptions) *ServerConfig {
	return &ServerConfig{
		cfg: NewConfig(opts),
	}
}

// New creates and starts a new Kratos application.
func New(server *ServerConfig) (app *kratos.App, cleanup func(), err error) {
	ctx := context.Background()
	opts := server.cfg.ServerRunOptions

	// App info
	appInfo := bootstrap.NewAppInfo(ID, opts.App.Name, opts.App.Version)
	appLogger := bootstrap.NewLogger(appInfo)

	databaseF, err := database.NewFactory(ctx, opts.Data.Database, database.WithNameFunc)
	if err != nil {
		return nil, nil, err
	}

	storageF, err := storage.NewFactory(opts.Data.Storage, storage.WithNameFunc)
	if err != nil {
		return nil, nil, err
	}

	emailF, err := email.NewFactory(opts.Email, email.WithNameFunc)
	if err != nil {
		return nil, nil, err
	}

	oauthF, err := oauth.NewFactory(opts.Auth.OAuth, oauth.WithNameFunc)
	if err != nil {
		return nil, nil, err
	}

	// Image processor
	imgProcessor := newImageProcessor(opts)
	imageWorker := worker.NewImageWorker(nil, worker.DefaultConfig(), appLogger)
	db, err := databaseF.MustGet("default")
	if err != nil {
		return nil, nil, err
	}
	oss, err := storageF.MustGet("default")
	if err != nil {
		return nil, nil, err
	}
	emailSvc, err := emailF.MustGet("default")
	if err != nil {
		return nil, nil, err
	}
	repo := store.New(db)

	// Biz layer - use OAuth Factory
	b := biz.New(opts.Auth.JWT, emailSvc, oauthF, repo, oss, imgProcessor, imageWorker)
	imageWorker.SetService(b.Image())
	imageWorker.Start()

	// Handler and HTTP server
	server.handler = handler.New(oauthF, b)
	server.storagePath = oss.GetBasePath()
	server.middlewares = mw.NewMiddlewares(&mw.Config{
		SkipAuthPaths: opts.Server.Middleware.AllowedPaths,
		JWTManager: jwt.New(
			jwt.WithSigningKey([]byte(opts.Auth.JWT.Secret)),
			jwt.WithExpired(opts.Auth.JWT.Expired.Duration),
		),
		AccountSecretService: b.AccountSecret(),
		Logger:               appLogger,
	})
	httpSrv := server.NewHTTPServer()
	appConfig := bootstrap.AppConfig{Info: appInfo, Logger: appLogger}
	app = bootstrap.NewApp(appConfig, transport.Server(httpSrv))

	cleanup = func() {
		imageWorker.Stop()
		err = databaseF.Close()
	}
	return app, cleanup, nil
}

// newImageProcessor creates image processor from options.
func newImageProcessor(opts *options.ServerRunOptions) *processor.Processor {
	cfg := processor.Config{
		BlurHashEnabled: opts.Server.Image.Processing.BlurHashEnabled,
		AllowedTypes:    opts.Server.Upload.AllowedTypes,
		MaxFileSize:     opts.Server.Upload.MaxFileSize,
	}
	for _, v := range opts.Server.Image.Processing.Variants {
		cfg.Variants = append(cfg.Variants, processor.VariantConfig{
			Name:      v.Name,
			MaxWidth:  v.MaxWidth,
			MaxHeight: v.MaxHeight,
		})
	}
	return processor.NewProcessor(cfg)
}
