// Package email provides transactional email delivery. The current concrete
// implementation is SMTP (works with a Gmail account + App Password), but
// callers depend only on a narrow Send method, so swapping in Resend/SES/etc.
// later is a drop-in change.
package email

import (
	"context"
	"fmt"
	"net/smtp"
	"strings"
)

// Config holds SMTP connection details. For Gmail: Host=smtp.gmail.com,
// Port=587, Username/From = your address, Password = a 16-char App Password
// (not your normal Google password).
type Config struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
	FromName string
}

// Configured reports whether enough is set to actually send mail. When false,
// the app runs without email (verification codes are generated but not sent),
// which keeps local dev working without SMTP credentials.
func (c Config) Configured() bool {
	return c.Host != "" && c.Port != 0 && c.From != ""
}

type SMTPMailer struct {
	cfg Config
}

func NewSMTPMailer(cfg Config) *SMTPMailer {
	return &SMTPMailer{cfg: cfg}
}

// Send delivers a multipart (plain-text + HTML) email synchronously. The
// context is accepted for interface symmetry; net/smtp itself is blocking.
func (m *SMTPMailer) Send(_ context.Context, to, subject, textBody, htmlBody string) error {
	addr := fmt.Sprintf("%s:%d", m.cfg.Host, m.cfg.Port)

	var auth smtp.Auth
	if m.cfg.Username != "" {
		auth = smtp.PlainAuth("", m.cfg.Username, m.cfg.Password, m.cfg.Host)
	}

	fromHeader := m.cfg.From
	if m.cfg.FromName != "" {
		fromHeader = fmt.Sprintf("%s <%s>", m.cfg.FromName, m.cfg.From)
	}

	msg := buildMIME(fromHeader, to, subject, textBody, htmlBody)
	// SendMail upgrades to STARTTLS automatically when the server advertises it
	// (Gmail on :587 does), so credentials aren't sent in the clear.
	return smtp.SendMail(addr, auth, m.cfg.From, []string{to}, msg)
}

// buildMIME assembles a multipart/alternative message so clients that can't
// render HTML fall back to the plain-text part.
func buildMIME(from, to, subject, textBody, htmlBody string) []byte {
	boundary := "tiq_boundary_9f2c1a7b"
	var b strings.Builder
	b.WriteString(fmt.Sprintf("From: %s\r\n", from))
	b.WriteString(fmt.Sprintf("To: %s\r\n", to))
	b.WriteString(fmt.Sprintf("Subject: %s\r\n", subject))
	b.WriteString("MIME-Version: 1.0\r\n")
	b.WriteString(fmt.Sprintf("Content-Type: multipart/alternative; boundary=%q\r\n", boundary))
	b.WriteString("\r\n")

	b.WriteString(fmt.Sprintf("--%s\r\n", boundary))
	b.WriteString("Content-Type: text/plain; charset=\"UTF-8\"\r\n\r\n")
	b.WriteString(textBody)
	b.WriteString("\r\n")

	b.WriteString(fmt.Sprintf("--%s\r\n", boundary))
	b.WriteString("Content-Type: text/html; charset=\"UTF-8\"\r\n\r\n")
	b.WriteString(htmlBody)
	b.WriteString("\r\n")

	b.WriteString(fmt.Sprintf("--%s--\r\n", boundary))
	return []byte(b.String())
}
