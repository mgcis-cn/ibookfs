package v1

import (
	"github.com/getkin/kin-openapi/openapi3"
)

// BuildOpenAPISpec constructs the full OpenAPI 3.0 specification programmatically.
func BuildOpenAPISpec() *openapi3.T {
	doc := &openapi3.T{
		OpenAPI: "3.0.3",
		Info: &openapi3.Info{
			Title:       "iBookFS API",
			Description: "iBookFS - Book and Image Management Service API",
			Version:     "1.0.0",
		},
		Servers: openapi3.Servers{
			{URL: "/", Description: "Default"},
		},
		Tags: openapi3.Tags{
			{Name: "Health", Description: "Health check"},
			{Name: "Auth", Description: "Authentication and user management"},
			{Name: "Books", Description: "Book management"},
			{Name: "Images", Description: "Image upload and management"},
			{Name: "ImageGroups", Description: "Image group management"},
			{Name: "AccountSecrets", Description: "API key / account secret management"},
		},
	}

	doc.Components = &openapi3.Components{
		SecuritySchemes: openapi3.SecuritySchemes{
			"BearerAuth": &openapi3.SecuritySchemeRef{
				Value: openapi3.NewSecurityScheme().WithType("http").WithScheme("bearer").WithBearerFormat("JWT"),
			},
		},
		Schemas: oapiBuildSchemas(),
	}

	doc.Paths = oapiBuildPaths()

	return doc
}

// ---- Schema helpers ----

func oapiStrProp(desc string) *openapi3.Schema {
	s := openapi3.NewStringSchema()
	s.Description = desc
	return s
}

func oapiIntProp(desc string) *openapi3.Schema {
	s := openapi3.NewIntegerSchema()
	s.Description = desc
	return s
}

func oapiInt64Prop(desc string) *openapi3.Schema {
	s := openapi3.NewInt64Schema()
	s.Description = desc
	return s
}

func oapiBoolProp(desc string) *openapi3.Schema {
	s := &openapi3.Schema{Type: &openapi3.Types{"boolean"}}
	s.Description = desc
	return s
}

func oapiDateProp(desc string) *openapi3.Schema {
	s := openapi3.NewStringSchema()
	s.Format = "date-time"
	s.Description = desc
	return s
}

func oapiRef(name string) *openapi3.SchemaRef {
	return openapi3.NewSchemaRef("#/components/schemas/"+name, nil)
}

func oapiArrayOf(itemRef *openapi3.SchemaRef) *openapi3.Schema {
	return &openapi3.Schema{
		Type:  &openapi3.Types{"array"},
		Items: itemRef,
	}
}

func oapiJSONContent(schemaRef *openapi3.SchemaRef) openapi3.Content {
	return openapi3.Content{
		"application/json": &openapi3.MediaType{Schema: schemaRef},
	}
}

func oapiJSONBody(schemaRef *openapi3.SchemaRef, required bool) *openapi3.RequestBodyRef {
	return &openapi3.RequestBodyRef{
		Value: &openapi3.RequestBody{
			Required: required,
			Content:  oapiJSONContent(schemaRef),
		},
	}
}

func oapiOKResponse(desc string, schemaRef *openapi3.SchemaRef) *openapi3.ResponseRef {
	r := openapi3.NewResponse().WithDescription(desc)
	if schemaRef != nil {
		r.Content = oapiJSONContent(schemaRef)
	}
	return &openapi3.ResponseRef{Value: r}
}

func oapiErrResponse(desc string) *openapi3.ResponseRef {
	return oapiOKResponse(desc, oapiRef("ErrorResponse"))
}

func oapiBearerSecurity() *openapi3.SecurityRequirements {
	s := openapi3.SecurityRequirements{
		{"BearerAuth": {}},
	}
	return &s
}

func oapiPathParam(name, desc string) *openapi3.ParameterRef {
	return &openapi3.ParameterRef{
		Value: &openapi3.Parameter{
			Name: name, In: "path", Required: true,
			Schema:      &openapi3.SchemaRef{Value: openapi3.NewIntegerSchema()},
			Description: desc,
		},
	}
}

func oapiQueryParam(name, desc string, schema *openapi3.Schema) *openapi3.ParameterRef {
	return &openapi3.ParameterRef{
		Value: &openapi3.Parameter{
			Name: name, In: "query",
			Schema:      &openapi3.SchemaRef{Value: schema},
			Description: desc,
		},
	}
}

func oapiQueryParamRequired(name, desc string, schema *openapi3.Schema) *openapi3.ParameterRef {
	return &openapi3.ParameterRef{
		Value: &openapi3.Parameter{
			Name: name, In: "query", Required: true,
			Schema:      &openapi3.SchemaRef{Value: schema},
			Description: desc,
		},
	}
}

// ---- Schemas ----

func oapiBuildSchemas() openapi3.Schemas {
	return openapi3.Schemas{
		"ErrorResponse":   {Value: oapiSchemaErrorResponse()},
		"MessageResponse": {Value: oapiSchemaMessageResponse()},

		// Auth
		"User":          {Value: oapiSchemaUser()},
		"LinkedAccount": {Value: oapiSchemaLinkedAccount()},

		// Books
		"Book": {Value: oapiSchemaBook()},

		// Images
		"Image":      {Value: oapiSchemaImage()},
		"ImageGroup": {Value: oapiSchemaImageGroup()},

		// Account Secrets
		"AccountSecret": {Value: oapiSchemaAccountSecret()},
	}
}

func oapiSchemaErrorResponse() *openapi3.Schema {
	return openapi3.NewObjectSchema().
		WithProperty("code", oapiIntProp("Application error code")).
		WithProperty("message", oapiStrProp("Error message"))
}

func oapiSchemaMessageResponse() *openapi3.Schema {
	return openapi3.NewObjectSchema().
		WithProperty("message", oapiStrProp(""))
}

func oapiSchemaUser() *openapi3.Schema {
	return openapi3.NewObjectSchema().
		WithProperty("id", oapiIntProp("")).
		WithProperty("email", oapiStrProp("")).
		WithProperty("first_name", oapiStrProp("")).
		WithProperty("last_name", oapiStrProp("")).
		WithProperty("display_name", oapiStrProp("")).
		WithProperty("avatar_url", oapiStrProp("")).
		WithProperty("language", oapiStrProp("")).
		WithProperty("timezone", oapiStrProp("")).
		WithProperty("status", oapiStrProp("")).
		WithProperty("created_at", oapiDateProp("")).
		WithProperty("updated_at", oapiDateProp(""))
}

func oapiSchemaLinkedAccount() *openapi3.Schema {
	return openapi3.NewObjectSchema().
		WithProperty("provider", oapiStrProp("")).
		WithProperty("provider_user_id", oapiStrProp("")).
		WithProperty("provider_username", oapiStrProp("")).
		WithProperty("provider_email", oapiStrProp("")).
		WithProperty("is_primary", oapiBoolProp("")).
		WithProperty("linked_at", oapiDateProp("")).
		WithProperty("last_used_at", oapiDateProp(""))
}

func oapiSchemaBook() *openapi3.Schema {
	return openapi3.NewObjectSchema().
		WithProperty("id", oapiIntProp("")).
		WithProperty("user_id", oapiIntProp("")).
		WithProperty("title", oapiStrProp("")).
		WithProperty("author", oapiStrProp("")).
		WithProperty("isbn", oapiStrProp("")).
		WithProperty("publisher", oapiStrProp("")).
		WithProperty("year", oapiIntProp("")).
		WithProperty("pages", oapiIntProp("")).
		WithProperty("uploaded_pages", oapiIntProp("")).
		WithProperty("status", oapiStrProp("")).
		WithProperty("cover", oapiStrProp("")).
		WithProperty("created_at", oapiDateProp("")).
		WithProperty("updated_at", oapiDateProp(""))
}

func oapiSchemaImage() *openapi3.Schema {
	return openapi3.NewObjectSchema().
		WithProperty("id", oapiIntProp("")).
		WithProperty("owner_id", oapiIntProp("")).
		WithProperty("file_name", oapiStrProp("")).
		WithProperty("original_name", oapiStrProp("")).
		WithProperty("file_size", oapiInt64Prop("")).
		WithProperty("mime_type", oapiStrProp("")).
		WithProperty("width", oapiIntProp("")).
		WithProperty("height", oapiIntProp("")).
		WithProperty("blurhash", oapiStrProp("")).
		WithProperty("storage_path", oapiStrProp("")).
		WithProperty("access_token", oapiStrProp("")).
		WithProperty("status", oapiStrProp("")).
		WithProperty("created_at", oapiDateProp("")).
		WithProperty("updated_at", oapiDateProp(""))
}

func oapiSchemaImageGroup() *openapi3.Schema {
	return openapi3.NewObjectSchema().
		WithProperty("id", oapiIntProp("")).
		WithProperty("owner_id", oapiIntProp("")).
		WithProperty("group_type", oapiStrProp("")).
		WithProperty("group_name", oapiStrProp("")).
		WithProperty("ref_id", oapiIntProp("")).
		WithProperty("ref_type", oapiStrProp("")).
		WithProperty("created_at", oapiDateProp("")).
		WithProperty("updated_at", oapiDateProp(""))
}

func oapiSchemaAccountSecret() *openapi3.Schema {
	return openapi3.NewObjectSchema().
		WithProperty("id", oapiIntProp("")).
		WithProperty("name", oapiStrProp("")).
		WithProperty("account_key", oapiStrProp("")).
		WithProperty("scope", oapiStrProp("")).
		WithProperty("resource_id", oapiIntProp("")).
		WithProperty("resource_type", oapiStrProp("")).
		WithProperty("permissions", oapiArrayOf(&openapi3.SchemaRef{Value: openapi3.NewStringSchema()})).
		WithProperty("is_enabled", oapiBoolProp("")).
		WithProperty("expires_at", oapiDateProp("")).
		WithProperty("created_at", oapiDateProp("")).
		WithProperty("last_used_at", oapiDateProp(""))
}
