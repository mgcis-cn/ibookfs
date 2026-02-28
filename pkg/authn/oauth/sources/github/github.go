package github

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/mgcis-cn/ibookfs/pkg/authn/oauth/sources"
	"github.com/mgcis-cn/ibookfs/pkg/options"
)

const SourceKind = "github"

func init() {
	sources.Register(SourceKind, func(cfg sources.Config) (sources.OAuth, error) {
		opts, ok := cfg.(*options.GithubOptions)
		if !ok {
			return nil, fmt.Errorf("github: invalid config type %T", cfg)
		}
		return New(opts.ClientID, opts.ClientSecret, opts.RedirectURL)
	})
}

type OAuth struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
}

func New(clientID, clientSecret, redirectURL string) (*OAuth, error) {
	if clientID == "" || clientSecret == "" {
		return nil, fmt.Errorf("github: clientID and clientSecret are required")
	}
	if redirectURL == "" {
		return nil, fmt.Errorf("github: redirectURL is required in config file")
	}
	return &OAuth{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
	}, nil
}

func (o *OAuth) Kind() string {
	return SourceKind
}

func (o *OAuth) GetAuthURL(state string) string {
	params := url.Values{}
	params.Set("client_id", o.ClientID)
	params.Set("redirect_uri", o.RedirectURL)
	params.Set("scope", "read:user user:email")
	params.Set("state", state)
	return "https://github.com/login/oauth/authorize?" + params.Encode()
}

func (o *OAuth) Exchange(ctx context.Context, code string) (string, error) {
	params := url.Values{}
	params.Set("client_id", o.ClientID)
	params.Set("client_secret", o.ClientSecret)
	params.Set("code", code)
	params.Set("redirect_uri", o.RedirectURL)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://github.com/login/oauth/access_token", nil)
	if err != nil {
		return "", err
	}
	req.URL.RawQuery = params.Encode()
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result struct {
		AccessToken      string `json:"access_token"`
		Error            string `json:"error"`
		ErrorDescription string `json:"error_description"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}
	if result.Error != "" {
		desc := result.ErrorDescription
		if desc == "" {
			desc = result.Error
		}
		return "", fmt.Errorf("github: %s", desc)
	}
	if result.AccessToken == "" {
		return "", fmt.Errorf("github: empty access token returned")
	}
	return result.AccessToken, nil
}

func (o *OAuth) GetUserInfo(ctx context.Context, token string) (*sources.UserInfo, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.github.com/user", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var user struct {
		ID        int64  `json:"id"`
		Login     string `json:"login"`
		Name      string `json:"name"`
		Email     string `json:"email"`
		AvatarURL string `json:"avatar_url"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, err
	}

	name := user.Name
	if name == "" {
		name = user.Login
	}

	// If email is not public, fetch from /user/emails API
	email := user.Email
	if email == "" {
		email = o.getPrimaryEmail(ctx, token)
	}

	return &sources.UserInfo{
		ID:       fmt.Sprintf("%d", user.ID),
		Name:     name,
		Email:    email,
		Avatar:   user.AvatarURL,
		Provider: string(SourceKind),
	}, nil
}

// getPrimaryEmail fetches the primary verified email from GitHub /user/emails API
func (o *OAuth) getPrimaryEmail(ctx context.Context, token string) string {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.github.com/user/emails", nil)
	if err != nil {
		return ""
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	var emails []struct {
		Email    string `json:"email"`
		Primary  bool   `json:"primary"`
		Verified bool   `json:"verified"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&emails); err != nil {
		return ""
	}

	// Find primary verified email
	for _, e := range emails {
		if e.Primary && e.Verified {
			return e.Email
		}
	}
	// Fallback to first verified email
	for _, e := range emails {
		if e.Verified {
			return e.Email
		}
	}
	return ""
}
