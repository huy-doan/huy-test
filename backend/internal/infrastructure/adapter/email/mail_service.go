package email

import (
	"fmt"
	"reflect"

	"github.com/huydq/test/internal/pkg/config"
	"github.com/huydq/test/internal/pkg/logger"
	"github.com/huydq/test/internal/pkg/utils"
)

// TwoFACodeEmailData contains data required for 2FA code emails
type TwoFACodeEmailData struct {
	Email          string
	ToName         string
	Token          string
	TokenExpiryMin int
}

// MailService handles email sending functionality by orchestrating various components
type MailService struct {
	config           *MailConfig
	logger           logger.Logger
	smtpClient       *SMTPClient
	composer         *MessageComposer
	templateRenderer *TemplateRenderer
}

// NewMailService creates a new mail service
func NewMailService() (*MailService, error) {
	// Get configuration from application config
	appConfig := config.GetConfig()
	appLogger := logger.GetLogger()

	// Create mail config
	mailConfig := &MailConfig{
		Host:        appConfig.SMTPHost,
		Port:        appConfig.SMTPPort,
		Username:    appConfig.SMTPUsername,
		Password:    appConfig.SMTPPassword,
		FromEmail:   appConfig.SMTPFromEmail,
		FromName:    appConfig.SMTPFromName,
		TemplateDir: appConfig.EmailTemplateDir,
		UseTLS:      appConfig.SMTPUseTLS,
		UseAuth:     appConfig.SMTPUseAuth,
	}

	// Create component instances
	smtpClient := NewSMTPClient(mailConfig, appLogger)
	composer := NewMessageComposer(mailConfig)
	templateRenderer := NewTemplateRenderer(mailConfig.TemplateDir, appLogger)

	return &MailService{
		config:           mailConfig,
		logger:           appLogger,
		smtpClient:       smtpClient,
		composer:         composer,
		templateRenderer: templateRenderer,
	}, nil
}

// SendEmail sends an email using the specified template and data
func (s *MailService) SendEmail(data EmailData) error {
	// Set defaults if not specified
	if data.ContentType == "" {
		data.ContentType = ContentTypeText
	}
	if data.Charset == "" {
		data.Charset = utils.CharsetUTF8
	}
	if data.Cc == nil {
		data.Cc = []string{}
	}
	if data.Bcc == nil {
		data.Bcc = []string{}
	}

	// Canonicalize email addresses and text
	for i := range data.To {
		data.To[i] = s.composer.CanonicalizeEmail(data.To[i])
	}
	for i := range data.Cc {
		data.Cc[i] = s.composer.CanonicalizeEmail(data.Cc[i])
	}
	for i := range data.Bcc {
		data.Bcc[i] = s.composer.CanonicalizeEmail(data.Bcc[i])
	}

	// Canonicalize subject
	data.Subject = s.composer.CanonicalizeText(data.Subject)

	// Render template content
	content, err := s.templateRenderer.RenderTemplate(data.TemplateFolder, data.TemplateFile, data.Data)
	if err != nil {
		return fmt.Errorf("failed to render template: %w", err)
	}

	// Canonicalize content
	content = s.composer.CanonicalizeText(content)

	// Convert to ISO-2022-JP if specified
	if data.Charset == utils.CharsetISO2022JP {
		content, err = utils.ToISO2022JPWithFallback(content)
		if err != nil {
			return fmt.Errorf("failed to convert content to ISO-2022-JP: %w", err)
		}
	}

	// Prepare from address
	fromName := s.composer.CanonicalizeSenderName(s.config.FromName)
	from, err := s.composer.FormatFromAddress(fromName, s.config.FromEmail, data.Charset)

	if err != nil {
		return fmt.Errorf("failed to format from address: %w", err)
	}

	// Compose complete message
	message, err := s.composer.ComposeMessage(from, data, content)
	if err != nil {
		return fmt.Errorf("failed to compose message: %w", err)
	}

	// Send the email through SMTP client
	if err := s.smtpClient.SendRawMessage(from, data.To, data.Cc, data.Bcc, message); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}

// SendHTMLEmail is a convenience method for sending HTML emails
func (s *MailService) SendHTMLEmail(data EmailData) error {
	data.ContentType = ContentTypeHTML
	return s.SendEmail(data)
}

// SendPlainTextEmail is a convenience method for sending plain text emails
func (s *MailService) SendPlainTextEmail(data EmailData) error {
	data.ContentType = ContentTypeText
	return s.SendEmail(data)
}

// TestConnection tests the SMTP connection
func (s *MailService) TestConnection() error {
	return s.smtpClient.TestConnection()
}

// sendMailByTemplateID sends an email using a template ID
func (s *MailService) SendMailByTemplateID(templateID string, payload any) {
	switch templateID {
	case TemplateID2FACode:
		emailData := s.prepareTwoFACodeEmail(payload)

		if !reflect.DeepEqual(emailData, EmailData{}) {
			err := s.SendEmail(emailData)

			if err != nil {
				s.logger.Error("[Send mail 2FA Code fail]", map[string]any{
					"error": err.Error(),
				})
			}
		}
	}
}

// prepareTwoFACodeEmail prepares the email data for the 2FA code email
func (s *MailService) prepareTwoFACodeEmail(payload any) EmailData {
	data, ok := payload.(TwoFACodeEmailData)
	if !ok {
		s.logger.Error("[Send mail 2FA Code fail]", map[string]any{
			"error": "failed to cast payload to TwoFACodeEmailData",
		})

		return EmailData{}
	}

	emailData := EmailData{
		To:             []string{data.Email},
		Subject:        Subject2FACode,
		TemplateFile:   TemplateFile2FACode,
		TemplateFolder: TemplateFolderAuth,
		Data: map[string]any{
			"ToName":    data.ToName,
			"Code":      data.Token,
			"ExpiresIn": data.TokenExpiryMin,
		},
	}

	return emailData
}
