package v1

type HealthRequest struct {
}

type HealthResponse struct {
	Status string `json:"status"`
}

type ListImagesRequest struct {
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
}

type ListImagesResponse struct {
	Data    any    `json:"data"`
	Message string `json:"message"`
}

type GetImageRequest struct {
	Id int `json:"id"`
}

type GetImageResponse struct {
	Data    any    `json:"data"`
	Message string `json:"message"`
}

type UploadImageRequest struct {
	BookId int `json:"book_id"`
}

type UploadImageResponse struct {
	Data    any    `json:"data"`
	Message string `json:"message"`
}

type DeleteImageRequest struct {
	Id int `json:"id"`
}

type DeleteImageResponse struct {
	Message string `json:"message"`
}

type ServeImageRequest struct {
	Id      int    `json:"id"`
	Variant string `json:"variant"`
	Token   string `json:"token"`
}

type ServeImageResponse struct{}

type CreateGroupRequest struct {
	GroupType string  `json:"group_type"`
	GroupName string  `json:"group_name"`
	RefID     *uint   `json:"ref_id"`
	RefType   *string `json:"ref_type"`
}

type CreateGroupResponse struct {
	Data    any    `json:"data"`
	Message string `json:"message"`
}

type AddImagesToGroupRequest struct {
	Id       int   `json:"id"`
	ImageIds []int `json:"image_ids"`
}

type AddImagesToGroupResponse struct {
	Message string `json:"message"`
}

type ListBooksRequest struct {
	Page     int    `json:"page"`
	PageSize int    `json:"page_size"`
	Status   string `json:"status"`
	Search   string `json:"search"`
}

type ListBooksResponse struct {
	Data    any    `json:"data"`
	Total   int64  `json:"total"`
	Message string `json:"message"`
}

type GetBookRequest struct {
	Id int `json:"id"`
}

type GetBookResponse struct {
	Data    any    `json:"data"`
	Message string `json:"message"`
}

type CreateBookRequest struct {
	Title     string `json:"title"`
	Author    string `json:"author"`
	ISBN      string `json:"isbn"`
	Publisher string `json:"publisher"`
	Year      int    `json:"year"`
	Pages     int    `json:"pages"`
}

type CreateBookResponse struct {
	Data    any    `json:"data"`
	Message string `json:"message"`
}

type UpdateBookRequest struct {
	Id     int    `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
	ISBN   string `json:"isbn"`
	Pages  int    `json:"pages"`
}

type UpdateBookResponse struct {
	Data    any    `json:"data"`
	Message string `json:"message"`
}

type DeleteBookRequest struct {
	Id int `json:"id"`
}

type DeleteBookResponse struct {
	Message string `json:"message"`
}

type ListAccountSecretsRequest struct {
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
}

type ListAccountSecretsResponse struct {
	Data     any    `json:"data"`
	Total    int64  `json:"total"`
	Page     int    `json:"page"`
	PageSize int    `json:"page_size"`
	Message  string `json:"message"`
}

type GetAccountSecretRequest struct {
	Id int `json:"id"`
}

type GetAccountSecretResponse struct {
	Data    any    `json:"data"`
	Message string `json:"message"`
}

type CreateAccountSecretRequest struct {
	Name          string   `json:"name"`
	Scope         string   `json:"scope"`
	ResourceID    *uint    `json:"resource_id"`
	ResourceType  *string  `json:"resource_type"`
	Permissions   []string `json:"permissions"`
	ExpiresInDays *int     `json:"expires_in_days"`
}

type CreateAccountSecretResponse struct {
	Data    any `json:"data"`
	Message any `json:"message"`
}

type DeleteAccountSecretRequest struct {
	Id int `json:"id"`
}

type DeleteAccountSecretResponse struct {
	Message string `json:"message"`
}

type DisableAccountSecretRequest struct {
	Id int `json:"id"`
}

type DisableAccountSecretResponse struct {
	Message string `json:"message"`
}

type EnableAccountSecretRequest struct {
	Id int `json:"id"`
}

type EnableAccountSecretResponse struct {
	Message string `json:"message"`
}

type SendCodeRequest struct {
	Email string `json:"email" binding:"required,email"`
	Type  string `json:"type" binding:"required"`
}

type SendCodeResponse struct {
	ExpiresIn int    `json:"expires_in"`
	Message   string `json:"message"`
}

type LoginRequest struct {
	Email string `json:"email" binding:"required,email"`
	Code  string `json:"code" binding:"required,len=6"`
}

type LoginResponse struct {
	User           any    `json:"user"`
	Token          string `json:"token"`
	RefreshToken   string `json:"refresh_token,omitempty"`
	ExpiresIn      int64  `json:"expires_in"`
	LinkedAccounts []any  `json:"linked_accounts,omitempty"`
}

type RegisterRequest struct {
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	Email     string `json:"email" binding:"required,email"`
	Code      string `json:"code" binding:"required,len=6"`
}

type RegisterResponse struct {
	User      any    `json:"user"`
	Token     string `json:"token"`
	ExpiresIn int64  `json:"expires_in"`
	IsNewUser bool   `json:"is_new_user,omitempty"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type RefreshTokenResponse struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
}

type LogoutResponse struct {
	Message string `json:"message"`
}

type GetCurrentUserResponse struct {
	User           any   `json:"user"`
	LinkedAccounts []any `json:"linked_accounts"`
}

type GetLinkedAccountsResponse struct {
	Accounts []any `json:"accounts"`
}

type LinkOAuthRequest struct {
	Provider string `json:"provider" binding:"required"`
	Code     string `json:"code" binding:"required"`
}

type LinkOAuthResponse struct {
	Message string `json:"message"`
}

type UnlinkOAuthRequest struct {
	Provider       string  `json:"provider" binding:"required"`
	ProviderUserID *string `json:"provider_user_id"`
}

type UnlinkOAuthResponse struct {
	Message string `json:"message"`
}

type BindEmailRequest struct {
	Email string `json:"email" binding:"required,email"`
	Code  string `json:"code" binding:"required,len=6"`
}

type BindEmailResponse struct {
	User any `json:"user"`
}

type SetPrimaryAccountRequest struct {
	Provider       string  `json:"provider" binding:"required"`
	ProviderUserID *string `json:"provider_user_id"`
}

type SetPrimaryAccountResponse struct {
	Message string `json:"message"`
}

type UpdateProfileRequest struct {
	FirstName   *string `json:"first_name"`
	LastName    *string `json:"last_name"`
	DisplayName *string `json:"display_name"`
	Language    *string `json:"language"`
	Timezone    *string `json:"timezone"`
}

type UpdateProfileResponse struct {
	User any `json:"user"`
}

type GetConfigResponse struct {
	OauthProviders []string               `json:"oauth_providers"`
	EmailEnabled   bool                   `json:"email_enabled"`
	OAuthEnabled   bool                   `json:"oauth_enabled"`
	OAuth          map[string]interface{} `json:"oauth"`
}

// Missing request types

type LogoutRequest struct {
	Token string `json:"token"`
}

type GetCurrentUserRequest struct {
	Token string `json:"token"`
}

type GetLinkedAccountsRequest struct {
	Token string `json:"token"`
}

type GetConfigRequest struct{}

type OAuthAuthorizeRequest struct {
	Provider    string `json:"provider" binding:"required"`
	RedirectURI string `json:"redirect_uri"`
	State       string `json:"state"`
}

type OAuthAuthorizeResponse struct {
	AuthorizeURL string `json:"authorize_url"`
	State        string `json:"state"`
}

type OAuthCallbackRequest struct {
	Provider string `json:"provider" binding:"required"`
	Code     string `json:"code" binding:"required"`
	State    string `json:"state" binding:"required"`
}

type OAuthCallbackResponse struct {
	Data any
}
