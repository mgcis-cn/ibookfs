package gitee

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/mgcis-cn/ibookfs/pkg/authn/oauth/sources"
	"github.com/mgcis-cn/ibookfs/pkg/options"
)

const SourceKind = "gitee"

func init() {
	sources.Register(SourceKind, func(cfg sources.Config) (sources.OAuth, error) {
		opts, ok := cfg.(*options.GiteeOptions)
		if !ok {
			return nil, fmt.Errorf("gitee: invalid config type %T", cfg)
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
		return nil, fmt.Errorf("gitee: clientID and clientSecret are required")
	}
	if redirectURL == "" {
		return nil, fmt.Errorf("gitee: redirectURL is required in config file")
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
	params.Set("response_type", "code")
	params.Set("scope", "user_info")
	params.Set("state", state)
	return "https://gitee.com/oauth/authorize?" + params.Encode()
}

func (o *OAuth) Exchange(ctx context.Context, code string) (string, error) {
	params := url.Values{}
	params.Set("grant_type", "authorization_code")
	params.Set("client_id", o.ClientID)
	params.Set("client_secret", o.ClientSecret)
	params.Set("code", code)
	params.Set("redirect_uri", o.RedirectURL)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://gitee.com/oauth/token", strings.NewReader(params.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result struct {
		AccessToken string `json:"access_token"`
		Error       string `json:"error"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}
	if result.Error != "" {
		return "", fmt.Errorf("gitee: %s", result.Error)
	}
	return result.AccessToken, nil
}

func (o *OAuth) GetUserInfo(ctx context.Context, token string) (*sources.UserInfo, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://gitee.com/api/v5/user?access_token="+token, nil)
	if err != nil {
		return nil, err
	}
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

	return &sources.UserInfo{
		ID:       fmt.Sprintf("%d", user.ID),
		Name:     name,
		Email:    user.Email,
		Avatar:   user.AvatarURL,
		Provider: string(SourceKind),
	}, nil
}
