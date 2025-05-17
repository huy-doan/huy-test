package email

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"net/mail"
	"strings"
	"time"

	"github.com/huydq/test/internal/pkg/utils"
	"github.com/huydq/test/internal/pkg/utils/static"
)

// MessageComposer handles email message preparation and composition
type MessageComposer struct {
	config *MailConfig
}

// NewMessageComposer creates a new instance of MessageComposer
func NewMessageComposer(config *MailConfig) *MessageComposer {
	return &MessageComposer{
		config: config,
	}
}

// ComposeMessage prepares a complete email message with headers and content
func (c *MessageComposer) ComposeMessage(from string, data EmailData, content string) ([]byte, error) {
	// Encode the content in base64
	encodedContent := base64.StdEncoding.EncodeToString([]byte(content))

	// Encode the subject
	subject, err := utils.EncodeSubject(data.Subject, data.Charset)
	if err != nil {
		return nil, fmt.Errorf("failed to encode subject: %w", err)
	}

	// Create message buffer
	var buf bytes.Buffer

	// Add headers
	c.writeHeaders(&buf, from, data.To, data.Cc, data.Bcc, subject, data.ContentType, data.Charset)

	// Add empty line to separate headers from body
	buf.WriteString("\r\n")

	// Add content with proper line wrapping
	buf.WriteString(utils.WrapBase64Lines(encodedContent))

	return buf.Bytes(), nil
}

// writeHeaders writes all the headers to the message buffer
func (c *MessageComposer) writeHeaders(buf *bytes.Buffer, from string, to, cc, bcc []string,
	encodedSubject, contentType string, charset utils.EmailCharset) {
	// From header
	fmt.Fprintf(buf, "From: %s\r\n", from)

	// To header
	toAddr := strings.Join(to, ",")
	if toAddr != "" {
		fmt.Fprintf(buf, "To: %s\r\n", toAddr)
	}

	// CC header
	ccAddr := strings.Join(cc, ",")
	if ccAddr != "" {
		fmt.Fprintf(buf, "Cc: %s\r\n", ccAddr)
	}

	// BCC header
	bccAddr := strings.Join(bcc, ",")
	if bccAddr != "" {
		fmt.Fprintf(buf, "Bcc: %s\r\n", bccAddr)
	}

	// Subject header
	fmt.Fprintf(buf, "Subject: %s\r\n", encodedSubject)

	// Date header - using RFC1123Z format for email compatibility
	fmt.Fprintf(buf, "Date: %s\r\n", time.Now().Format(time.RFC1123Z))

	// MIME headers
	buf.WriteString("MIME-Version: 1.0\r\n")
	fmt.Fprintf(buf, "Content-Type: %s;charset=\"%s\"\r\n", contentType, charset.String())
	buf.WriteString("Content-Transfer-Encoding: base64\r\n")
}

// FormatFromAddress formats a sender name and address into a valid From header
func (c *MessageComposer) FormatFromAddress(name, address string, charset utils.EmailCharset) (string, error) {
	// If no name, just return the address
	if name == "" {
		return address, nil
	}

	// For UTF-8, use Go's standard mail package
	if charset == utils.CharsetUTF8 {
		return (&mail.Address{
			Name:    name,
			Address: address,
		}).String(), nil
	}

	// For ISO-2022-JP, we need to encode the name separately
	if charset == utils.CharsetISO2022JP {
		encodedName, err := utils.EncodeMIMEHeader(name, charset)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("%s <%s>", encodedName, address), nil
	}

	return address, nil
}

// CanonicalizeEmail normalizes email addresses
func (c *MessageComposer) CanonicalizeEmail(email string) string {
	return utils.ConvertKana(email, []string{static.CONV_HALFWIDTHEISU_FROM_FULLWIDTHEISU})
}

// CanonicalizeText normalizes text content
func (c *MessageComposer) CanonicalizeText(text string) string {
	return utils.ConvertKana(text, []string{static.CONV_FULLWIDTHKATA_FROM_HALFWIDTHKATA_AND_DAKUTEN})
}

// CanonicalizeSenderName normalizes sender names
func (c *MessageComposer) CanonicalizeSenderName(name string) string {
	if name == "" {
		return name
	}
	return utils.ConvertKana(name, []string{
		static.CONV_FULLWIDTHKATA_FROM_HALFWIDTHKATA_AND_DAKUTEN,
		static.CONV_HALFWIDTHEISU_FROM_FULLWIDTHEISU,
	})
}
