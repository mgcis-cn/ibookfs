package v1

import (
	"github.com/getkin/kin-openapi/openapi3"
)

func oapiBuildPaths() *openapi3.Paths {
	paths := openapi3.NewPaths()

	// Health (served at root, not under /api/v1)
	paths.Set("/health", oapiHealthPath())

	// Auth - public
	paths.Set("/api/v1/auth/config", oapiAuthConfigPath())
	paths.Set("/api/v1/auth/send-code", oapiAuthSendCodePath())
	paths.Set("/api/v1/auth/login", oapiAuthLoginPath())
	paths.Set("/api/v1/auth/register", oapiAuthRegisterPath())
	paths.Set("/api/v1/auth/oauth/authorize", oapiAuthOAuthAuthorizePath())
	paths.Set("/api/v1/auth/oauth/callback", oapiAuthOAuthCallbackPath())
	paths.Set("/api/v1/auth/refresh", oapiAuthRefreshPath())

	// Auth - protected
	paths.Set("/api/v1/auth/logout", oapiAuthLogoutPath())
	paths.Set("/api/v1/auth/me", oapiAuthMePath())
	paths.Set("/api/v1/auth/accounts", oapiAuthAccountsPath())
	paths.Set("/api/v1/auth/accounts/primary", oapiAuthAccountsPrimaryPath())
	paths.Set("/api/v1/auth/oauth/link", oapiAuthOAuthLinkPath())
	paths.Set("/api/v1/auth/oauth/unlink", oapiAuthOAuthUnlinkPath())
	paths.Set("/api/v1/auth/bind-email", oapiAuthBindEmailPath())

	// Books
	paths.Set("/api/v1/books", oapiBooksPath())
	paths.Set("/api/v1/books/{id}", oapiBookByIDPath())

	// Images
	paths.Set("/api/v1/images", oapiImagesPath())
	paths.Set("/api/v1/images/{id}", oapiImageByIDPath())
	paths.Set("/api/v1/images/{id}/file", oapiImageFilePath())

	// Image Groups
	paths.Set("/api/v1/images/groups", oapiImageGroupsPath())
	paths.Set("/api/v1/images/groups/{id}/images", oapiImageGroupImagesPath())

	// Account Secrets
	paths.Set("/api/v1/account-secrets", oapiAccountSecretsPath())
	paths.Set("/api/v1/account-secrets/{id}", oapiAccountSecretByIDPath())
	paths.Set("/api/v1/account-secrets/{id}/disable", oapiAccountSecretDisablePath())
	paths.Set("/api/v1/account-secrets/{id}/enable", oapiAccountSecretEnablePath())

	return paths
}

// ---- Health ----

func oapiHealthPath() *openapi3.PathItem {
	return &openapi3.PathItem{
		Get: &openapi3.Operation{
			Tags:        []string{"Health"},
			Summary:     "Health check",
			OperationID: "health",
			Responses: openapi3.NewResponses(
				openapi3.WithStatus(200, oapiOKResponse("OK", &openapi3.SchemaRef{
					Value: openapi3.NewObjectSchema().WithProperty("status", oapiStrProp("")),
				})),
			),
		},
	}
}

// ---- Auth public ----

func oapiAuthConfigPath() *openapi3.PathItem {
	return &openapi3.PathItem{
		Get: &openapi3.Operation{
			Tags:        []string{"Auth"},
			Summary:     "Get authentication configuration",
			OperationID: "getAuthConfig",
			Responses: openapi3.NewResponses(
				openapi3.WithStatus(200, oapiOKResponse("OK", &openapi3.SchemaRef{
					Value: openapi3.NewObjectSchema().
						WithProperty("oauth_providers", oapiArrayOf(&openapi3.SchemaRef{Value: openapi3.NewStringSchema()})).
						WithProperty("email_enabled", oapiBoolProp("")).
						WithProperty("oauth_enabled", oapiBoolProp("")).
						WithProperty("oauth", &openapi3.Schema{Type: &openapi3.Types{"object"}}),
				})),
			),
		},
	}
}

func oapiAuthSendCodePath() *openapi3.PathItem {
	reqSchema := openapi3.NewObjectSchema().
		WithProperty("email", oapiStrProp("Email address")).
		WithProperty("type", oapiStrProp("Code type: login, register, bind_email, reset_password"))
	reqSchema.Required = []string{"email", "type"}

	return &openapi3.PathItem{
		Post: &openapi3.Operation{
			Tags:        []string{"Auth"},
			Summary:     "Send verification code to email",
			OperationID: "sendCode",
			RequestBody: oapiJSONBody(&openapi3.SchemaRef{Value: reqSchema}, true),
			Responses: openapi3.NewResponses(
				openapi3.WithStatus(200, oapiOKResponse("OK", &openapi3.SchemaRef{
					Value: openapi3.NewObjectSchema().
						WithProperty("expires_in", oapiIntProp("")).
						WithProperty("message", oapiStrProp("")),
				})),
				openapi3.WithStatus(400, oapiErrResponse("Bad Request")),
			),
		},
	}
}

func oapiAuthLoginPath() *openapi3.PathItem {
	reqSchema := openapi3.NewObjectSchema().
		WithProperty("email", oapiStrProp("")).
		WithProperty("code", oapiStrProp("6-digit verification code"))
	reqSchema.Required = []string{"email", "code"}

	return &openapi3.PathItem{
		Post: &openapi3.Operation{
			Tags:        []string{"Auth"},
			Summary:     "Login with email and verification code",
			OperationID: "login",
			RequestBody: oapiJSONBody(&openapi3.SchemaRef{Value: reqSchema}, true),
			Responses: openapi3.NewResponses(
				openapi3.WithStatus(200, oapiOKResponse("OK", &openapi3.SchemaRef{
					Value: openapi3.NewObjectSchema().
						WithPropertyRef("user", oapiRef("User")).
						WithProperty("token", oapiStrProp("")).
						WithProperty("refresh_token", oapiStrProp("")).
						WithProperty("expires_in", oapiInt64Prop("")).
						WithProperty("linked_accounts", oapiArrayOf(oapiRef("LinkedAccount"))),
				})),
				openapi3.WithStatus(400, oapiErrResponse("Bad Request")),
			),
		},
	}
}

func oapiAuthRegisterPath() *openapi3.PathItem {
	reqSchema := openapi3.NewObjectSchema().
		WithProperty("first_name", oapiStrProp("")).
		WithProperty("last_name", oapiStrProp("")).
		WithProperty("email", oapiStrProp("")).
		WithProperty("code", oapiStrProp("6-digit verification code"))
	reqSchema.Required = []string{"first_name", "last_name", "email", "code"}

	return &openapi3.PathItem{
		Post: &openapi3.Operation{
			Tags:        []string{"Auth"},
			Summary:     "Register with email and verification code",
			OperationID: "register",
			RequestBody: oapiJSONBody(&openapi3.SchemaRef{Value: reqSchema}, true),
			Responses: openapi3.NewResponses(
				openapi3.WithStatus(200, oapiOKResponse("OK", &openapi3.SchemaRef{
					Value: openapi3.NewObjectSchema().
						WithPropertyRef("user", oapiRef("User")).
						WithProperty("token", oapiStrProp("")).
						WithProperty("expires_in", oapiInt64Prop("")).
						WithProperty("is_new_user", oapiBoolProp("")),
				})),
				openapi3.WithStatus(400, oapiErrResponse("Bad Request")),
			),
		},
	}
}

func oapiAuthOAuthAuthorizePath() *openapi3.PathItem {
	return &openapi3.PathItem{
		Get: &openapi3.Operation{
			Tags:        []string{"Auth"},
			Summary:     "Get OAuth authorization URL",
			OperationID: "oauthAuthorize",
			Parameters: openapi3.Parameters{
				oapiQueryParamRequired("provider", "OAuth provider", openapi3.NewStringSchema()),
				oapiQueryParam("redirect_uri", "Redirect URI", openapi3.NewStringSchema()),
				oapiQueryParam("state", "State parameter", openapi3.NewStringSchema()),
			},
			Responses: openapi3.NewResponses(
				openapi3.WithStatus(200, oapiOKResponse("OK", &openapi3.SchemaRef{
					Value: openapi3.NewObjectSchema().
						WithProperty("authorize_url", oapiStrProp("")).
						WithProperty("state", oapiStrProp("")),
				})),
			),
		},
	}
}

func oapiAuthOAuthCallbackPath() *openapi3.PathItem {
	reqSchema := openapi3.NewObjectSchema().
		WithProperty("provider", oapiStrProp("")).
		WithProperty("code", oapiStrProp("")).
		WithProperty("state", oapiStrProp(""))
	reqSchema.Required = []string{"provider", "code", "state"}

	return &openapi3.PathItem{
		Post: &openapi3.Operation{
			Tags:        []string{"Auth"},
			Summary:     "Handle OAuth callback",
			OperationID: "oauthCallback",
			RequestBody: oapiJSONBody(&openapi3.SchemaRef{Value: reqSchema}, true),
			Responses: openapi3.NewResponses(
				openapi3.WithStatus(200, oapiOKResponse("OK", nil)),
			),
		},
	}
}

func oapiAuthRefreshPath() *openapi3.PathItem {
	reqSchema := openapi3.NewObjectSchema().
		WithProperty("refresh_token", oapiStrProp(""))
	reqSchema.Required = []string{"refresh_token"}

	return &openapi3.PathItem{
		Post: &openapi3.Operation{
			Tags:        []string{"Auth"},
			Summary:     "Refresh JWT token",
			OperationID: "refreshToken",
			RequestBody: oapiJSONBody(&openapi3.SchemaRef{Value: reqSchema}, true),
			Responses: openapi3.NewResponses(
				openapi3.WithStatus(200, oapiOKResponse("OK", &openapi3.SchemaRef{
					Value: openapi3.NewObjectSchema().
						WithProperty("token", oapiStrProp("")).
						WithProperty("refresh_token", oapiStrProp("")).
						WithProperty("expires_in", oapiInt64Prop("")),
				})),
			),
		},
	}
}

// ---- Auth protected ----

func oapiAuthLogoutPath() *openapi3.PathItem {
	return &openapi3.PathItem{
		Post: &openapi3.Operation{
			Tags:        []string{"Auth"},
			Summary:     "Logout",
			OperationID: "logout",
			Security:    oapiBearerSecurity(),
			Responses: openapi3.NewResponses(
				openapi3.WithStatus(200, oapiOKResponse("OK", oapiRef("MessageResponse"))),
			),
		},
	}
}

func oapiAuthMePath() *openapi3.PathItem {
	return &openapi3.PathItem{
		Get: &openapi3.Operation{
			Tags:        []string{"Auth"},
			Summary:     "Get current user",
			OperationID: "getCurrentUser",
			Security:    oapiBearerSecurity(),
			Responses: openapi3.NewResponses(
				openapi3.WithStatus(200, oapiOKResponse("OK", &openapi3.SchemaRef{
					Value: openapi3.NewObjectSchema().
						WithPropertyRef("user", oapiRef("User")).
						WithProperty("linked_accounts", oapiArrayOf(oapiRef("LinkedAccount"))),
				})),
			),
		},
		Put: &openapi3.Operation{
			Tags:        []string{"Auth"},
			Summary:     "Update profile",
			OperationID: "updateProfile",
			Security:    oapiBearerSecurity(),
			RequestBody: oapiJSONBody(&openapi3.SchemaRef{
				Value: openapi3.NewObjectSchema().
					WithProperty("first_name", oapiStrProp("")).
					WithProperty("last_name", oapiStrProp("")).
					WithProperty("display_name", oapiStrProp("")).
					WithProperty("language", oapiStrProp("")).
					WithProperty("timezone", oapiStrProp("")),
			}, true),
			Responses: openapi3.NewResponses(
				openapi3.WithStatus(200, oapiOKResponse("OK", &openapi3.SchemaRef{
					Value: openapi3.NewObjectSchema().WithPropertyRef("user", oapiRef("User")),
				})),
			),
		},
	}
}

func oapiAuthAccountsPath() *openapi3.PathItem {
	return &openapi3.PathItem{
		Get: &openapi3.Operation{
			Tags:        []string{"Auth"},
			Summary:     "Get linked accounts",
			OperationID: "getLinkedAccounts",
			Security:    oapiBearerSecurity(),
			Responses: openapi3.NewResponses(
				openapi3.WithStatus(200, oapiOKResponse("OK", &openapi3.SchemaRef{
					Value: openapi3.NewObjectSchema().
						WithProperty("accounts", oapiArrayOf(oapiRef("LinkedAccount"))),
				})),
			),
		},
	}
}

func oapiAuthAccountsPrimaryPath() *openapi3.PathItem {
	reqSchema := openapi3.NewObjectSchema().
		WithProperty("provider", oapiStrProp("")).
		WithProperty("provider_user_id", oapiStrProp(""))
	reqSchema.Required = []string{"provider"}

	return &openapi3.PathItem{
		Put: &openapi3.Operation{
			Tags:        []string{"Auth"},
			Summary:     "Set primary account",
			OperationID: "setPrimaryAccount",
			Security:    oapiBearerSecurity(),
			RequestBody: oapiJSONBody(&openapi3.SchemaRef{Value: reqSchema}, true),
			Responses: openapi3.NewResponses(
				openapi3.WithStatus(200, oapiOKResponse("OK", oapiRef("MessageResponse"))),
			),
		},
	}
}

func oapiAuthOAuthLinkPath() *openapi3.PathItem {
	reqSchema := openapi3.NewObjectSchema().
		WithProperty("provider", oapiStrProp("")).
		WithProperty("code", oapiStrProp(""))
	reqSchema.Required = []string{"provider", "code"}

	return &openapi3.PathItem{
		Post: &openapi3.Operation{
			Tags:        []string{"Auth"},
			Summary:     "Link OAuth account",
			OperationID: "linkOAuth",
			Security:    oapiBearerSecurity(),
			RequestBody: oapiJSONBody(&openapi3.SchemaRef{Value: reqSchema}, true),
			Responses: openapi3.NewResponses(
				openapi3.WithStatus(200, oapiOKResponse("OK", oapiRef("MessageResponse"))),
			),
		},
	}
}

func oapiAuthOAuthUnlinkPath() *openapi3.PathItem {
	reqSchema := openapi3.NewObjectSchema().
		WithProperty("provider", oapiStrProp("")).
		WithProperty("provider_user_id", oapiStrProp(""))
	reqSchema.Required = []string{"provider"}

	return &openapi3.PathItem{
		Delete: &openapi3.Operation{
			Tags:        []string{"Auth"},
			Summary:     "Unlink OAuth account",
			OperationID: "unlinkOAuth",
			Security:    oapiBearerSecurity(),
			RequestBody: oapiJSONBody(&openapi3.SchemaRef{Value: reqSchema}, true),
			Responses: openapi3.NewResponses(
				openapi3.WithStatus(200, oapiOKResponse("OK", oapiRef("MessageResponse"))),
			),
		},
	}
}

func oapiAuthBindEmailPath() *openapi3.PathItem {
	reqSchema := openapi3.NewObjectSchema().
		WithProperty("email", oapiStrProp("")).
		WithProperty("code", oapiStrProp("6-digit verification code"))
	reqSchema.Required = []string{"email", "code"}

	return &openapi3.PathItem{
		Post: &openapi3.Operation{
			Tags:        []string{"Auth"},
			Summary:     "Bind email to account",
			OperationID: "bindEmail",
			Security:    oapiBearerSecurity(),
			RequestBody: oapiJSONBody(&openapi3.SchemaRef{Value: reqSchema}, true),
			Responses: openapi3.NewResponses(
				openapi3.WithStatus(200, oapiOKResponse("OK", &openapi3.SchemaRef{
					Value: openapi3.NewObjectSchema().WithPropertyRef("user", oapiRef("User")),
				})),
			),
		},
	}
}

// ---- Books ----

func oapiBooksPath() *openapi3.PathItem {
	createReq := openapi3.NewObjectSchema().
		WithProperty("title", oapiStrProp("")).
		WithProperty("author", oapiStrProp("")).
		WithProperty("isbn", oapiStrProp("")).
		WithProperty("publisher", oapiStrProp("")).
		WithProperty("year", oapiIntProp(""))
	createReq.Required = []string{"title"}

	return &openapi3.PathItem{
		Get: &openapi3.Operation{
			Tags:        []string{"Books"},
			Summary:     "List books",
			OperationID: "listBooks",
			Security:    oapiBearerSecurity(),
			Parameters: openapi3.Parameters{
				oapiQueryParam("page", "Page number", openapi3.NewIntegerSchema()),
				oapiQueryParam("page_size", "Page size", openapi3.NewIntegerSchema()),
				oapiQueryParam("status", "Filter by status", openapi3.NewStringSchema()),
				oapiQueryParam("search", "Search keyword", openapi3.NewStringSchema()),
			},
			Responses: openapi3.NewResponses(
				openapi3.WithStatus(200, oapiOKResponse("OK", &openapi3.SchemaRef{
					Value: openapi3.NewObjectSchema().
						WithProperty("data", oapiArrayOf(oapiRef("Book"))).
						WithProperty("total", oapiInt64Prop("")).
						WithProperty("message", oapiStrProp("")),
				})),
			),
		},
		Post: &openapi3.Operation{
			Tags:        []string{"Books"},
			Summary:     "Create a book",
			OperationID: "createBook",
			Security:    oapiBearerSecurity(),
			RequestBody: oapiJSONBody(&openapi3.SchemaRef{Value: createReq}, true),
			Responses: openapi3.NewResponses(
				openapi3.WithStatus(200, oapiOKResponse("OK", &openapi3.SchemaRef{
					Value: openapi3.NewObjectSchema().
						WithPropertyRef("data", oapiRef("Book")).
						WithProperty("message", oapiStrProp("")),
				})),
			),
		},
	}
}

func oapiBookByIDPath() *openapi3.PathItem {
	updateReq := openapi3.NewObjectSchema().
		WithProperty("title", oapiStrProp("")).
		WithProperty("author", oapiStrProp("")).
		WithProperty("isbn", oapiStrProp("")).
		WithProperty("publisher", oapiStrProp("")).
		WithProperty("year", oapiIntProp(""))

	idParam := openapi3.Parameters{oapiPathParam("id", "Book ID")}

	return &openapi3.PathItem{
		Get: &openapi3.Operation{
			Tags:        []string{"Books"},
			Summary:     "Get a book by ID",
			OperationID: "getBook",
			Security:    oapiBearerSecurity(),
			Parameters:  idParam,
			Responses: openapi3.NewResponses(
				openapi3.WithStatus(200, oapiOKResponse("OK", &openapi3.SchemaRef{
					Value: openapi3.NewObjectSchema().
						WithPropertyRef("data", oapiRef("Book")).
						WithProperty("message", oapiStrProp("")),
				})),
				openapi3.WithStatus(404, oapiErrResponse("Not Found")),
			),
		},
		Put: &openapi3.Operation{
			Tags:        []string{"Books"},
			Summary:     "Update a book",
			OperationID: "updateBook",
			Security:    oapiBearerSecurity(),
			Parameters:  idParam,
			RequestBody: oapiJSONBody(&openapi3.SchemaRef{Value: updateReq}, true),
			Responses: openapi3.NewResponses(
				openapi3.WithStatus(200, oapiOKResponse("OK", &openapi3.SchemaRef{
					Value: openapi3.NewObjectSchema().
						WithPropertyRef("data", oapiRef("Book")).
						WithProperty("message", oapiStrProp("")),
				})),
			),
		},
		Delete: &openapi3.Operation{
			Tags:        []string{"Books"},
			Summary:     "Delete a book",
			OperationID: "deleteBook",
			Security:    oapiBearerSecurity(),
			Parameters:  idParam,
			Responses: openapi3.NewResponses(
				openapi3.WithStatus(200, oapiOKResponse("OK", oapiRef("MessageResponse"))),
			),
		},
	}
}

// ---- Images ----

func oapiImagesPath() *openapi3.PathItem {
	return &openapi3.PathItem{
		Get: &openapi3.Operation{
			Tags:        []string{"Images"},
			Summary:     "List images",
			OperationID: "listImages",
			Security:    oapiBearerSecurity(),
			Parameters: openapi3.Parameters{
				oapiQueryParam("page", "Page number", openapi3.NewIntegerSchema()),
				oapiQueryParam("page_size", "Page size", openapi3.NewIntegerSchema()),
			},
			Responses: openapi3.NewResponses(
				openapi3.WithStatus(200, oapiOKResponse("OK", &openapi3.SchemaRef{
					Value: openapi3.NewObjectSchema().
						WithProperty("data", openapi3.NewObjectSchema().
							WithProperty("items", oapiArrayOf(oapiRef("Image"))).
							WithProperty("total", oapiInt64Prop(""))).
						WithProperty("message", oapiStrProp("")),
				})),
			),
		},
		Post: &openapi3.Operation{
			Tags:        []string{"Images"},
			Summary:     "Upload an image",
			OperationID: "uploadImage",
			Security:    oapiBearerSecurity(),
			RequestBody: &openapi3.RequestBodyRef{
				Value: &openapi3.RequestBody{
					Required: true,
					Content: openapi3.Content{
						"multipart/form-data": &openapi3.MediaType{
							Schema: &openapi3.SchemaRef{
								Value: openapi3.NewObjectSchema().
									WithProperty("file", &openapi3.Schema{Type: &openapi3.Types{"string"}, Format: "binary"}).
									WithProperty("group_id", oapiIntProp("")),
							},
						},
					},
				},
			},
			Responses: openapi3.NewResponses(
				openapi3.WithStatus(200, oapiOKResponse("OK", &openapi3.SchemaRef{
					Value: openapi3.NewObjectSchema().
						WithPropertyRef("data", oapiRef("Image")).
						WithProperty("message", oapiStrProp("")),
				})),
			),
		},
	}
}

func oapiImageByIDPath() *openapi3.PathItem {
	idParam := openapi3.Parameters{oapiPathParam("id", "Image ID")}

	return &openapi3.PathItem{
		Get: &openapi3.Operation{
			Tags:        []string{"Images"},
			Summary:     "Get image metadata",
			OperationID: "getImage",
			Security:    oapiBearerSecurity(),
			Parameters:  idParam,
			Responses: openapi3.NewResponses(
				openapi3.WithStatus(200, oapiOKResponse("OK", &openapi3.SchemaRef{
					Value: openapi3.NewObjectSchema().
						WithPropertyRef("data", oapiRef("Image")).
						WithProperty("message", oapiStrProp("")),
				})),
			),
		},
		Delete: &openapi3.Operation{
			Tags:        []string{"Images"},
			Summary:     "Delete an image",
			OperationID: "deleteImage",
			Security:    oapiBearerSecurity(),
			Parameters:  idParam,
			Responses: openapi3.NewResponses(
				openapi3.WithStatus(200, oapiOKResponse("OK", oapiRef("MessageResponse"))),
			),
		},
	}
}

func oapiImageFilePath() *openapi3.PathItem {
	return &openapi3.PathItem{
		Get: &openapi3.Operation{
			Tags:        []string{"Images"},
			Summary:     "Serve image file",
			OperationID: "serveImage",
			Security:    oapiBearerSecurity(),
			Parameters: openapi3.Parameters{
				oapiPathParam("id", "Image ID"),
				oapiQueryParam("variant", "Image variant (e.g. thumbnail)", openapi3.NewStringSchema()),
			},
			Responses: openapi3.NewResponses(
				openapi3.WithStatus(200, &openapi3.ResponseRef{
					Value: openapi3.NewResponse().WithDescription("Image binary"),
				}),
			),
		},
	}
}

// ---- Image Groups ----

func oapiImageGroupsPath() *openapi3.PathItem {
	reqSchema := openapi3.NewObjectSchema().
		WithProperty("group_type", oapiStrProp("book, album, or custom")).
		WithProperty("group_name", oapiStrProp("")).
		WithProperty("ref_id", oapiIntProp("")).
		WithProperty("ref_type", oapiStrProp(""))
	reqSchema.Required = []string{"group_type", "group_name"}

	return &openapi3.PathItem{
		Post: &openapi3.Operation{
			Tags:        []string{"ImageGroups"},
			Summary:     "Create an image group",
			OperationID: "createGroup",
			Security:    oapiBearerSecurity(),
			RequestBody: oapiJSONBody(&openapi3.SchemaRef{Value: reqSchema}, true),
			Responses: openapi3.NewResponses(
				openapi3.WithStatus(200, oapiOKResponse("OK", &openapi3.SchemaRef{
					Value: openapi3.NewObjectSchema().
						WithPropertyRef("data", oapiRef("ImageGroup")).
						WithProperty("message", oapiStrProp("")),
				})),
			),
		},
	}
}

func oapiImageGroupImagesPath() *openapi3.PathItem {
	reqSchema := openapi3.NewObjectSchema().
		WithProperty("image_ids", oapiArrayOf(&openapi3.SchemaRef{Value: openapi3.NewIntegerSchema()}))
	reqSchema.Required = []string{"image_ids"}

	return &openapi3.PathItem{
		Post: &openapi3.Operation{
			Tags:        []string{"ImageGroups"},
			Summary:     "Add images to a group",
			OperationID: "addImagesToGroup",
			Security:    oapiBearerSecurity(),
			Parameters:  openapi3.Parameters{oapiPathParam("id", "Group ID")},
			RequestBody: oapiJSONBody(&openapi3.SchemaRef{Value: reqSchema}, true),
			Responses: openapi3.NewResponses(
				openapi3.WithStatus(200, oapiOKResponse("OK", oapiRef("MessageResponse"))),
			),
		},
	}
}

// ---- Account Secrets ----

func oapiAccountSecretsPath() *openapi3.PathItem {
	createReq := openapi3.NewObjectSchema().
		WithProperty("name", oapiStrProp("")).
		WithProperty("scope", oapiStrProp("")).
		WithProperty("resource_id", oapiIntProp("")).
		WithProperty("resource_type", oapiStrProp("")).
		WithProperty("permissions", oapiArrayOf(&openapi3.SchemaRef{Value: openapi3.NewStringSchema()})).
		WithProperty("expires_in_days", oapiIntProp(""))
	createReq.Required = []string{"name", "scope"}

	return &openapi3.PathItem{
		Get: &openapi3.Operation{
			Tags:        []string{"AccountSecrets"},
			Summary:     "List account secrets",
			OperationID: "listAccountSecrets",
			Security:    oapiBearerSecurity(),
			Parameters: openapi3.Parameters{
				oapiQueryParam("page", "Page number", openapi3.NewIntegerSchema()),
				oapiQueryParam("page_size", "Page size", openapi3.NewIntegerSchema()),
			},
			Responses: openapi3.NewResponses(
				openapi3.WithStatus(200, oapiOKResponse("OK", &openapi3.SchemaRef{
					Value: openapi3.NewObjectSchema().
						WithProperty("data", openapi3.NewObjectSchema().
							WithProperty("items", oapiArrayOf(oapiRef("AccountSecret"))).
							WithProperty("total", oapiInt64Prop("")).
							WithProperty("page", oapiIntProp("")).
							WithProperty("page_size", oapiIntProp(""))).
						WithProperty("total", oapiInt64Prop("")).
						WithProperty("page", oapiIntProp("")).
						WithProperty("page_size", oapiIntProp("")).
						WithProperty("message", oapiStrProp("")),
				})),
			),
		},
		Post: &openapi3.Operation{
			Tags:        []string{"AccountSecrets"},
			Summary:     "Create an account secret",
			OperationID: "createAccountSecret",
			Security:    oapiBearerSecurity(),
			RequestBody: oapiJSONBody(&openapi3.SchemaRef{Value: createReq}, true),
			Responses: openapi3.NewResponses(
				openapi3.WithStatus(200, oapiOKResponse("OK", &openapi3.SchemaRef{
					Value: openapi3.NewObjectSchema().
						WithProperty("data", &openapi3.Schema{Type: &openapi3.Types{"object"}, Description: "Secret details including the secret key (shown only once)"}).
						WithProperty("message", oapiStrProp("")),
				})),
			),
		},
	}
}

func oapiAccountSecretByIDPath() *openapi3.PathItem {
	idParam := openapi3.Parameters{oapiPathParam("id", "Account secret ID")}

	return &openapi3.PathItem{
		Get: &openapi3.Operation{
			Tags:        []string{"AccountSecrets"},
			Summary:     "Get account secret detail",
			OperationID: "getAccountSecret",
			Security:    oapiBearerSecurity(),
			Parameters:  idParam,
			Responses: openapi3.NewResponses(
				openapi3.WithStatus(200, oapiOKResponse("OK", &openapi3.SchemaRef{
					Value: openapi3.NewObjectSchema().
						WithPropertyRef("data", oapiRef("AccountSecret")).
						WithProperty("message", oapiStrProp("")),
				})),
			),
		},
		Delete: &openapi3.Operation{
			Tags:        []string{"AccountSecrets"},
			Summary:     "Delete an account secret",
			OperationID: "deleteAccountSecret",
			Security:    oapiBearerSecurity(),
			Parameters:  idParam,
			Responses: openapi3.NewResponses(
				openapi3.WithStatus(200, oapiOKResponse("OK", oapiRef("MessageResponse"))),
			),
		},
	}
}

func oapiAccountSecretDisablePath() *openapi3.PathItem {
	return &openapi3.PathItem{
		Post: &openapi3.Operation{
			Tags:        []string{"AccountSecrets"},
			Summary:     "Disable an account secret",
			OperationID: "disableAccountSecret",
			Security:    oapiBearerSecurity(),
			Parameters:  openapi3.Parameters{oapiPathParam("id", "Account secret ID")},
			Responses: openapi3.NewResponses(
				openapi3.WithStatus(200, oapiOKResponse("OK", oapiRef("MessageResponse"))),
			),
		},
	}
}

func oapiAccountSecretEnablePath() *openapi3.PathItem {
	return &openapi3.PathItem{
		Post: &openapi3.Operation{
			Tags:        []string{"AccountSecrets"},
			Summary:     "Enable an account secret",
			OperationID: "enableAccountSecret",
			Security:    oapiBearerSecurity(),
			Parameters:  openapi3.Parameters{oapiPathParam("id", "Account secret ID")},
			Responses: openapi3.NewResponses(
				openapi3.WithStatus(200, oapiOKResponse("OK", oapiRef("MessageResponse"))),
			),
		},
	}
}
