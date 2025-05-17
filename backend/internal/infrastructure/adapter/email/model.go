package email

import (
	"github.com/huydq/test/internal/pkg/utils"
)

// Constants for email content types
const (
	ContentTypeText = "text/plain"
	ContentTypeHTML = "text/html"
)

// Email template IDs
const (
	TemplateID2FACode = "send_2fa_code"
)

// Email template files
const (
	TemplateFile2FACode = "2fa_code.tmpl"
)

// Email template folders
const (
	TemplateFolderAuth = "auth"
)

// Email subjects
const (
	Subject2FACode = "件名: ログイン確認コードのご案内"
)

// IEmailService defines the contract for email services
type IEmailService interface {
	// SendEmail sends an email using the specified data
	SendEmail(data EmailData) error

	// SendHTMLEmail is a convenience method for sending HTML emails
	SendHTMLEmail(data EmailData) error

	// SendPlainTextEmail is a convenience method for sending plain text emails
	SendPlainTextEmail(data EmailData) error

	// TestConnection tests the SMTP connection
	TestConnection() error
}

// MailConfig holds SMTP server configuration
type MailConfig struct {
	Host        string
	Port        int
	Username    string
	Password    string
	FromEmail   string
	FromName    string
	TemplateDir string
	UseTLS      bool
	UseAuth     bool
}

// EmailData contains the data needed to send an email
type EmailData struct {
	To             []string
	ToName         string
	Subject        string
	TemplateFile   string
	TemplateFolder string
	Data           map[string]any
	Cc             []string
	Bcc            []string
	ContentType    string             // "text/plain" or "text/html"
	Charset        utils.EmailCharset // "utf-8" or "ISO-2022-JP"
}

// NewEmailData creates a new EmailData with default values
func NewEmailData() EmailData {
	return EmailData{
		To:          []string{},
		Cc:          []string{},
		Bcc:         []string{},
		ContentType: ContentTypeText,
		Charset:     utils.CharsetUTF8,
		Data:        map[string]any{},
	}
}

// NewHTMLEmailData creates a new EmailData with HTML content type
func NewHTMLEmailData() EmailData {
	data := NewEmailData()
	data.ContentType = ContentTypeHTML
	return data
}
