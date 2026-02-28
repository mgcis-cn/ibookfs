package biz

import (
	"github.com/mgcis-cn/ibookfs/internal/apiserver/biz/accountsecret"
	"github.com/mgcis-cn/ibookfs/internal/apiserver/biz/auth"
	"github.com/mgcis-cn/ibookfs/internal/apiserver/biz/book"
	"github.com/mgcis-cn/ibookfs/internal/apiserver/biz/email"
	"github.com/mgcis-cn/ibookfs/internal/apiserver/biz/image"
	oauthBiz "github.com/mgcis-cn/ibookfs/internal/apiserver/biz/oauth"
	"github.com/mgcis-cn/ibookfs/internal/apiserver/biz/processor"
	"github.com/mgcis-cn/ibookfs/internal/apiserver/store"
	"github.com/mgcis-cn/ibookfs/pkg/authn/jwt"
	"github.com/mgcis-cn/ibookfs/pkg/authn/oauth"
	emailSources "github.com/mgcis-cn/ibookfs/pkg/email/sources"
	"github.com/mgcis-cn/ibookfs/pkg/options"
	"github.com/mgcis-cn/ibookfs/pkg/storage/sources"
)

type Biz interface {
	AccountSecret() accountsecret.AccountSecretBiz
	Auth() auth.AuthBiz
	Book() book.BookBiz
	Email() email.EmailBiz
	Image() image.ImageBiz
	OAuth() oauthBiz.OAuthBiz
}

type biz struct {
	oauthFactory *oauth.Factory
	repo         store.IStore
	storage      sources.Storage
	imgProcessor *processor.Processor
	imageWorker  image.ImageWorker
	jwtManager   *jwt.Auth
	email        emailSources.Email
}

var _ Biz = (*biz)(nil)

func New(
	jwtOptions *options.JWTOptions,
	email emailSources.Email,
	oauthFactory *oauth.Factory,
	repo store.IStore,
	storage sources.Storage,
	imgProcessor *processor.Processor,
	imageWorker image.ImageWorker,
) *biz {
	// Create JWT manager
	jwtManager := jwt.New(
		jwt.WithSigningKey([]byte(jwtOptions.Secret)),
		jwt.WithExpired(jwtOptions.Expired.Duration),
	)

	return &biz{
		oauthFactory: oauthFactory,
		repo:         repo,
		storage:      storage,
		imgProcessor: imgProcessor,
		imageWorker:  imageWorker,
		jwtManager:   jwtManager,
		email:        email,
	}
}

func (b *biz) AccountSecret() accountsecret.AccountSecretBiz {
	return accountsecret.NewAccountSecretBiz(b.repo)
}

func (b *biz) Auth() auth.AuthBiz {
	return auth.NewAuthBiz(b.jwtManager, b.Email(), b.repo)
}

func (b *biz) Book() book.BookBiz {
	return book.NewBookBiz(b.repo)
}

func (b *biz) Email() email.EmailBiz {
	return email.NewEmailBiz(b.email, "http://localhost:3000")
}

func (b *biz) Image() image.ImageBiz {
	return image.NewImageBiz(b.storage, b.imgProcessor, b.repo, b.imageWorker)
}

func (b *biz) OAuth() oauthBiz.OAuthBiz {
	return oauthBiz.NewOAuthBiz(b.Auth(), b.oauthFactory, b.repo)
}
