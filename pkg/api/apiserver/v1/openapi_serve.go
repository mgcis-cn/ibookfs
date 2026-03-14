package v1

import (
	"net/http"
	"sync"
)

var (
	openapiSpecOnce sync.Once
	openapiSpecJSON []byte
)

// SetOpenAPISpec sets the generated OpenAPI spec JSON bytes (called once at startup).
func SetOpenAPISpec(jsonBytes []byte) {
	openapiSpecOnce.Do(func() {
		openapiSpecJSON = jsonBytes
	})
}

// GetOpenAPISpec returns the generated OpenAPI spec JSON bytes.
func GetOpenAPISpec() []byte {
	return openapiSpecJSON
}

// OpenAPISwaggerUIHandler returns an http.HandlerFunc that serves the Swagger UI page.
func OpenAPISwaggerUIHandler(specURL string) http.HandlerFunc {
	html := []byte(`<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <title>iBookFS API - Swagger UI</title>
  <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5/swagger-ui.css">
  <style>
    html { box-sizing: border-box; overflow-y: scroll; }
    *, *:before, *:after { box-sizing: inherit; }
    body { margin: 0; background: #fafafa; }
  </style>
</head>
<body>
  <div id="swagger-ui"></div>
  <script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
  <script>
    SwaggerUIBundle({
      url: "` + specURL + `",
      dom_id: '#swagger-ui',
      deepLinking: true,
      presets: [
        SwaggerUIBundle.presets.apis,
        SwaggerUIBundle.SwaggerUIStandalonePreset
      ],
      layout: "BaseLayout"
    });
  </script>
</body>
</html>`)
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write(html)
	}
}

// OpenAPISpecHandler returns an http.HandlerFunc that serves the OpenAPI spec as JSON.
func OpenAPISpecHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Write(openapiSpecJSON)
	}
}
