package email

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/psyconnect/notificationservice/pkg/config"
	"gopkg.in/gomail.v2"
)

type EmailService struct {
	config    *config.Config
	dialer    *gomail.Dialer
	templates map[string]string
}

func NewEmailService(cfg *config.Config) *EmailService {
	dialer := gomail.NewDialer(
		cfg.SMTPHost,
		cfg.SMTPPort,
		cfg.SMTPUsername,
		cfg.SMTPPassword,
	)

	service := &EmailService{
		config:    cfg,
		dialer:    dialer,
		templates: make(map[string]string),
	}

	// Load templates
	service.loadTemplates()

	return service
}

// loadTemplates loads HTML email templates from files
func (s *EmailService) loadTemplates() {
	templates := map[string]string{
		"verified":       "templates/verified.html",
		"reset-password": "templates/reset-password.html",
		"account-update": "templates/account-update.html",
	}

	for name, path := range templates {
		content, err := os.ReadFile(path)
		if err != nil {
			log.Printf("Failed to load template %s: %v", name, err)
			continue
		}
		s.templates[name] = string(content)
		log.Printf("Loaded template: %s", name)
	}
}

// SendActivationEmail sends account activation email
func (s *EmailService) SendActivationEmail(email, username, code, fullname string) error {
	subject := "Verify Your Account - PsyConnect"

	// Use template if available, otherwise fallback
	var body string
	if template, ok := s.templates["verified"]; ok {
		body = strings.ReplaceAll(template, "{USERNAME}", username)
		body = strings.ReplaceAll(body, "{CODE}", code)
	} else {
		body = s.buildActivationEmailFallback(username, code, fullname)
	}

	return s.sendEmail(email, subject, body)
}

// SendPasswordResetEmail sends password reset email
func (s *EmailService) SendPasswordResetEmail(email, username, code string) error {
	subject := "Reset Your Password - PsyConnect"

	// Use template if available, otherwise fallback
	var body string
	if template, ok := s.templates["reset-password"]; ok {
		body = strings.ReplaceAll(template, "{EMAIL}", email)
		body = strings.ReplaceAll(body, "{TOKEN}", code)
		body = strings.ReplaceAll(body, "{USERNAME}", username)
	} else {
		body = s.buildPasswordResetEmailFallback(username, code)
	}

	return s.sendEmail(email, subject, body)
}

// SendAccountChangeEmail sends account change notification
func (s *EmailService) SendAccountChangeEmail(email, username string) error {
	subject := "Account Update Notification - PsyConnect"

	// Use template if available, otherwise fallback
	var body string
	if template, ok := s.templates["account-update"]; ok {
		body = strings.ReplaceAll(template, "{USERNAME}", username)
	} else {
		body = s.buildAccountChangeEmailFallback(username)
	}

	return s.sendEmail(email, subject, body)
}

// sendEmail is the core email sending function
func (s *EmailService) sendEmail(to, subject, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", s.config.SMTPFrom)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	if err := s.dialer.DialAndSend(m); err != nil {
		log.Printf("ailed to send email to %s: %v", to, err)
		return err
	}

	log.Printf("Email sent successfully to %s", to)
	return nil
}

// Fallback templates (in case files are not found)

func (s *EmailService) buildActivationEmailFallback(username, code, fullname string) string {
	return fmt.Sprintf(`
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <title>Verify Your Account</title>
    <style>
      body {
        font-family: "Quicksand", sans-serif;
        background-color: #f4f4f5;
        margin: 0;
        padding: 0;
        color: #171717;
      }
      .email__container {
        max-width: 520px;
        margin: 40px auto;
        background: #ffffff;
        border-radius: 8px;
        box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
        padding: 24px 32px;
        text-align: center;
      }
      .email__header {
        font-size: 1.25rem;
        font-weight: 600;
        color: #18181b;
        margin-bottom: 16px;
      }
      .email__text {
        font-size: 0.95rem;
        color: #52525b;
        line-height: 1.6;
        margin: 8px 0;
      }
      .email__code {
        display: inline-block;
        background: #f4f4f5;
        border-radius: 6px;
        padding: 10px 20px;
        font-size: 1.25rem;
        font-weight: 600;
        color: #4a90e2;
        margin-top: 16px;
        letter-spacing: 2px;
      }
      .email__footer {
        border-top: 1px solid #e4e4e7;
        padding-top: 12px;
        font-size: 0.8rem;
        color: #71717a;
        margin-top: 24px;
      }
    </style>
  </head>
  <body>
    <div class="email__container">
      <div class="email__header">Verify Your Account</div>
      <p class="email__text">Hello <strong>%s</strong>,</p>
      <p class="email__text">
        Here is your verification code (expires in 10 minutes):
      </p>
      <div class="email__code">%s</div>
      <p class="email__text">
        If you did not request this code, please ignore this email.
      </p>
      <div class="email__footer">
        &copy; 2025 PsyConnect. All rights reserved.
      </div>
    </div>
  </body>
</html>
`, username, code)
}

func (s *EmailService) buildPasswordResetEmailFallback(username, code string) string {
	return fmt.Sprintf(`
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <title>Reset Your Password</title>
    <style>
      body {
        font-family: "Quicksand", sans-serif;
        background-color: #f4f4f5;
        margin: 0;
        padding: 0;
        color: #171717;
      }
      .email__container {
        max-width: 520px;
        margin: 40px auto;
        background: #fff;
        border-radius: 8px;
        box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
        padding: 24px 32px;
        text-align: center;
      }
      .email__header {
        font-size: 1.25rem;
        font-weight: 600;
        color: #18181b;
        margin-bottom: 12px;
      }
      .email__text {
        font-size: 0.95rem;
        color: #52525b;
        line-height: 1.6;
        margin: 8px 0;
      }
      .email__code {
        display: inline-block;
        background: #f4f4f5;
        border-radius: 6px;
        padding: 10px 20px;
        font-size: 1.25rem;
        font-weight: 600;
        color: #dc2626;
        margin-top: 16px;
        letter-spacing: 2px;
      }
      .email__footer {
        border-top: 1px solid #e4e4e7;
        padding-top: 12px;
        font-size: 0.8rem;
        color: #71717a;
        margin-top: 24px;
      }
    </style>
  </head>
  <body>
    <div class="email__container">
      <div class="email__header">Reset Your Password</div>
      <p class="email__text">
        We received a request to reset your password. Use the code below:
      </p>
      <div class="email__code">%s</div>
      <p class="email__text">This code will expire in 10 minutes.</p>
      <p class="email__text">
        If you did not request a password reset, you can safely ignore this email.
      </p>
      <div class="email__footer">
        &copy; 2025 PsyConnect. All rights reserved.
      </div>
    </div>
  </body>
</html>
`, code)
}

func (s *EmailService) buildAccountChangeEmailFallback(username string) string {
	return fmt.Sprintf(`
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <title>Account Update Notification</title>
    <style>
      body {
        font-family: "Quicksand", sans-serif;
        background-color: #f4f4f5;
        margin: 0;
        padding: 0;
        color: #171717;
      }
      .email__container {
        max-width: 600px;
        margin: 40px auto;
        background: #ffffff;
        border-radius: 8px;
        box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
        padding: 24px 32px;
      }
      .email__header {
        border-bottom: 1px solid #e4e4e7;
        padding-bottom: 12px;
        margin-bottom: 20px;
      }
      .email__header h2 {
        font-size: 1.25rem;
        font-weight: 600;
        color: #18181b;
        margin: 0;
      }
      .email__highlight {
        color: #4a90e2;
        font-weight: 600;
      }
      .email__body p {
        font-size: 0.95rem;
        color: #52525b;
        line-height: 1.6;
        margin: 0 0 12px;
      }
      .email__footer {
        border-top: 1px solid #e4e4e7;
        padding-top: 12px;
        font-size: 0.8rem;
        color: #71717a;
        text-align: center;
        margin-top: 24px;
      }
    </style>
  </head>
  <body>
    <div class="email__container">
      <div class="email__header">
        <h2>Hi <span class="email__highlight">%s</span>,</h2>
      </div>
      <div class="email__body">
        <p>
          Your account <strong>%s</strong> has recently been updated.
        </p>
        <p>
          If you didn't request this, please contact our support team immediately.
        </p>
        <p>Thank you for trusting <strong>PsyConnect</strong>.</p>
      </div>
      <div class="email__footer">
        &copy; 2025 PsyConnect. All rights reserved.
      </div>
    </div>
  </body>
</html>
`, username, username)
}
