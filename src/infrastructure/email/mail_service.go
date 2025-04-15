// src/infrastructure/email/mail_service.go

package email

import (
	"bytes"
	"fmt"
	"net/smtp"
	"os"
	"path/filepath"
	"text/template"

	"github.com/vnlab/makeshop-payment/src/infrastructure/config"
	"github.com/vnlab/makeshop-payment/src/infrastructure/logger"
)

// MailConfig holds SMTP server configuration
type MailConfig struct {
	Host        string
	Port        int
	Username    string
	Password    string
	FromEmail   string
	FromName    string
	TemplateDir string
}

// EmailData contains the data needed to send an email
type EmailData struct {
	To             string
	ToName         string
	Subject        string
	TemplateFile   string
	TemplateFolder string
	Data           map[string]interface{}
}

// MailService handles email sending functionality
type MailService struct {
	config *MailConfig
	logger logger.Logger
}

// NewMailService creates a new mail service with the given logger
func NewMailService(logger logger.Logger) (*MailService, error) {
	appConfig := config.GetConfig()

	mailConfig := &MailConfig{
		Host:        appConfig.SMTPHost,
		Port:        appConfig.SMTPPort,
		Username:    appConfig.SMTPUsername,
		Password:    appConfig.SMTPPassword,
		FromEmail:   appConfig.SMTPFromEmail,
		FromName:    appConfig.SMTPFromName,
		TemplateDir: appConfig.EmailTemplateDir,
	}

	return &MailService{
		config: mailConfig,
		logger: logger,
	}, nil
}

// SendEmail sends an email using the specified template and data
func (s *MailService) SendEmail(data EmailData) error {
	// Render template content
	content, err := s.renderEmailContent(data)
	if err != nil {
		return err
	}

	// Prepare email message
	message := s.prepareEmailMessage(data, content)
	
	// Send the email
	return s.sendSMTPEmail(data, message)
}

// prepareEmailMessage constructs the complete email with headers and content
func (s *MailService) prepareEmailMessage(data EmailData, content string) string {
	from := s.config.FromEmail
	if s.config.FromName != "" {
		from = fmt.Sprintf("%s <%s>", s.config.FromName, s.config.FromEmail)
	}

	// Build email headers
	headers := map[string]string{
		"From":         from,
		"To":           data.To,
		"Subject":      data.Subject,
		"MIME-Version": "1.0",
		"Content-Type": "text/plain; charset=UTF-8",
	}

	message := ""
	for key, value := range headers {
		message += fmt.Sprintf("%s: %s\r\n", key, value)
	}
	message += "\r\n" + content
	
	return message
}

// sendSMTPEmail handles the actual SMTP sending logic
func (s *MailService) sendSMTPEmail(data EmailData, message string) error {
	address := fmt.Sprintf("%s:%d", s.config.Host, s.config.Port)

	// Use different sending method based on host
	if s.isLocalhost() {
		return s.sendToLocalSMTP(address, data.To, message)
	} 
	
	// Use authenticated SMTP for production servers
	auth := smtp.PlainAuth("", s.config.Username, s.config.Password, s.config.Host)
	err := smtp.SendMail(address, auth, s.config.FromEmail, []string{data.To}, []byte(message))
	
	if err != nil {
		s.logger.Error("Failed to send email", map[string]interface{}{
			"to":    data.To,
			"error": err.Error(),
		})
		return err
	}

	s.logEmailSent(data)
	return nil
}

// isLocalhost checks if the SMTP server is a local development server
func (s *MailService) isLocalhost() bool {
	return s.config.Host == "mailhog" || s.config.Host == "localhost"
}

// sendToLocalSMTP handles sending to local SMTP servers (like mailhog)
func (s *MailService) sendToLocalSMTP(address, recipient, message string) error {
	c, err := smtp.Dial(address)
	if err != nil {
		s.logger.Error("Failed to connect to SMTP server", map[string]interface{}{
			"error": err.Error(),
		})
		return err
	}
	defer c.Close()

	if err = c.Mail(s.config.FromEmail); err != nil {
		s.logger.Error("Failed to set sender", map[string]interface{}{
			"error": err.Error(),
		})
		return err
	}

	if err = c.Rcpt(recipient); err != nil {
		s.logger.Error("Failed to set recipient", map[string]interface{}{
			"error": err.Error(),
		})
		return err
	}

	wc, err := c.Data()
	if err != nil {
		s.logger.Error("Failed to start data transmission", map[string]interface{}{
			"error": err.Error(),
		})
		return err
	}
	defer wc.Close()

	_, err = wc.Write([]byte(message))
	if err != nil {
		s.logger.Error("Failed to send email data", map[string]interface{}{
			"error": err.Error(),
		})
		return err
	}
	
	s.logEmailSent(EmailData{To: recipient})
	return nil
}

// logEmailSent logs successful email sending
func (s *MailService) logEmailSent(data EmailData) {
	s.logger.Info("Email sent successfully", map[string]interface{}{
		"to":       data.To,
		"subject":  data.Subject,
		"template": data.TemplateFile,
	})
}

// renderEmailContent renders an email template with the provided data
func (s *MailService) renderEmailContent(emailData EmailData) (string, error) {
	mainTemplatePath := filepath.Join(s.config.TemplateDir, emailData.TemplateFolder, emailData.TemplateFile)
	footerTemplatePath := filepath.Join(s.config.TemplateDir, "common", "footer.tmpl")

	// Read template files
	mainContent, err := os.ReadFile(mainTemplatePath)
	if err != nil {
		s.logger.Error("Failed to read main template", map[string]interface{}{
			"path":  mainTemplatePath,
			"error": err.Error(),
		})
		return "", fmt.Errorf("failed to read main template: %w", err)
	}

	// Footer is optional
	footerContent, err := os.ReadFile(footerTemplatePath)
	if err != nil {
		s.logger.Error("Failed to read footer template", map[string]interface{}{
			"path":  footerTemplatePath,
			"error": err.Error(),
		})
		footerContent = []byte("")
	}

	// Combine templates
	combinedContent := fmt.Sprintf(`{{define "main"}}%s{{end}}
{{define "footer"}}%s{{end}}

{{template "main" .}}
{{template "footer" .}}`, string(mainContent), string(footerContent))

	// Parse and execute template
	tmpl, err := template.New("email").Parse(combinedContent)
	if err != nil {
		s.logger.Error("Failed to parse template", map[string]interface{}{
			"error": err.Error(),
		})
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, emailData.Data); err != nil {
		s.logger.Error("Failed to execute template", map[string]interface{}{
			"error": err.Error(),
		})
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.String(), nil
}

// TestConnection tests the SMTP connection
func (s *MailService) TestConnection() error {
	address := fmt.Sprintf("%s:%d", s.config.Host, s.config.Port)
	client, err := smtp.Dial(address)
	if err != nil {
		return fmt.Errorf("failed to connect to SMTP server: %w", err)
	}
	defer client.Close()
	
	return nil
}
