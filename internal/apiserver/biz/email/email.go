// Package email provides email business logic.
package email

import (
	"context"
	"fmt"

	apierr "github.com/mgcis-cn/ibookfs/internal/apiserver/errors"
	"github.com/mgcis-cn/ibookfs/pkg/email/sources"
)

// EmailBiz defines the interface for email operations.
type EmailBiz interface {
	SendVerificationCode(to, code string, codeType string) error
	IsEmailEnabled() bool
}

// emailBiz is the concrete implementation of EmailBiz.
type emailBiz struct {
	email   sources.Email
	baseURL string
}

// NewEmailBiz creates a new email business logic.
func NewEmailBiz(email sources.Email, baseURL string) EmailBiz {
	return &emailBiz{
		email:   email,
		baseURL: baseURL,
	}
}

// IsEmailEnabled checks if email is properly configured.
func (s *emailBiz) IsEmailEnabled() bool {
	return s.email != nil && s.email.IsEnabled()
}

// SendVerificationCode sends a verification code email.
func (s *emailBiz) SendVerificationCode(to, code string, codeType string) error {
	if !s.IsEmailEnabled() {
		return apierr.ErrEmailDisabled
	}
	subject := s.getSubject(codeType)
	body := s.getEmailBody(code, codeType)
	return s.email.Send(context.Background(), to, subject, body)
}

// getSubject returns email subject based on code type.
func (s *emailBiz) getSubject(codeType string) string {
	switch codeType {
	case "login":
		return "ibookfs-登录"
	case "register":
		return "ibookfs-注册"
	case "bind_email":
		return "ibookfs-绑定邮箱"
	case "reset_password":
		return "ibookfs-重置密码"
	default:
		return "ibookfs-邮箱验证"
	}
}

// getEmailBody returns email body with verification code.
func (s *emailBiz) getEmailBody(code, codeType string) string {
	action := "验证"
	switch codeType {
	case "login":
		action = "登录"
	case "register":
		action = "注册"
	case "bind_email":
		action = "绑定邮箱"
	case "reset_password":
		action = "重置密码"
	}
	logoURL := s.baseURL + "/assets/logo.svg"
	if s.baseURL == "" {
		logoURL = "/assets/logo.svg"
	}
	return fmt.Sprintf(`<!DOCTYPE html>
<html><head><meta charset="UTF-8"></head>
<body style="font-family:Arial,sans-serif;line-height:1.6;color:#333">
<div style="max-width:600px;margin:0 auto;padding:20px">
<div style="background:#f9f9f9;padding:30px;border-radius:10px">
<p>您好，您正在使用验证码进行%s操作。</p>
<p>您的验证码是：</p>
<div style="background:white;padding:20px;text-align:center;font-size:32px;letter-spacing:5px;font-weight:bold;margin:20px 0;border-radius:5px;border:2px dashed #667eea">%s</div>
<div style="text-align:center;margin:30px 0"><img src="%s" alt="iBookFS" style="width:100px;height:auto;opacity:0.85"/></div>
<p>验证码有效期为 <strong>5分钟</strong>，请尽快完成验证。</p>
<p>如果这不是您的操作，请忽略此邮件。</p>
</div>
<div style="text-align:center;margin-top:30px;color:#999;font-size:12px">
<p>此邮件由系统自动发送，请勿直接回复。</p>
<p>&copy; 2025 iBookFS. All rights reserved.</p>
</div></div></body></html>`, action, code, logoURL)
}
