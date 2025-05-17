package email

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net"
	"net/smtp"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"text/template"
	"time"

	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"

	"github.com/huydq/test/src/infrastructure/config"
	"github.com/huydq/test/src/infrastructure/logger"
	"github.com/huydq/test/src/lib/utils"
	"github.com/huydq/test/src/lib/utils/static"
)

// Constants for email content types
const (
	ContentTypeText = "text/plain"
	ContentTypeHTML = "text/html"
)

// Constants for character sets
const (
	CharsetUTF8      = "utf-8"
	CharsetISO2022JP = "ISO-2022-JP"
)

// Constants for email formatting
const (
	MaxLineLength = 76 // RFC 5322 recommended line length
)

// メール一行の最大文字数（RFC 5322をもとに設定）
const (
	characterLimitForOneLineOfEmail = 990
)

type (
	contentType string
	charset     string
)

func (c contentType) String() string {
	return string(c)
}

func (char charset) String() string {
	return string(char)
}

func (char charset) Length() int {
	switch char {
	case CharsetUTF8:
		return characterLimitForOneLineOfEmail / 6 // Unicodeでは一文字が最大6バイトになるため
	case CharsetISO2022JP:
		return 36
	}
	panic("unreachable panic")
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
	To             string
	ToName         string
	Subject        string
	TemplateFile   string
	TemplateFolder string
	Data           map[string]interface{}
	Cc             []string
	Bcc            []string
	ContentType    string  // "text/plain" or "text/html"
	Charset        charset // "utf-8" or "ISO-2022-JP"
}

// MailService handles email sending functionality
type MailService struct {
	config *MailConfig
	logger logger.Logger
	mutex  sync.Mutex
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
		UseTLS:      appConfig.SMTPUseTLS,
		UseAuth:     appConfig.SMTPUseAuth,
	}

	return &MailService{
		config: mailConfig,
		logger: logger,
	}, nil
}

func (s *MailService) encode(char charset, str string) (ret string, err error) {
	if char == CharsetISO2022JP {
		str, err = s.toISO2022JP(str)
		if err != nil {
			return "", err
		}
	}
	ret = base64.StdEncoding.EncodeToString([]byte(str))
	return
}

// toISO2022JP converts UTF-8 text to ISO-2022-JP encoding
func (s *MailService) toISO2022JP(text string) (string, error) {
	r := transform.NewReader(strings.NewReader(text), japanese.ISO2022JP.NewEncoder())
	b, err := io.ReadAll(r)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// toIso2022JpOrReplacedChar converts text with character replacement
func (s *MailService) toIso2022JpOrReplacedChar(text string) (string, error) {
	body := &bytes.Buffer{}
	// Replace characters that can't be converted with question mark
	w := &utils.RuneWriter{W: transform.NewWriter(body, japanese.ISO2022JP.NewEncoder())}
	_, err := io.Copy(w, strings.NewReader(text))
	if err != nil {
		return "", err
	}
	return body.String(), nil
}

// SendEmail sends an email using the specified template and data
func (s *MailService) SendEmail(data EmailData) error {
	if data.ContentType == "" {
		data.ContentType = ContentTypeText
	}
	if data.Charset == "" {
		data.Charset = CharsetUTF8
	}
	if data.Cc == nil {
		data.Cc = []string{}
	}
	if data.Bcc == nil {
		data.Bcc = []string{}
	}

	// Canonicalize email addresses, subject and content
	data.To = s.canonicalizeEmail(data.To)
	for i := range data.Cc {
		data.Cc[i] = s.canonicalizeEmail(data.Cc[i])
	}
	for i := range data.Bcc {
		data.Bcc[i] = s.canonicalizeEmail(data.Bcc[i])
	}
	data.Subject = s.canonicalizeText(data.Subject)

	// Render template content
	content, err := s.renderEmailContent(data)
	if err != nil {
		return err
	}

	// Canonicalize content
	content = s.canonicalizeText(content)

	// Convert to ISO-2022-JP if specified
	if data.Charset == CharsetISO2022JP {
		content, err = s.toIso2022JpOrReplacedChar(content)
		if err != nil {
			return fmt.Errorf("failed to convert content to ISO-2022-JP: %w", err)
		}
	}

	// Prepare email message
	message, err := s.prepareEmailMessage(data, content)
	if err != nil {
		return err
	}

	// Send the email
	return s.sendSMTPEmail(data, message)
}

// prepareEmailMessage constructs the complete email with headers and content
func (s *MailService) prepareEmailMessage(data EmailData, content string) (string, error) {
	from := s.formatFromAddress(s.canonicalizeSenderName(s.config.FromName), s.config.FromEmail, data.Charset.String())

	// Encode subject
	subject, err := s.encodeSubject(data.Subject, data.Charset)
	if err != nil {
		return "", err
	}

	// Build email headers
	headers := s.makeHeader(from, data.To, data.Cc, data.Bcc, subject, data.ContentType, data.Charset.String())

	message := ""
	for key, value := range headers {
		message += fmt.Sprintf("%s: %s\r\n", key, value)
	}

	// Add content encoding header and content
	message += "Content-Transfer-Encoding: base64\r\n\r\n"

	// Encode and format content
	encodedContent := base64.StdEncoding.EncodeToString([]byte(content))
	message += s.cutAndAddCrlf(encodedContent)

	return message, nil
}

// makeHeader creates a map of email headers
func (s *MailService) makeHeader(from, to string, cc, bcc []string, subject, contentType, charset string) map[string]string {
	headers := map[string]string{
		"From":         from,
		"To":           to,
		"Subject":      subject,
		"MIME-Version": "1.0",
		"Content-Type": fmt.Sprintf("%s; charset=%s", contentType, charset),
		"Date":         time.Now().Format(time.RFC1123Z),
	}

	if len(cc) > 0 {
		headers["Cc"] = strings.Join(cc, ", ")
	}

	if len(bcc) > 0 {
		headers["Bcc"] = strings.Join(bcc, ", ")
	}

	return headers
}

// formatFromAddress formats the From address with proper encoding
func (s *MailService) formatFromAddress(name, email, charset string) string {
	if name == "" {
		return email
	}

	if charset == CharsetISO2022JP {
		japName, err := s.toISO2022JP(name)
		if err == nil {
			encodedName := base64.StdEncoding.EncodeToString([]byte(japName))
			return fmt.Sprintf("=?%s?B?%s?= <%s>", charset, encodedName, email)
		} else {
			return fmt.Sprintf("%s <%s>", name, email)
		}
	}

	// Check if name needs encoding (contains non-ASCII characters)
	needsEncoding := false
	for _, r := range name {
		if r > 127 {
			needsEncoding = true
			break
		}
	}

	if needsEncoding {
		// Encode UTF-8 names
		encodedName := base64.StdEncoding.EncodeToString([]byte(name))
		return fmt.Sprintf("=?UTF-8?B?%s?= <%s>", encodedName, email)
	}

	// Use standard format for ASCII names
	return fmt.Sprintf("%s <%s>", name, email)
}

func utf8Split(utf8String string, charset charset) []string {
	var result []string
	buffer := &bytes.Buffer{}
	switch charset {
	case CharsetUTF8:
		for k, c := range strings.Split(utf8String, "") {
			buffer.WriteString(c)
			if (k+1)%charset.Length() == 0 {
				result = append(result, buffer.String())
				buffer.Reset()
			}
		}
	case CharsetISO2022JP:
		for _, c := range strings.Split(utf8String, "") {
			if buffer.Len()+len(c) > charset.Length() {
				result = append(result, buffer.String())
				buffer.Reset()
			}
			buffer.WriteString(c)
		}
	}

	if buffer.Len() > 0 {
		result = append(result, buffer.String())
	}
	return result
}

// encodeSubject encodes a subject line
func (s *MailService) encodeSubject(subject string, charset charset) (string, error) {
	buffer := &bytes.Buffer{}
	splits := utf8Split(subject, charset)
	for i, line := range splits {
		if i > 0 {
			buffer.WriteString(" ")
		}
		buffer.WriteString("=?" + charset.String() + "?B?")
		encodedLine, err := s.encode(charset, line)
		if err != nil {
			return "", err
		}
		buffer.WriteString(encodedLine)
		buffer.WriteString("?=")
		if i != len(splits)-1 {
			buffer.WriteString("\r\n")
		}
	}
	return buffer.String(), nil
}

// splitEncodedString splits an encoded string into parts
func (s *MailService) splitEncodedString(encoded string, maxLength int) []string {
	var parts []string
	for i := 0; i < len(encoded); i += maxLength {
		end := i + maxLength
		if end > len(encoded) {
			end = len(encoded)
		}
		parts = append(parts, encoded[i:end])
	}
	return parts
}

// cutAndAddCrlf cuts text into lines with CRLF
// Direct implementation from send_mail.go
func (s *MailService) cutAndAddCrlf(msg string) string {
	buffer := bytes.Buffer{}
	for k, c := range strings.Split(msg, "") {
		buffer.WriteString(c)
		if (k+1)%MaxLineLength == 0 {
			buffer.WriteString("\r\n")
		}
	}
	return buffer.String()
}

// sendSMTPEmail handles the actual SMTP sending logic
func (s *MailService) sendSMTPEmail(data EmailData, message string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	addr := fmt.Sprintf("%s:%d", s.config.Host, s.config.Port)

	// All recipients
	recipients := []string{data.To}
	recipients = append(recipients, data.Cc...)
	recipients = append(recipients, data.Bcc...)

	// Use direct SMTP API for both local and production environments
	// Connect to the server
	c, err := smtp.Dial(addr)
	if err != nil {
		s.logger.Error("Failed to connect to SMTP server", map[string]interface{}{
			"error": err.Error(),
		})
		return err
	}
	defer c.Close()

	if err = c.Hello("localhost"); err != nil {
		return err
	}

	// Try STARTTLS if configured and supported
	if s.config.UseTLS && !s.isLocalhost() {
		if ok, _ := c.Extension("STARTTLS"); ok {
			serverName, _, _ := net.SplitHostPort(addr)
			tlsConfig := &tls.Config{ServerName: serverName}

			if err = c.StartTLS(tlsConfig); err != nil {
				s.logger.Error("Failed to start TLS", map[string]interface{}{
					"error": err.Error(),
				})
				return err
			}
		} else {
			s.logger.Warn("STARTTLS not supported by server, continuing without TLS", nil)
		}
	}

	// Try authentication if configured
	if s.config.UseAuth {
		if ok, _ := c.Extension("AUTH"); ok {
			auth := smtp.PlainAuth("", s.config.Username, s.config.Password, s.config.Host)
			if err = c.Auth(auth); err != nil {
				// On production, consider this an error
				if !s.isLocalhost() {
					s.logger.Error("Authentication failed", map[string]interface{}{
						"error": err.Error(),
					})
					return err
				}
				// For localhost, just log a warning
				s.logger.Warn("Authentication failed for local server, continuing", map[string]interface{}{
					"error": err.Error(),
				})
			}
		} else {
			s.logger.Warn("AUTH not supported by server", nil)
		}
	}

	// Set sender
	if err = c.Mail(s.config.FromEmail); err != nil {
		s.logger.Error("Failed to set sender", map[string]interface{}{
			"error": err.Error(),
		})
		return err
	}

	// Set all recipients
	for _, recipient := range recipients {
		if err = c.Rcpt(recipient); err != nil {
			s.logger.Error("Failed to add recipient", map[string]interface{}{
				"recipient": recipient,
				"error":     err.Error(),
			})
			return err
		}
	}

	// Send the email data
	wc, err := c.Data()
	if err != nil {
		s.logger.Error("Failed to start data transmission", map[string]interface{}{
			"error": err.Error(),
		})
		return err
	}

	_, err = wc.Write([]byte(message))
	if err != nil {
		s.logger.Error("Failed to write email data", map[string]interface{}{
			"error": err.Error(),
		})
		wc.Close()
		return err
	}

	if err = wc.Close(); err != nil {
		s.logger.Error("Failed to close data writer", map[string]interface{}{
			"error": err.Error(),
		})
		return err
	}

	// Properly close the connection - parse quit error from send_mail.go
	if err = c.Quit(); err != nil {
		s.logger.Warn("Failed to quit SMTP connection cleanly", map[string]interface{}{
			"error": err.Error(),
		})
		// Parse SMTP quit error - from send_mail.go
		err = s.parseQuitError(err)
		if err != nil {
			s.logger.Warn("SMTP quit error", map[string]interface{}{
				"error": err.Error(),
			})
		}
	}

	s.logger.Info("Email sent successfully", map[string]interface{}{
		"to":       data.To,
		"subject":  data.Subject,
		"template": data.TemplateFile,
	})

	return nil
}

// parseQuitError parses SMTP quit command errors - from send_mail.go
func (s *MailService) parseQuitError(err error) error {
	if err == nil {
		return nil
	}
	str := strings.Split(err.Error(), " ")
	if len(str) > 0 {
		if code, err2 := strconv.Atoi(str[0]); err2 == nil {
			if 200 <= code && code < 399 {
				return nil
			} else if 400 <= code && code < 500 {
				return fmt.Errorf("maybe soft bounce: %w", err)
			} else if 500 <= code && code < 600 {
				return fmt.Errorf("maybe hard bounce: %w", err)
			}
		}
	}

	return err
}

// isLocalhost checks if the SMTP server is a local development server
func (s *MailService) isLocalhost() bool {
	return s.config.Host == "localhost" || s.config.Host == "mailhog" || s.config.Host == "127.0.0.1"
}

// canonicalizeEmail normalizes email addresses
func (s *MailService) canonicalizeEmail(email string) string {
	return utils.ConvertKana(email, []string{static.CONV_HALFWIDTHEISU_FROM_FULLWIDTHEISU})
}

// canonicalizeText normalizes text content
func (s *MailService) canonicalizeText(text string) string {
	return utils.ConvertKana(text, []string{static.CONV_FULLWIDTHKATA_FROM_HALFWIDTHKATA_AND_DAKUTEN})
}

// canonicalizeSenderName normalizes sender names
func (s *MailService) canonicalizeSenderName(name string) string {
	if name == "" {
		return name
	}
	return utils.ConvertKana(name, []string{
		static.CONV_FULLWIDTHKATA_FROM_HALFWIDTHKATA_AND_DAKUTEN,
		static.CONV_HALFWIDTHEISU_FROM_FULLWIDTHEISU,
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

	// Try to establish connection
	client, err := smtp.Dial(address)
	if err != nil {
		return fmt.Errorf("failed to connect to SMTP server: %w", err)
	}
	defer client.Close()

	// Test TLS if configured
	if s.config.UseTLS {
		if ok, _ := client.Extension("STARTTLS"); ok {
			serverName, _, _ := net.SplitHostPort(address)
			tlsConfig := &tls.Config{ServerName: serverName}

			if err = client.StartTLS(tlsConfig); err != nil {
				return fmt.Errorf("STARTTLS failed: %w", err)
			}
		} else {
			s.logger.Warn("STARTTLS extension not supported by server", nil)
		}
	}

	// Test authentication if configured
	if s.config.UseAuth {
		if ok, _ := client.Extension("AUTH"); ok {
			auth := smtp.PlainAuth("", s.config.Username, s.config.Password, s.config.Host)
			if err = client.Auth(auth); err != nil {
				return fmt.Errorf("authentication failed: %w", err)
			}
		} else {
			s.logger.Warn("AUTH extension not supported by server", nil)
		}
	}

	s.logger.Info("SMTP connection test successful", map[string]interface{}{
		"host": s.config.Host,
		"port": s.config.Port,
	})

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

// BulkSendEmail sends the same email to multiple recipients
func (s *MailService) BulkSendEmail(recipients []string, data EmailData) (int, error) {
	successCount := 0

	for _, recipient := range recipients {
		emailData := data
		emailData.To = recipient

		err := s.SendEmail(emailData)
		if err != nil {
			s.logger.Error("Failed to send email to recipient", map[string]interface{}{
				"recipient": recipient,
				"error":     err.Error(),
			})
		} else {
			successCount++
		}
	}

	if successCount == 0 && len(recipients) > 0 {
		return 0, errors.New("failed to send email to any recipients")
	}

	return successCount, nil
}
