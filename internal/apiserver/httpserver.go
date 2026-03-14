package apiserver

import (
	"encoding/json"
	"log"
	nethttp "net/http"

	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/gorilla/handlers"

	v1 "github.com/mgcis-cn/ibookfs/pkg/api/apiserver/v1"
	pkgerr "github.com/mgcis-cn/ibookfs/pkg/errors"
)

// apiErrorResponse wraps all error responses for the frontend.
type apiErrorResponse struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message"`
}

// responseEncoder encodes handler responses directly as JSON.
func responseEncoder(w nethttp.ResponseWriter, r *nethttp.Request, v interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(v)
}

// errorEncoder encodes handler errors as {code, message}.
// It extracts the HTTP status and error code from pkgerr.Error if available;
// otherwise falls back to 400 Bad Request.
func errorEncoder(w nethttp.ResponseWriter, r *nethttp.Request, err error) {
	w.Header().Set("Content-Type", "application/json")

	httpStatus := nethttp.StatusBadRequest
	errCode := 0
	message := err.Error()

	if e := pkgerr.FromError(err); e != nil {
		httpStatus = e.HTTPStatus()
		errCode = e.Code()
		message = e.Message()
	}

	w.WriteHeader(httpStatus)
	json.NewEncoder(w).Encode(apiErrorResponse{
		Code:    errCode,
		Message: message,
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
			handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS", "DELETE"}),
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

	r := srv.Route("/")

	// Generate OpenAPI spec at startup
	doc := v1.BuildOpenAPISpec()
	specJSON, err := doc.MarshalJSON()
	if err != nil {
		log.Printf("[WARN] failed to generate OpenAPI spec: %v", err)
	} else {
		v1.SetOpenAPISpec(specJSON)
	}

	// Serve OpenAPI spec and Swagger UI
	r.GET("/openapi.json", func(ctx http.Context) error {
		v1.OpenAPISpecHandler()(ctx.Response(), ctx.Request())
		return nil
	})
	r.GET("/swagger", func(ctx http.Context) error {
		v1.OpenAPISwaggerUIHandler("/openapi.json")(ctx.Response(), ctx.Request())
		return nil
	})

	return srv
}
