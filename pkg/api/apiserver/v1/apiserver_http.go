// Package apiserver provides HTTP server registration for the API server.
package v1

import (
	"context"

	http "github.com/go-kratos/kratos/v2/transport/http"
	contextx "github.com/mgcis-cn/ibookfs/pkg/context"
)

// ApiServerHTTPServer defines the interface for API server HTTP handlers.
type ApiServerHTTPServer interface {
	Health(ctx context.Context, req *HealthRequest) (*HealthResponse, error)

	// Images
	ListImages(ctx context.Context, req *ListImagesRequest) (*ListImagesResponse, error)
	GetImage(ctx context.Context, req *GetImageRequest) (*GetImageResponse, error)
	UploadImage(ctx context.Context, req *UploadImageRequest) (*UploadImageResponse, error)
	DeleteImage(ctx context.Context, req *DeleteImageRequest) (*DeleteImageResponse, error)
	ServeImage(ctx context.Context, req *ServeImageRequest) (*ServeImageResponse, error)

	// Image Groups
	CreateGroup(ctx context.Context, req *CreateGroupRequest) (*CreateGroupResponse, error)
	AddImagesToGroup(ctx context.Context, req *AddImagesToGroupRequest) (*AddImagesToGroupResponse, error)

	// Books
	ListBooks(ctx context.Context, req *ListBooksRequest) (*ListBooksResponse, error)
	GetBook(ctx context.Context, req *GetBookRequest) (*GetBookResponse, error)
	CreateBook(ctx context.Context, req *CreateBookRequest) (*CreateBookResponse, error)
	UpdateBook(ctx context.Context, req *UpdateBookRequest) (*UpdateBookResponse, error)
	DeleteBook(ctx context.Context, req *DeleteBookRequest) (*DeleteBookResponse, error)

	// Account Secrets
	ListAccountSecrets(ctx context.Context, req *ListAccountSecretsRequest) (*ListAccountSecretsResponse, error)
	GetAccountSecret(ctx context.Context, req *GetAccountSecretRequest) (*GetAccountSecretResponse, error)
	CreateAccountSecret(ctx context.Context, req *CreateAccountSecretRequest) (*CreateAccountSecretResponse, error)
	DeleteAccountSecret(ctx context.Context, req *DeleteAccountSecretRequest) (*DeleteAccountSecretResponse, error)
	DisableAccountSecret(ctx context.Context, req *DisableAccountSecretRequest) (*DisableAccountSecretResponse, error)
	EnableAccountSecret(ctx context.Context, req *EnableAccountSecretRequest) (*EnableAccountSecretResponse, error)

	// Auth
	SendCode(ctx context.Context, req *SendCodeRequest) (*SendCodeResponse, error)
	Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error)
	Register(ctx context.Context, req *RegisterRequest) (*RegisterResponse, error)
	RefreshToken(ctx context.Context, req *RefreshTokenRequest) (*RefreshTokenResponse, error)
	Logout(ctx context.Context, req *LogoutRequest) (*LogoutResponse, error)
	GetCurrentUser(ctx context.Context, req *GetCurrentUserRequest) (*GetCurrentUserResponse, error)
	GetLinkedAccounts(ctx context.Context, req *GetLinkedAccountsRequest) (*GetLinkedAccountsResponse, error)
	LinkOAuth(ctx context.Context, req *LinkOAuthRequest) (*LinkOAuthResponse, error)
	UnlinkOAuth(ctx context.Context, req *UnlinkOAuthRequest) (*UnlinkOAuthResponse, error)
	BindEmail(ctx context.Context, req *BindEmailRequest) (*BindEmailResponse, error)
	SetPrimaryAccount(ctx context.Context, req *SetPrimaryAccountRequest) (*SetPrimaryAccountResponse, error)
	UpdateProfile(ctx context.Context, req *UpdateProfileRequest) (*UpdateProfileResponse, error)
	GetConfig(ctx context.Context, req *GetConfigRequest) (*GetConfigResponse, error)
	OAuthAuthorize(ctx context.Context, req *OAuthAuthorizeRequest) (*OAuthAuthorizeResponse, error)
	OAuthCallback(ctx context.Context, req *OAuthCallbackRequest) (*OAuthCallbackResponse, error)
}

// RegisterApiServerHTTPServer registers the API server HTTP handlers with the given server.
func RegisterApiServerHTTPServer(s *http.Server, srv ApiServerHTTPServer) {
	r := s.Route("/")
	r.GET("/health", healthHttpHandler(srv))

	// images
	r.GET("/api/v1/images/{id}/file", serveImageHttpHandler(srv))
	r.GET("/api/v1/images/{id}", getImageHttpHandler(srv))
	r.GET("/api/v1/images", listImagesHttpHandler(srv))
	r.POST("/api/v1/images", uploadImageHttpHandler(srv))
	r.DELETE("/api/v1/images/{id}", deleteImageHttpHandler(srv))

	// imageGroup
	r.POST("/api/v1/images/groups/{id}/images", addImagesToGroupHttpHandler(srv))
	r.POST("/api/v1/images/groups", createGroupHttpHandler(srv))

	// books
	r.GET("/api/v1/books/{id}", getBookHttpHandler(srv))
	r.GET("/api/v1/books", listBooksHttpHandler(srv))
	r.POST("/api/v1/books", createBookHttpHandler(srv))
	r.PUT("/api/v1/books/{id}", updateBookHttpHandler(srv))
	r.DELETE("/api/v1/books/{id}", deleteBookHttpHandler(srv))

	// account secrets
	r.GET("/api/v1/account-secrets/{id}", getAccountSecretHttpHandler(srv))
	r.GET("/api/v1/account-secrets", listAccountSecretsHttpHandler(srv))
	r.POST("/api/v1/account-secrets/{id}/disable", disableAccountSecretHttpHandler(srv))
	r.POST("/api/v1/account-secrets/{id}/enable", enableAccountSecretHttpHandler(srv))
	r.POST("/api/v1/account-secrets", createAccountSecretHttpHandler(srv))
	r.DELETE("/api/v1/account-secrets/{id}", deleteAccountSecretHttpHandler(srv))

	// auth - public routes
	r.GET("/api/v1/auth/config", getConfigHttpHandler(srv))
	r.POST("/api/v1/auth/send-code", sendCodeHttpHandler(srv))
	r.POST("/api/v1/auth/login", loginHttpHandler(srv))
	r.POST("/api/v1/auth/register", registerHttpHandler(srv))
	r.GET("/api/v1/auth/oauth/authorize", oAuthAuthorizeHttpHandler(srv))
	r.POST("/api/v1/auth/oauth/callback", oAuthCallbackHttpHandler(srv))
	r.POST("/api/v1/auth/refresh", refreshTokenHttpHandler(srv))

	// auth - protected routes
	r.POST("/api/v1/auth/logout", logoutHttpHandler(srv))
	r.GET("/api/v1/auth/me", getCurrentUserHttpHandler(srv))
	r.GET("/api/v1/auth/accounts", getLinkedAccountsHttpHandler(srv))
	r.POST("/api/v1/auth/oauth/link", linkOAuthHttpHandler(srv))
	r.DELETE("/api/v1/auth/oauth/unlink", unlinkOAuthHttpHandler(srv))
	r.POST("/api/v1/auth/bind-email", bindEmailHttpHandler(srv))
	r.PUT("/api/v1/auth/accounts/primary", setPrimaryAccountHttpHandler(srv))
	r.PUT("/api/v1/auth/me", updateProfileHttpHandler(srv))
}

func serveImageHttpHandler(srv ApiServerHTTPServer) http.HandlerFunc {
	return func(ctx http.Context) error {
		var in ServeImageRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		if err := ctx.BindVars(&in); err != nil {
			return err
		}
		cresp := ctx.Response()
		h := ctx.Middleware(func(ctx context.Context, req any) (any, error) {
			ctx = contextx.WithResponse(ctx, cresp)
			return srv.ServeImage(ctx, req.(*ServeImageRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*AddImagesToGroupResponse)
		return ctx.Result(200, reply)
	}
}

func addImagesToGroupHttpHandler(srv ApiServerHTTPServer) http.HandlerFunc {
	return func(ctx http.Context) error {
		var in AddImagesToGroupRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		h := ctx.Middleware(func(ctx context.Context, req any) (any, error) {
			return srv.AddImagesToGroup(ctx, req.(*AddImagesToGroupRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*AddImagesToGroupResponse)
		return ctx.Result(200, reply)
	}
}

func createGroupHttpHandler(srv ApiServerHTTPServer) http.HandlerFunc {
	return func(ctx http.Context) error {
		var in CreateGroupRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		h := ctx.Middleware(func(ctx context.Context, req any) (any, error) {
			return srv.CreateGroup(ctx, req.(*CreateGroupRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*CreateGroupResponse)
		return ctx.Result(200, reply)
	}
}

func deleteImageHttpHandler(srv ApiServerHTTPServer) http.HandlerFunc {
	return func(ctx http.Context) (err error) {
		var in DeleteImageRequest
		if err := ctx.BindVars(&in); err != nil {
			return err
		}
		h := ctx.Middleware(func(ctx context.Context, req any) (any, error) {
			return srv.DeleteImage(ctx, req.(*DeleteImageRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*DeleteImageResponse)
		return ctx.Result(200, reply)
	}
}

func uploadImageHttpHandler(srv ApiServerHTTPServer) http.HandlerFunc {
	return func(ctx http.Context) (err error) {
		var in UploadImageRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		creq := ctx.Request()
		h := ctx.Middleware(func(ctx context.Context, req any) (any, error) {
			ctx = contextx.WithRequest(ctx, creq)
			return srv.UploadImage(ctx, req.(*UploadImageRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*UploadImageResponse)
		return ctx.Result(200, reply)
	}
}

func getImageHttpHandler(srv ApiServerHTTPServer) http.HandlerFunc {
	return func(ctx http.Context) error {
		var in GetImageRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		if err := ctx.BindVars(&in); err != nil {
			return err
		}
		h := ctx.Middleware(func(ctx context.Context, req any) (any, error) {
			return srv.GetImage(ctx, req.(*GetImageRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*GetImageResponse)
		return ctx.Result(200, reply)
	}
}

func listImagesHttpHandler(srv ApiServerHTTPServer) http.HandlerFunc {
	return func(ctx http.Context) error {
		var in ListImagesRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		h := ctx.Middleware(func(ctx context.Context, req any) (any, error) {
			return srv.ListImages(ctx, req.(*ListImagesRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*ListImagesResponse)
		return ctx.Result(200, reply)
	}
}

func healthHttpHandler(srv ApiServerHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in HealthRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		h := ctx.Middleware(func(ctx context.Context, req any) (any, error) {
			return srv.Health(ctx, req.(*HealthRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*HealthResponse)
		return ctx.Result(200, reply)
	}
}

// Book HTTP handlers

func listBooksHttpHandler(srv ApiServerHTTPServer) http.HandlerFunc {
	return func(ctx http.Context) error {
		var in ListBooksRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		h := ctx.Middleware(func(ctx context.Context, req any) (any, error) {
			return srv.ListBooks(ctx, req.(*ListBooksRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*ListBooksResponse)
		return ctx.Result(200, reply)
	}
}

func getBookHttpHandler(srv ApiServerHTTPServer) http.HandlerFunc {
	return func(ctx http.Context) error {
		var in GetBookRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		if err := ctx.BindVars(&in); err != nil {
			return err
		}
		h := ctx.Middleware(func(ctx context.Context, req any) (any, error) {
			return srv.GetBook(ctx, req.(*GetBookRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*GetBookResponse)
		return ctx.Result(200, reply)
	}
}

func createBookHttpHandler(srv ApiServerHTTPServer) http.HandlerFunc {
	return func(ctx http.Context) error {
		var in CreateBookRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		h := ctx.Middleware(func(ctx context.Context, req any) (any, error) {
			return srv.CreateBook(ctx, req.(*CreateBookRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*CreateBookResponse)
		return ctx.Result(200, reply)
	}
}

func updateBookHttpHandler(srv ApiServerHTTPServer) http.HandlerFunc {
	return func(ctx http.Context) error {
		var in UpdateBookRequest
		if err := ctx.BindVars(&in); err != nil {
			return err
		}
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		h := ctx.Middleware(func(ctx context.Context, req any) (any, error) {
			return srv.UpdateBook(ctx, req.(*UpdateBookRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*UpdateBookResponse)
		return ctx.Result(200, reply)
	}
}

func deleteBookHttpHandler(srv ApiServerHTTPServer) http.HandlerFunc {
	return func(ctx http.Context) error {
		var in DeleteBookRequest
		if err := ctx.BindVars(&in); err != nil {
			return err
		}
		h := ctx.Middleware(func(ctx context.Context, req any) (any, error) {
			return srv.DeleteBook(ctx, req.(*DeleteBookRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*DeleteBookResponse)
		return ctx.Result(200, reply)
	}
}

// Account Secret HTTP handlers

func listAccountSecretsHttpHandler(srv ApiServerHTTPServer) http.HandlerFunc {
	return func(ctx http.Context) error {
		var in ListAccountSecretsRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		h := ctx.Middleware(func(ctx context.Context, req any) (any, error) {
			return srv.ListAccountSecrets(ctx, req.(*ListAccountSecretsRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*ListAccountSecretsResponse)
		return ctx.Result(200, reply)
	}
}

func getAccountSecretHttpHandler(srv ApiServerHTTPServer) http.HandlerFunc {
	return func(ctx http.Context) error {
		var in GetAccountSecretRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		if err := ctx.BindVars(&in); err != nil {
			return err
		}
		h := ctx.Middleware(func(ctx context.Context, req any) (any, error) {
			return srv.GetAccountSecret(ctx, req.(*GetAccountSecretRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*GetAccountSecretResponse)
		return ctx.Result(200, reply)
	}
}

func createAccountSecretHttpHandler(srv ApiServerHTTPServer) http.HandlerFunc {
	return func(ctx http.Context) error {
		var in CreateAccountSecretRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		h := ctx.Middleware(func(ctx context.Context, req any) (any, error) {
			return srv.CreateAccountSecret(ctx, req.(*CreateAccountSecretRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*CreateAccountSecretResponse)
		return ctx.Result(200, reply)
	}
}

func deleteAccountSecretHttpHandler(srv ApiServerHTTPServer) http.HandlerFunc {
	return func(ctx http.Context) error {
		var in DeleteAccountSecretRequest
		if err := ctx.BindVars(&in); err != nil {
			return err
		}
		h := ctx.Middleware(func(ctx context.Context, req any) (any, error) {
			return srv.DeleteAccountSecret(ctx, req.(*DeleteAccountSecretRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*DeleteAccountSecretResponse)
		return ctx.Result(200, reply)
	}
}

func disableAccountSecretHttpHandler(srv ApiServerHTTPServer) http.HandlerFunc {
	return func(ctx http.Context) error {
		var in DisableAccountSecretRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		h := ctx.Middleware(func(ctx context.Context, req any) (any, error) {
			return srv.DisableAccountSecret(ctx, req.(*DisableAccountSecretRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*DisableAccountSecretResponse)
		return ctx.Result(200, reply)
	}
}

func enableAccountSecretHttpHandler(srv ApiServerHTTPServer) http.HandlerFunc {
	return func(ctx http.Context) error {
		var in EnableAccountSecretRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		h := ctx.Middleware(func(ctx context.Context, req any) (any, error) {
			return srv.EnableAccountSecret(ctx, req.(*EnableAccountSecretRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*EnableAccountSecretResponse)
		return ctx.Result(200, reply)
	}
}

// Auth HTTP handlers

func sendCodeHttpHandler(srv ApiServerHTTPServer) http.HandlerFunc {
	return func(ctx http.Context) error {
		var in SendCodeRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		h := ctx.Middleware(func(ctx context.Context, req any) (any, error) {
			return srv.SendCode(ctx, req.(*SendCodeRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*SendCodeResponse)
		return ctx.Result(200, reply)
	}
}

func loginHttpHandler(srv ApiServerHTTPServer) http.HandlerFunc {
	return func(ctx http.Context) error {
		var in LoginRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		creq := ctx.Request()
		h := ctx.Middleware(func(ctx context.Context, req any) (any, error) {
			ctx = contextx.WithRequest(ctx, creq)
			return srv.Login(ctx, req.(*LoginRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*LoginResponse)
		return ctx.Result(200, reply)
	}
}

func registerHttpHandler(srv ApiServerHTTPServer) http.HandlerFunc {
	return func(ctx http.Context) error {
		var in RegisterRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		creq := ctx.Request()
		h := ctx.Middleware(func(ctx context.Context, req any) (any, error) {
			ctx = contextx.WithRequest(ctx, creq)
			return srv.Register(ctx, req.(*RegisterRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*RegisterResponse)
		return ctx.Result(200, reply)
	}
}

func refreshTokenHttpHandler(srv ApiServerHTTPServer) http.HandlerFunc {
	return func(ctx http.Context) error {
		var in RefreshTokenRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		h := ctx.Middleware(func(ctx context.Context, req any) (any, error) {
			return srv.RefreshToken(ctx, req.(*RefreshTokenRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*RefreshTokenResponse)
		return ctx.Result(200, reply)
	}
}

func logoutHttpHandler(srv ApiServerHTTPServer) http.HandlerFunc {
	return func(ctx http.Context) error {
		var in LogoutRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		h := ctx.Middleware(func(ctx context.Context, req any) (any, error) {
			return srv.Logout(ctx, req.(*LogoutRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*LogoutResponse)
		return ctx.Result(200, reply)
	}
}

func getCurrentUserHttpHandler(srv ApiServerHTTPServer) http.HandlerFunc {
	return func(ctx http.Context) error {
		var in GetCurrentUserRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		h := ctx.Middleware(func(ctx context.Context, req any) (any, error) {
			return srv.GetCurrentUser(ctx, req.(*GetCurrentUserRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*GetCurrentUserResponse)
		return ctx.Result(200, reply)
	}
}

func getLinkedAccountsHttpHandler(srv ApiServerHTTPServer) http.HandlerFunc {
	return func(ctx http.Context) error {
		var in GetLinkedAccountsRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		h := ctx.Middleware(func(ctx context.Context, req any) (any, error) {
			return srv.GetLinkedAccounts(ctx, req.(*GetLinkedAccountsRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*GetLinkedAccountsResponse)
		return ctx.Result(200, reply)
	}
}

func linkOAuthHttpHandler(srv ApiServerHTTPServer) http.HandlerFunc {
	return func(ctx http.Context) error {
		var in LinkOAuthRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		h := ctx.Middleware(func(ctx context.Context, req any) (any, error) {
			return srv.LinkOAuth(ctx, req.(*LinkOAuthRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*LinkOAuthResponse)
		return ctx.Result(200, reply)
	}
}

func unlinkOAuthHttpHandler(srv ApiServerHTTPServer) http.HandlerFunc {
	return func(ctx http.Context) error {
		var in UnlinkOAuthRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		h := ctx.Middleware(func(ctx context.Context, req any) (any, error) {
			return srv.UnlinkOAuth(ctx, req.(*UnlinkOAuthRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*UnlinkOAuthResponse)
		return ctx.Result(200, reply)
	}
}

func bindEmailHttpHandler(srv ApiServerHTTPServer) http.HandlerFunc {
	return func(ctx http.Context) error {
		var in BindEmailRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		h := ctx.Middleware(func(ctx context.Context, req any) (any, error) {
			return srv.BindEmail(ctx, req.(*BindEmailRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*BindEmailResponse)
		return ctx.Result(200, reply)
	}
}

func setPrimaryAccountHttpHandler(srv ApiServerHTTPServer) http.HandlerFunc {
	return func(ctx http.Context) error {
		var in SetPrimaryAccountRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		h := ctx.Middleware(func(ctx context.Context, req any) (any, error) {
			return srv.SetPrimaryAccount(ctx, req.(*SetPrimaryAccountRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*SetPrimaryAccountResponse)
		return ctx.Result(200, reply)
	}
}

func updateProfileHttpHandler(srv ApiServerHTTPServer) http.HandlerFunc {
	return func(ctx http.Context) error {
		var in UpdateProfileRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		h := ctx.Middleware(func(ctx context.Context, req any) (any, error) {
			return srv.UpdateProfile(ctx, req.(*UpdateProfileRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*UpdateProfileResponse)
		return ctx.Result(200, reply)
	}
}

func getConfigHttpHandler(srv ApiServerHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in GetConfigRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, "/api/v1/auth/config")
		h := ctx.Middleware(func(ctx context.Context, req any) (any, error) {
			return srv.GetConfig(ctx, req.(*GetConfigRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*GetConfigResponse)
		return ctx.Result(200, reply)
	}
}

func oAuthAuthorizeHttpHandler(srv ApiServerHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in OAuthAuthorizeRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		h := ctx.Middleware(func(ctx context.Context, req any) (any, error) {
			return srv.OAuthAuthorize(ctx, req.(*OAuthAuthorizeRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*OAuthAuthorizeResponse)
		return ctx.Result(200, reply)
	}
}

func oAuthCallbackHttpHandler(srv ApiServerHTTPServer) http.HandlerFunc {
	return func(ctx http.Context) error {
		var in OAuthCallbackRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		creq := ctx.Request()
		h := ctx.Middleware(func(ctx context.Context, req any) (any, error) {
			ctx = contextx.WithRequest(ctx, creq)
			return srv.OAuthCallback(ctx, req.(*OAuthCallbackRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*OAuthCallbackResponse)
		return ctx.Result(200, reply)
	}
}
