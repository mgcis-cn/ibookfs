# iBookFS Authentication API Design

## Overview

This document describes the REST API for multi-provider authentication in iBookFS.

**Base URL**: `/api/v1/auth`

**Supported Providers**: Email (verification code), GitHub, Gitee, iCloud

---

## Authentication Flow

```
┌─────────────────────────────────────────────────────────────────────────┐
│                         AUTHENTICATION FLOWS                            │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  ┌─────────────┐     ┌──────────────┐     ┌─────────────────────┐    │
│  │ Email Login │ →   │ Send Code    │ →   │ Verify & Login      │    │
│  └─────────────┘     └──────────────┘     └─────────────────────┘    │
│                                                                         │
│  ┌─────────────┐     ┌──────────────┐     ┌─────────────────────┐    │
│  │ OAuth Login │ →   │ Provider     │ →   │ Link or Create      │    │
│  │             │     │ Redirect     │     │ Account             │    │
│  └─────────────┘     └──────────────┘     └─────────────────────┘    │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

---

## API Endpoints

### 1. Send Verification Code

Initiates email verification by sending a 6-digit code.

**Endpoint**: `POST /send-code`

**Request Body**:
```json
{
  "email": "user@example.com",
  "type": "login" | "register" | "bind_email" | "reset_password"
}
```

**Response**:
```json
{
  "success": true,
  "data": {
    "expires_in": 300,
    "message": "验证码已发送"
  }
}
```

**Error Responses**:
| Status | Code | Description |
|--------|------|-------------|
| 429 | RATE_LIMIT_EXCEEDED | Too many requests, wait 60s |
| 400 | INVALID_EMAIL | Invalid email format |
| 400 | EMAIL_ALREADY_EXISTS | Email already registered (for register) |

---

### 2. Login with Email Code

Authenticates user using email and verification code.

**Endpoint**: `POST /login`

**Request Body**:
```json
{
  "email": "user@example.com",
  "code": "123456"
}
```

**Response**:
```json
{
  "success": true,
  "data": {
    "user": {
      "id": "123",
      "email": "user@example.com",
      "first_name": "张",
      "last_name": "三",
      "display_name": "张三",
      "avatar_url": null,
      "created_at": "2024-01-01T00:00:00Z"
    },
    "token": "eyJhbGciOiJIUzI1NiIs...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
    "expires_in": 86400,
    "linked_accounts": [
      {
        "provider": "email",
        "provider_email": "user@example.com",
        "is_primary": true
      }
    ]
  }
}
```

**Error Responses**:
| Status | Code | Description |
|--------|------|-------------|
| 401 | INVALID_CODE | Verification code invalid |
| 401 | CODE_EXPIRED | Verification code expired |
| 404 | USER_NOT_FOUND | User not found (use register) |

---

### 3. Register with Email

Creates a new account using email verification.

**Endpoint**: `POST /register`

**Request Body**:
```json
{
  "email": "user@example.com",
  "code": "123456",
  "first_name": "张",
  "last_name": "三"
}
```

**Response**:
```json
{
  "success": true,
  "data": {
    "user": {
      "id": "123",
      "email": "user@example.com",
      "first_name": "张",
      "last_name": "三",
      "created_at": "2024-01-01T00:00:00Z"
    },
    "token": "eyJhbGciOiJIUzI1NiIs...",
    "expires_in": 86400
  }
}
```

**Error Responses**:
| Status | Code | Description |
|--------|------|-------------|
| 400 | EMAIL_EXISTS | Email already registered |
| 401 | INVALID_CODE | Verification code invalid |

---

### 4. OAuth Authorization URLs

Get authorization URLs for OAuth providers.

**Endpoint**: `GET /oauth/authorize`

**Query Parameters**:
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| provider | string | Yes | `github`, `gitee`, `icloud` |
| redirect_uri | string | Yes | URL to redirect after auth |
| state | string | No | CSRF protection token |

**Response**:
```json
{
  "success": true,
  "data": {
    "authorize_url": "https://github.com/login/oauth/authorize?client_id=...",
    "state": "random_state_string"
  }
}
```

---

### 5. OAuth Callback

Handle OAuth callback from provider.

**Endpoint**: `POST /oauth/callback`

**Request Body**:
```json
{
  "provider": "github",
  "code": "authorization_code_from_provider",
  "state": "state_from_authorize"
}
```

**Response** (New Account):
```json
{
  "success": true,
  "data": {
    "user": {
      "id": "123",
      "email": "user@example.com",
      "first_name": "张",
      "last_name": "三",
      "display_name": "张三",
      "avatar_url": "https://avatars.githubusercontent.com/u/123",
      "created_at": "2024-01-01T00:00:00Z"
    },
    "token": "eyJhbGciOiJIUzI1NiIs...",
    "is_new_user": true
  }
}
```

**Response** (Existing Account):
```json
{
  "success": true,
  "data": {
    "user": { ... },
    "token": "eyJhbGciOiJIUzI1NiIs...",
    "is_new_user": false
  }
}
```

**OAuth Flow States**:

| Scenario | Action |
|----------|--------|
| New user, no email match | Create account, link OAuth |
| Existing user, OAuth already linked | Login, return token |
| Existing user, OAuth not linked, emails match | Link OAuth to account, login |
| Existing user, OAuth not linked, emails differ | Ask user to confirm linking |

---

### 6. Link OAuth Account

Link an OAuth provider to an authenticated user's account.

**Endpoint**: `POST /oauth/link`

**Headers**: `Authorization: Bearer {token}`

**Request Body**:
```json
{
  "provider": "github",
  "code": "authorization_code_from_provider"
}
```

**Response**:
```json
{
  "success": true,
  "data": {
    "linked_account": {
      "provider": "github",
      "provider_username": "username",
      "provider_email": "user@example.com",
      "is_primary": false,
      "linked_at": "2024-01-01T00:00:00Z"
    }
  }
}
```

**Error Responses**:
| Status | Code | Description |
|--------|------|-------------|
| 400 | ALREADY_LINKED | OAuth already linked to this account |
| 400 | LINKED_TO_OTHER | OAuth already linked to another account |
| 401 | UNAUTHORIZED | Invalid or expired token |

---

### 7. Unlink OAuth Account

Remove an OAuth provider from user's account.

**Endpoint**: `DELETE /oauth/unlink`

**Headers**: `Authorization: Bearer {token}`

**Request Body**:
```json
{
  "provider": "github"
}
```

**Response**:
```json
{
  "success": true,
  "data": {
    "message": "OAuth account unlinked successfully"
  }
}
```

**Error Responses**:
| Status | Code | Description |
|--------|------|-------------|
| 400 | CANNOT_UNLINK_PRIMARY | Cannot unlink primary login method |
| 404 | NOT_LINKED | OAuth not linked to this account |

---

### 8. Bind Email to Account

Add email as a login method for OAuth-only users.

**Endpoint**: `POST /bind-email`

**Headers**: `Authorization: Bearer {token}`

**Request Body**:
```json
{
  "email": "user@example.com",
  "code": "123456"
}
```

**Response**:
```json
{
  "success": true,
  "data": {
    "user": {
      "id": "123",
      "email": "user@example.com",
      "email_verified": true,
      "first_name": "张",
      "last_name": "三"
    }
  }
}
```

---

### 9. Get Linked Accounts

Get all linked authentication methods for current user.

**Endpoint**: `GET /accounts`

**Headers**: `Authorization: Bearer {token}`

**Response**:
```json
{
  "success": true,
  "data": {
    "accounts": [
      {
        "provider": "email",
        "provider_email": "user@example.com",
        "is_primary": true,
        "linked_at": "2024-01-01T00:00:00Z"
      },
      {
        "provider": "github",
        "provider_username": "username",
        "provider_email": "user@example.com",
        "is_primary": false,
        "linked_at": "2024-01-02T00:00:00Z",
        "last_used_at": "2024-01-05T00:00:00Z"
      },
      {
        "provider": "gitee",
        "provider_username": "username",
        "is_primary": false,
        "linked_at": "2024-01-03T00:00:00Z"
      }
    ]
  }
}
```

---

### 10. Set Primary Account

Set a specific login method as primary.

**Endpoint**: `PUT /accounts/primary`

**Headers**: `Authorization: Bearer {token}`

**Request Body**:
```json
{
  "provider": "github",
  "provider_user_id": "12345678"
}
```

**Response**:
```json
{
  "success": true,
  "data": {
    "message": "Primary account updated"
  }
}
```

---

### 11. Refresh Token

Refresh access token using refresh token.

**Endpoint**: `POST /refresh`

**Request Body**:
```json
{
  "refresh_token": "eyJhbGciOiJIUzI1NiIs..."
}
```

**Response**:
```json
{
  "success": true,
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIs...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
    "expires_in": 86400
  }
}
```

---

### 12. Logout

Invalidate current session.

**Endpoint**: `POST /logout`

**Headers**: `Authorization: Bearer {token}`

**Response**:
```json
{
  "success": true,
  "data": {
    "message": "Logged out successfully"
  }
}
```

---

### 13. Get Current User

Get authenticated user's profile.

**Endpoint**: `GET /me`

**Headers**: `Authorization: Bearer {token}`

**Response**:
```json
{
  "success": true,
  "data": {
    "id": "123",
    "email": "user@example.com",
    "email_verified": true,
    "first_name": "张",
    "last_name": "三",
    "display_name": "张三",
    "avatar_url": "https://...",
    "bio": null,
    "status": "active",
    "role": "user",
    "preferences": {
      "theme": "light"
    },
    "language": "zh-CN",
    "timezone": "Asia/Shanghai",
    "last_login_at": "2024-01-01T00:00:00Z",
    "created_at": "2024-01-01T00:00:00Z",
    "linked_accounts": [
      {
        "provider": "email",
        "is_primary": true
      },
      {
        "provider": "github",
        "provider_username": "username",
        "is_primary": false
      }
    ]
  }
}
```

---

### 14. Update Profile

Update user's profile information.

**Endpoint**: `PUT /me`

**Headers**: `Authorization: Bearer {token}`

**Request Body**:
```json
{
  "first_name": "李",
  "last_name": "四",
  "display_name": "李四",
  "bio": "热爱阅读",
  "language": "en-US",
  "timezone": "America/New_York"
}
```

**Response**:
```json
{
  "success": true,
  "data": {
    "user": { ... }
  }
}
```

---

## OAuth Provider Configuration

### GitHub

| Parameter | Value |
|-----------|-------|
| Authorization URL | `https://github.com/login/oauth/authorize` |
| Token URL | `https://github.com/login/oauth/access_token` |
| User API | `https://api.github.com/user` |
| Scopes | `read:user`, `user:email` |

### Gitee

| Parameter | Value |
|-----------|-------|
| Authorization URL | `https://gitee.com/oauth/authorize` |
| Token URL | `https://gitee.com/oauth/token` |
| User API | `https://gitee.com/api/v5/user` |
| Scopes | `user_info`, `emails` |

### iCloud

| Parameter | Value |
|-----------|-------|
| Authorization URL | `https://idmsa.apple.com/appleauth/auth` |
| Notes | Requires Apple Developer account, uses Sign in with Apple |

---

## Error Response Format

All error responses follow this format:

```json
{
  "success": false,
  "error": {
    "code": "ERROR_CODE",
    "message": "Human readable error message",
    "details": {}
  }
}
```

### Common Error Codes

| Code | Status | Description |
|------|--------|-------------|
| UNAUTHORIZED | 401 | Invalid or missing authentication |
| TOKEN_EXPIRED | 401 | JWT token has expired |
| INVALID_INPUT | 400 | Request validation failed |
| RATE_LIMITED | 429 | Too many requests |
| INTERNAL_ERROR | 500 | Server error |

---

## Security Considerations

1. **JWT Tokens**: Use RS256 for production, HS256 for development
2. **Token Storage**: HttpOnly, Secure, SameSite cookies for web
3. **Rate Limiting**:
   - Send code: 1 request per minute per email
   - Login attempts: 10 attempts per 15 minutes per IP
4. **CORS**: Configure allowed origins properly
5. **State Parameter**: Always use and verify state in OAuth flow
6. **Token Encryption**: Encrypt OAuth tokens in database

---

## Webhook Events (Future)

For integration with webhook system:

| Event | Description |
|-------|-------------|
| `user.created` | New user registered |
| `user.logged_in` | User logged in |
| `user.oauth_linked` | OAuth account linked |
| `user.oauth_unlinked` | OAuth account unlinked |
| `user.email_verified` | User email verified |
