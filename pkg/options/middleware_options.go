package options

// MiddlewareOptions holds middleware-related configuration.
type MiddlewareOptions struct {
	AllowedOrigins []string `json:"allowed_origins" yaml:"allowed_origins" mapstructure:"allowed_origins"`
	AllowedPaths   []string `json:"allowed_paths" yaml:"allowed_paths" mapstructure:"allowed_paths"`
}

// NewMiddlewareOptions returns default middleware options.
func NewMiddlewareOptions() *MiddlewareOptions {
	return &MiddlewareOptions{
		AllowedOrigins: []string{"http://localhost:3000"},
		AllowedPaths: []string{
			"/health",
			"/swagger",
			"/openapi.json",
			"/storages/",
			"/api/v1/auth/send-code",
			"/api/v1/auth/login",
			"/api/v1/auth/register",
			"/api/v1/auth/oauth/authorize",
			"/api/v1/auth/oauth/callback",
			"/api/v1/auth/refresh",
			"/api/v1/auth/config",
		},
	}
}
