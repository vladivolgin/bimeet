package mailer

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/smtp"
	"time"
)

const (
	smtpDialTimeout = 10 * time.Second
	httpTimeout     = 10 * time.Second
	resendEndpoint  = "https://api.resend.com/emails"
)

type Config struct {
	// Resend HTTP API (preferred when set — works on any network, port 443)
	ResendAPIKey string
	// MailFrom — canonical "From:" address used for both Resend and SMTP.
	// Falls back to SMTPFrom if empty.
	MailFrom string
	// SMTP (used only if ResendAPIKey is empty; if SMTPHost is also empty, mail is only logged)
	SMTPHost string
	SMTPPort string
	SMTPUser string
	SMTPPass string
	SMTPFrom string

	FrontendURL string
}

type Mailer struct {
	cfg  Config
	from string
	http *http.Client
}

func New(cfg Config) *Mailer {
	from := cfg.MailFrom
	if from == "" {
		from = cfg.SMTPFrom
	}
	return &Mailer{
		cfg:  cfg,
		from: from,
		http: &http.Client{Timeout: httpTimeout},
	}
}

func (m *Mailer) SendInvite(toEmail, eventTitle, organizerName, inviteToken string) error {
	inviteURL := m.cfg.FrontendURL + "/invite/" + inviteToken
	subject := fmt.Sprintf("Вас пригласили на встречу «%s»", eventTitle)
	body := fmt.Sprintf(
		"Привет!\n\n%s приглашает вас на встречу «%s».\n\n"+
			"Перейдите по ссылке, чтобы принять или отклонить приглашение:\n%s",
		organizerName, eventTitle, inviteURL,
	)
	return m.send(toEmail, subject, body)
}

func (m *Mailer) SendPasswordReset(toEmail, resetToken string) error {
	resetURL := m.cfg.FrontendURL + "/reset-password?token=" + resetToken
	subject := "Восстановление пароля"
	body := fmt.Sprintf(
		"Привет!\n\nВы запросили восстановление пароля.\n\n"+
			"Перейдите по ссылке, чтобы задать новый пароль:\n%s\n\n"+
			"Если вы не запрашивали восстановление, просто проигнорируйте это письмо.",
		resetURL,
	)
	return m.send(toEmail, subject, body)
}

func (m *Mailer) send(toEmail, subject, body string) error {
	switch {
	case m.cfg.ResendAPIKey != "":
		return m.sendResend(toEmail, subject, body)
	case m.cfg.SMTPHost != "":
		return m.sendSMTP(toEmail, subject, body)
	default:
		fmt.Printf("[mailer] no transport configured — would send to %s: %s\n", toEmail, subject)
		return nil
	}
}

// ─── Resend HTTP API ───────────────────────────────────────────────────────

type resendRequest struct {
	From    string   `json:"from"`
	To      []string `json:"to"`
	Subject string   `json:"subject"`
	Text    string   `json:"text"`
}

type resendError struct {
	Name    string `json:"name"`
	Message string `json:"message"`
}

func (m *Mailer) sendResend(toEmail, subject, body string) error {
	if m.from == "" {
		return fmt.Errorf("MAIL_FROM is required for Resend")
	}

	payload, err := json.Marshal(resendRequest{
		From:    m.from,
		To:      []string{toEmail},
		Subject: subject,
		Text:    body,
	})
	if err != nil {
		return fmt.Errorf("marshal resend payload: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, resendEndpoint, bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("new request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+m.cfg.ResendAPIKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := m.http.Do(req)
	if err != nil {
		return fmt.Errorf("resend request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return nil
	}

	respBody, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
	var apiErr resendError
	if err := json.Unmarshal(respBody, &apiErr); err == nil && apiErr.Message != "" {
		return fmt.Errorf("resend %d %s: %s", resp.StatusCode, apiErr.Name, apiErr.Message)
	}
	return fmt.Errorf("resend %d: %s", resp.StatusCode, string(respBody))
}

// ─── SMTP ─────────────────────────────────────────────────────────────────

func (m *Mailer) sendSMTP(toEmail, subject, body string) error {
	from := m.from
	if from == "" {
		from = m.cfg.SMTPUser
	}

	msg := fmt.Sprintf(
		"From: %s\r\nTo: %s\r\nSubject: %s\r\nContent-Type: text/plain; charset=UTF-8\r\n\r\n%s",
		from, toEmail, subject, body,
	)

	addr := m.cfg.SMTPHost + ":" + m.cfg.SMTPPort

	conn, err := net.DialTimeout("tcp", addr, smtpDialTimeout)
	if err != nil {
		return fmt.Errorf("dial %s: %w", addr, err)
	}

	c, err := smtp.NewClient(conn, m.cfg.SMTPHost)
	if err != nil {
		_ = conn.Close()
		return fmt.Errorf("smtp client: %w", err)
	}
	defer c.Close()

	if ok, _ := c.Extension("STARTTLS"); ok {
		if err := c.StartTLS(&tls.Config{ServerName: m.cfg.SMTPHost}); err != nil {
			return fmt.Errorf("starttls: %w", err)
		}
	}

	if m.cfg.SMTPUser != "" {
		auth := smtp.PlainAuth("", m.cfg.SMTPUser, m.cfg.SMTPPass, m.cfg.SMTPHost)
		if err := c.Auth(auth); err != nil {
			return fmt.Errorf("auth: %w", err)
		}
	}

	if err := c.Mail(from); err != nil {
		return fmt.Errorf("mail from: %w", err)
	}
	if err := c.Rcpt(toEmail); err != nil {
		return fmt.Errorf("rcpt to: %w", err)
	}

	w, err := c.Data()
	if err != nil {
		return fmt.Errorf("data: %w", err)
	}
	if _, err := w.Write([]byte(msg)); err != nil {
		return fmt.Errorf("write body: %w", err)
	}
	if err := w.Close(); err != nil {
		return fmt.Errorf("close data: %w", err)
	}

	return c.Quit()
}
