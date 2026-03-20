package netease

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/smtp"
	"strings"

	"github.com/mgcis-cn/ibookfs/pkg/email/sources"
	"github.com/mgcis-cn/ibookfs/pkg/options"
)

const SourceKind sources.Kind = "netease"

func init() {
	sources.Register(SourceKind, func(cfg sources.Config) (sources.Email, error) {
		opts, ok := cfg.(*options.NeteaseOptions)
		if !ok {
			return nil, fmt.Errorf("netease: invalid config type %T", cfg)
		}
		return New(opts)
	})
}

type Email struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	From     string `json:"from"`
	FromName string `json:"from_name"`
}

func New(opts *options.NeteaseOptions) (*Email, error) {
	if opts == nil {
		return nil, fmt.Errorf("netease: options is nil")
	}
	return &Email{
		Host:     opts.Host,
		Port:     opts.Port,
		User:     opts.User,
		Password: opts.Password,
		From:     opts.From,
		FromName: opts.FromName,
	}, nil
}

func (e *Email) Kind() sources.Kind {
	return SourceKind
}

// IsEnabled checks if the email provider is properly configured.
func (e *Email) IsEnabled() bool {
	return e.Host != "" && e.Port > 0
}

// Send sends an email with the given parameters.
func (e *Email) Send(ctx context.Context, to, subject, body string) error {
	if !e.IsEnabled() {
		return fmt.Errorf("netease email is disabled")
	}
	addr := fmt.Sprintf("%s:%d", e.Host, e.Port)
	msg := e.buildMessage(to, subject, body)
	return e.sendSMTP(addr, to, msg)
}

// sendSMTP sends email via SMTP with SSL/TLS support.
func (e *Email) sendSMTP(addr, to string, msg []byte) error {
	var auth smtp.Auth
	var client *smtp.Client
	var err error

	if e.Port == 465 || e.Port == 994 || e.Port == 587 {
		tlsConfig := &tls.Config{ServerName: e.Host, MinVersion: tls.VersionTLS12}
		conn, err := tls.Dial("tcp", addr, tlsConfig)
		if err != nil {
			return fmt.Errorf("failed to connect: %w", err)
		}
		client, err = smtp.NewClient(conn, e.Host)
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}
	} else {
		client, err = smtp.Dial(addr)
		if err != nil {
			return fmt.Errorf("failed to connect: %w", err)
		}
		if ok, _ := client.Extension("STARTTLS"); ok {
			if err := client.StartTLS(&tls.Config{ServerName: e.Host, MinVersion: tls.VersionTLS12}); err != nil {
				return fmt.Errorf("failed to start TLS: %w", err)
			}
		}
	}
	defer client.Close()

	auth = smtp.PlainAuth("", e.User, e.Password, e.Host)
	if err := client.Auth(auth); err != nil {
		return fmt.Errorf("auth failed: %w", err)
	}
	if err := client.Mail(e.From); err != nil {
		return fmt.Errorf("mail failed: %w", err)
	}
	if err := client.Rcpt(to); err != nil {
		return fmt.Errorf("rcpt failed: %w", err)
	}
	writer, err := client.Data()
	if err != nil {
		return fmt.Errorf("data failed: %w", err)
	}
	defer writer.Close()
	_, err = writer.Write(msg)
	return err
}

// buildMessage builds a raw email message.
func (e *Email) buildMessage(to, subject, body string) []byte {
	var msg strings.Builder
	msg.WriteString(fmt.Sprintf("From: %s <%s>\r\n", e.FromName, e.From))
	msg.WriteString(fmt.Sprintf("To: %s\r\n", to))
	msg.WriteString(fmt.Sprintf("Subject: %s\r\n", subject))
	msg.WriteString("MIME-Version: 1.0\r\n")
	msg.WriteString("Content-Type: text/html; charset=UTF-8\r\n\r\n")
	msg.WriteString(body)
	return []byte(msg.String())
}
