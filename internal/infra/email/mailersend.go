package email

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// mailerSendEndpoint is MailerSend's single-transactional-email API.
// Docs: https://developers.mailersend.com/api/v1/email.html
const mailerSendEndpoint = "https://api.mailersend.com/v1/email"

// MailerSendConfig holds what the MailerSend HTTP API needs. APIKey is an API
// token created in the MailerSend dashboard; From must be an address on a
// domain that has been added and verified there, or every send is rejected.
type MailerSendConfig struct {
	APIKey   string
	From     string
	FromName string
}

// Configured reports whether enough is set to send through MailerSend.
func (c MailerSendConfig) Configured() bool {
	return c.APIKey != "" && c.From != ""
}

// MailerSendMailer delivers mail over MailerSend's REST API instead of SMTP,
// satisfying the same narrow Send interface the usecases depend on.
type MailerSendMailer struct {
	cfg    MailerSendConfig
	client *http.Client
}

func NewMailerSendMailer(cfg MailerSendConfig) *MailerSendMailer {
	return &MailerSendMailer{
		cfg: cfg,
		// The API normally answers in well under a second; the timeout only
		// bounds how long a caller can be stuck when MailerSend is degraded.
		client: &http.Client{Timeout: 15 * time.Second},
	}
}

// Send posts one email to the MailerSend API. MailerSend accepts the message
// with 202 (delivery happens asynchronously on their side); anything else is
// returned as an error with the response body, which carries their validation
// details (e.g. unverified sender domain, trial-account recipient limits).
func (m *MailerSendMailer) Send(ctx context.Context, to, subject, textBody, htmlBody string) error {
	payload := map[string]any{
		"from":    map[string]string{"email": m.cfg.From, "name": m.cfg.FromName},
		"to":      []map[string]string{{"email": to}},
		"subject": subject,
		"text":    textBody,
		"html":    htmlBody,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("mailersend: encoding request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, mailerSendEndpoint, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("mailersend: building request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+m.cfg.APIKey)

	resp, err := m.client.Do(req)
	if err != nil {
		return fmt.Errorf("mailersend: sending request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusAccepted || resp.StatusCode == http.StatusOK {
		return nil
	}
	detail, _ := io.ReadAll(io.LimitReader(resp.Body, 2048))
	return fmt.Errorf("mailersend: API returned %d: %s", resp.StatusCode, detail)
}
