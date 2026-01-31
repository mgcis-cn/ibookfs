// Package service provides business logic.
package service

import (
	"crypto/tls"
	"fmt"
	"net/smtp"
	"strings"

	"github.com/mgcis-cn/ibookfs/internal/config"
)

// EmailService handles email operations.
type EmailService struct {
	cfg     *config.EmailConfig
	baseURL string // Base URL for static assets
}

// NewEmailService creates a new email service.
func NewEmailService(cfg *config.EmailConfig, baseURL string) *EmailService {
	return &EmailService{
		cfg:     cfg,
		baseURL: baseURL,
	}
}

// SendVerificationCode sends a verification code email.
func (s *EmailService) SendVerificationCode(to, code string, codeType string) error {
	if !s.cfg.IsEmailEnabled() {
		return fmt.Errorf("email is disabled")
	}

	subject := s.getSubject(codeType)
	body := s.getEmailBody(code, codeType)

	return s.sendEmail(to, subject, body)
}

// sendEmail sends an email using SMTP.
func (s *EmailService) sendEmail(to, subject, body string) error {
	cfg := s.cfg.GetEmailConfig()

	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)

	// Compose email message
	msg := s.buildEmailMessage(cfg.From, cfg.FromName, to, subject, body)

	// Send email
	switch cfg.Kind {
	case "netease", "smtp":
		return s.sendSMTP(addr, cfg.User, cfg.Password, cfg.From, to, msg)
	default:
		return fmt.Errorf("unsupported email kind: %s", cfg.Kind)
	}
}

// sendSMTP sends email via SMTP with SSL/TLS support.
func (s *EmailService) sendSMTP(addr, username, password, from, to string, msg []byte) error {
	cfg := s.cfg.GetEmailConfig()

	// For port 465 (SSL) or port 587 (STARTTLS)
	var auth smtp.Auth
	var client *smtp.Client
	var err error

	if cfg.Port == 465 {
		// SSL connection (port 465)
		tlsConfig := &tls.Config{
			ServerName: cfg.Host,
			MinVersion: tls.VersionTLS12,
		}

		conn, err := tls.Dial("tcp", addr, tlsConfig)
		if err != nil {
			return fmt.Errorf("failed to connect to SMTP server: %w", err)
		}

		client, err = smtp.NewClient(conn, cfg.Host)
		if err != nil {
			return fmt.Errorf("failed to create SMTP client: %w", err)
		}
	} else {
		// STARTTLS or plain connection (port 587 or 25)
		client, err = smtp.Dial(addr)
		if err != nil {
			return fmt.Errorf("failed to connect to SMTP server: %w", err)
		}

		// Start TLS if available
		if ok, _ := client.Extension("STARTTLS"); ok {
			if err := client.StartTLS(&tls.Config{
				ServerName: cfg.Host,
				MinVersion: tls.VersionTLS12,
			}); err != nil {
				return fmt.Errorf("failed to start TLS: %w", err)
			}
		}
	}
	defer client.Close()

	// Auth
	auth = smtp.PlainAuth("", username, password, cfg.Host)
	if err := client.Auth(auth); err != nil {
		return fmt.Errorf("failed to authenticate: %w", err)
	}

	// Set sender
	if err := client.Mail(from); err != nil {
		return fmt.Errorf("failed to set sender: %w", err)
	}

	// Set recipient
	if err := client.Rcpt(to); err != nil {
		return fmt.Errorf("failed to set recipient: %w", err)
	}

	// Send data
	writer, err := client.Data()
	if err != nil {
		return fmt.Errorf("failed to get data writer: %w", err)
	}
	defer writer.Close()

	_, err = writer.Write(msg)
	if err != nil {
		return fmt.Errorf("failed to write message: %w", err)
	}

	return nil
}

// buildEmailMessage builds a raw email message.
func (s *EmailService) buildEmailMessage(from, fromName, to, subject, body string) []byte {
	cfg := s.cfg.GetEmailConfig()
	if fromName == "" {
		fromName = cfg.FromName
	}

	var msg strings.Builder

	// Headers
	msg.WriteString(fmt.Sprintf("From: %s <%s>\r\n", fromName, from))
	msg.WriteString(fmt.Sprintf("To: %s\r\n", to))
	msg.WriteString(fmt.Sprintf("Subject: %s\r\n", subject))
	msg.WriteString("MIME-Version: 1.0\r\n")
	msg.WriteString("Content-Type: text/html; charset=UTF-8\r\n")
	msg.WriteString("\r\n")

	// Body
	msg.WriteString(body)

	return []byte(msg.String())
}

// getSubject returns email subject based on code type.
func (s *EmailService) getSubject(codeType string) string {
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
func (s *EmailService) getEmailBody(code, codeType string) string {
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

	// Build logo URL
	logoURL := s.baseURL + "/assets/logo.svg"
	if s.baseURL == "" {
		logoURL = "/assets/logo.svg" // Fallback for local development
	}

	return fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .content { background: #f9f9f9; padding: 30px; border-radius: 0 0 10px 10px; }
        .code { background: white; padding: 20px; text-align: center; font-size: 32px; letter-spacing: 5px; font-weight: bold; margin: 20px 0; border-radius: 5px; border: 2px dashed #667eea; }
        .footer { text-align: center; margin-top: 30px; color: #999; font-size: 12px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="content">
            <p>您好，您正在使用验证码进行%s操作。</p>
            <p>您的验证码是：</p>
            <div class="code">%s</div>
            <div style="text-align: center; margin: 30px 0;">
                <img src="%s" alt="iBookFS" style="width: 100px; height: auto; opacity: 0.85;" />
            </div>
            <p>验证码有效期为 <strong>5分钟</strong>，请尽快完成验证。</p>
            <p>如果这不是您的操作，请忽略此邮件。</p>
        </div>
        <div class="footer">
            <p>此邮件由系统自动发送，请勿直接回复。</p>
            <p>&copy; 2025 iBookFS. All rights reserved.</p>
        </div>
    </div>
</body>
</html>`, action, code, logoURL)
}
