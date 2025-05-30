package mail

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/smtp"
	"strings"

	"github.com/trysourcetool/onprem-portal/internal/config"
	"github.com/trysourcetool/onprem-portal/internal/logger"
)

type input struct {
	To       []string
	From     string
	FromName string
	Subject  string
	Body     string
}

func send(ctx context.Context, in input) error {
	// Build From header with proper format
	fromHeader := in.From
	if in.FromName != "" {
		fromHeader = fmt.Sprintf("%s <%s>", in.FromName, in.From)
	}

	msg := fmt.Sprintf("From: %s\r\n"+
		"To: %s\r\n"+
		"Subject: %s\r\n"+
		"MIME-Version: 1.0\r\n"+
		"Content-Type: text/plain; charset=UTF-8\r\n"+
		"Content-Transfer-Encoding: 7bit\r\n"+
		"\r\n"+
		"%s\r\n", fromHeader, strings.Join(in.To, ","), in.Subject, in.Body)

	if config.Config.Env == config.EnvLocal {
		// In local environment, just log the email content
		logger.Logger.Sugar().Debug("================= EMAIL CONTENT =================")
		logger.Logger.Sugar().Debug(msg)
		logger.Logger.Sugar().Debug("================= EMAIL CONTENT =================")

		// Don't actually send in local environment
		return nil
	}

	cfg := config.Config.SMTP
	auth := smtp.PlainAuth("", cfg.Username, cfg.Password, cfg.Host)
	addr := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)
	tlsConf := &tls.Config{
		ServerName: cfg.Host,
		MinVersion: tls.VersionTLS12,
	}

	var client *smtp.Client
	var err error

	if cfg.UseTLS {
		// Implicit TLS (usually port 465)
		dialer := &tls.Dialer{
			Config: tlsConf,
		}
		conn, err := dialer.DialContext(ctx, "tcp", addr)
		if err != nil {
			return fmt.Errorf("failed to create TLS connection: %w", err)
		}
		defer conn.Close()

		client, err = smtp.NewClient(conn, cfg.Host)
		if err != nil {
			return fmt.Errorf("failed to create SMTP client: %w", err)
		}
	} else {
		// Start with plain connection and upgrade to TLS using STARTTLS (usually port 587)
		d := net.Dialer{}
		conn, err := d.DialContext(ctx, "tcp", addr)
		if err != nil {
			return fmt.Errorf("failed to create connection: %w", err)
		}
		defer conn.Close()

		client, err = smtp.NewClient(conn, cfg.Host)
		if err != nil {
			return fmt.Errorf("failed to create SMTP client: %w", err)
		}

		// Try STARTTLS if available
		if ok, _ := client.Extension("STARTTLS"); ok {
			if err := client.StartTLS(tlsConf); err != nil {
				return fmt.Errorf("failed to start TLS: %w", err)
			}
		}
	}
	defer client.Quit()

	if err = client.Auth(auth); err != nil {
		return fmt.Errorf("failed to authenticate: %w", err)
	}

	if err := client.Mail(in.From); err != nil {
		return fmt.Errorf("failed to set FROM address: %w", err)
	}

	for _, addr := range in.To {
		if err := client.Rcpt(addr); err != nil {
			return fmt.Errorf("failed to set TO address: %w", err)
		}
	}

	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("failed to create message writer: %w", err)
	}
	defer w.Close()

	if _, err = w.Write([]byte(msg)); err != nil {
		return fmt.Errorf("failed to write message: %w", err)
	}

	return nil
}

const fromName = "Sourcetool Team"

func SendMagicLinkEmail(ctx context.Context, email, firstName, url string) error {
	subject := "Log in to Sourcetool On-premise portal"

	content := fmt.Sprintf(`Hi %s,

Here's your magic link to log in to your Sourcetool On-premise portal. Click the link below to access your account securely without a password:

%s

- This link will expire in 15 minutes for security reasons.
- If you didn't request this link, you can safely ignore this email.

Thank you for using Sourcetool!

The Sourcetool Team`, firstName, url)

	if err := send(ctx, input{
		From:     config.Config.SMTP.FromEmail,
		FromName: fromName,
		To:       []string{email},
		Subject:  subject,
		Body:     content,
	}); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}

func SendUpdateEmailInstructions(ctx context.Context, to, firstName, url string) error {
	subject := "[Sourcetool] Confirm your new email address"
	content := fmt.Sprintf(`Hi %s,

We received a request to change the email address associated with your Sourcetool account. To ensure the security of your account, we need you to verify your new email address.

Please click the following link within the next 24 hours to confirm your email change:
%s

Thank you for being a part of the Sourcetool community!
Regards,

The Sourcetool Team`,
		firstName,
		url,
	)

	return send(ctx, input{
		From:     config.Config.SMTP.FromEmail,
		FromName: fromName,
		To:       []string{to},
		Subject:  subject,
		Body:     content,
	})
}
