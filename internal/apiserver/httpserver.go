package apiserver

import (
	"encoding/json"
	nethttp "net/http"

	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/gorilla/handlers"

	v1 "github.com/mgcis-cn/ibookfs/pkg/api/apiserver/v1"
)

// apiResponse wraps all successful responses for the frontend.
type apiResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
}

// apiErrorResponse wraps all error responses for the frontend.
type apiErrorResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// responseEncoder wraps handler responses with {success: true, data: ...}.
func responseEncoder(w nethttp.ResponseWriter, r *nethttp.Request, v interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(apiResponse{
		Success: true,
		Data:    v,
	})
}

// errorEncoder wraps handler errors with {success: false, message: ...}.
func errorEncoder(w nethttp.ResponseWriter, r *nethttp.Request, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(nethttp.StatusBadRequest)
	json.NewEncoder(w).Encode(apiErrorResponse{
		Success: false,
		Message: err.Error(),
	})
}

// NewHTTPServer creates a new HTTP server with configuration.
func (c *ServerConfig) NewHTTPServer() *http.Server {
	opts := []http.ServerOption{
		http.ResponseEncoder(responseEncoder),
		http.ErrorEncoder(errorEncoder),
		http.Filter(handlers.CORS(
			handlers.AllowedHeaders([]string{
				"X-Requested-With",
				"Content-Type",
				"Authorization",
				"X-Idempotent-ID",
				"X-Account-Key",
				"X-Signature",
				"X-Expires",
			}),
			handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"}),
			handlers.AllowedOrigins([]string{"*"}),
		)),
	}

	// Apply Kratos middleware
	if len(c.middlewares) > 0 {
		opts = append(opts, http.Middleware(c.middlewares...))
	}

	if c.cfg.Server.HTTP.Addr != "" {
		opts = append(opts, http.Address(c.cfg.Server.HTTP.Addr))
	}
	if c.cfg.Server.HTTP.Timeout.Duration != 0 {
		opts = append(opts, http.Timeout(c.cfg.Server.HTTP.Timeout.Duration))
	}
	srv := http.NewServer(opts...)
	v1.RegisterApiServerHTTPServer(srv, c.handler)
	return srv
}
