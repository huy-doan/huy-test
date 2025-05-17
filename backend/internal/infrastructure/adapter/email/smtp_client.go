package email

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"net/smtp"
	"strconv"
	"strings"
	"sync"

	"github.com/huydq/test/internal/pkg/logger"
)

// SMTPClient handles communication with the SMTP server
type SMTPClient struct {
	config *MailConfig
	logger logger.Logger
	mutex  sync.Mutex
}

// NewSMTPClient creates a new instance of SMTPClient
func NewSMTPClient(config *MailConfig, logger logger.Logger) *SMTPClient {
	return &SMTPClient{
		config: config,
		logger: logger,
	}
}

// SendRawMessage sends a raw email message to the SMTP server
func (c *SMTPClient) SendRawMessage(from string, to, cc, bcc []string, message []byte) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	// Create server address
	addr := fmt.Sprintf("%s:%d", c.config.Host, c.config.Port)

	// Connect to the server
	client, err := smtp.Dial(addr)
	if err != nil {
		return fmt.Errorf("failed to connect to SMTP server: %w", err)
	}
	// nolint
	defer client.Close()

	// Send HELO command
	if err = client.Hello("localhost"); err != nil {
		return fmt.Errorf("HELO command failed: %w", err)
	}

	// Set up TLS if configured
	if c.config.UseTLS {
		if ok, _ := client.Extension("STARTTLS"); !ok {
			return errors.New("smtp: server doesn't support STARTTLS")
		}

		serverName, _, err := net.SplitHostPort(addr)
		if err != nil {
			return fmt.Errorf("failed to parse server name: %w", err)
		}

		tlsConfig := &tls.Config{ServerName: serverName}
		if err = client.StartTLS(tlsConfig); err != nil {
			c.logger.Error("Failed to start TLS", map[string]any{
				"error": err.Error(),
			})
			return fmt.Errorf("TLS setup failed: %w", err)
		}
	}

	// Authenticate if required
	if c.config.UseAuth {
		if ok, _ := client.Extension("AUTH"); !ok {
			return errors.New("smtp: server doesn't support AUTH")
		}

		auth := smtp.PlainAuth("", c.config.Username, c.config.Password, c.config.Host)
		if err = client.Auth(auth); err != nil {
			return fmt.Errorf("authentication failed: %w", err)
		}
	}

	// Set sender
	if err = client.Mail(from); err != nil {
		c.logger.Error("Failed to set sender", map[string]any{
			"error": err.Error(),
		})
		return fmt.Errorf("failed to set sender: %w", err)
	}

	// Set recipients
	for _, recipient := range to {
		if err := client.Rcpt(recipient); err != nil {
			return fmt.Errorf("failed to add recipient (%s): %w", recipient, err)
		}
	}

	for _, recipient := range cc {
		if err := client.Rcpt(recipient); err != nil {
			return fmt.Errorf("failed to add CC recipient (%s): %w", recipient, err)
		}
	}

	for _, recipient := range bcc {
		if err := client.Rcpt(recipient); err != nil {
			return fmt.Errorf("failed to add BCC recipient (%s): %w", recipient, err)
		}
	}

	// Send the message
	wc, err := client.Data()
	if err != nil {
		return fmt.Errorf("failed to start data transaction: %w", err)
	}

	defer func() {
		if err := c.parseQuitError(wc.Close()); err != nil {
			c.logger.Error("Failed to close data writer", map[string]any{
				"error": err.Error(),
			})
		}
	}()

	if _, err := wc.Write(message); err != nil {
		return fmt.Errorf("failed to write message: %w", err)
	}

	// Quit the session
	if err = c.parseQuitError(client.Quit()); err != nil {
		c.logger.Warn("SMTP quit error", map[string]any{
			"error": err.Error(),
		})
	}

	return nil
}

// TestConnection tests the SMTP connection
func (c *SMTPClient) TestConnection() error {
	address := fmt.Sprintf("%s:%d", c.config.Host, c.config.Port)

	// Try to establish connection
	client, err := smtp.Dial(address)
	if err != nil {
		return fmt.Errorf("failed to connect to SMTP server: %w", err)
	}
	// nolint
	defer client.Close()

	// Test TLS if configured
	if c.config.UseTLS {
		if ok, _ := client.Extension("STARTTLS"); ok {
			serverName, _, _ := net.SplitHostPort(address)
			tlsConfig := &tls.Config{ServerName: serverName}

			if err = client.StartTLS(tlsConfig); err != nil {
				return fmt.Errorf("STARTTLS failed: %w", err)
			}
		} else {
			c.logger.Warn("STARTTLS extension not supported by server", nil)
		}
	}

	// Test authentication if configured
	if c.config.UseAuth {
		if ok, _ := client.Extension("AUTH"); ok {
			auth := smtp.PlainAuth("", c.config.Username, c.config.Password, c.config.Host)
			if err = client.Auth(auth); err != nil {
				return fmt.Errorf("authentication failed: %w", err)
			}
		} else {
			c.logger.Warn("AUTH extension not supported by server", nil)
		}
	}

	c.logger.Info("SMTP connection test successful", map[string]any{
		"host": c.config.Host,
		"port": c.config.Port,
	})

	return nil
}

// parseQuitError analyzes SMTP error codes to determine the error type
func (c *SMTPClient) parseQuitError(err error) error {
	if err == nil {
		return nil
	}

	str := strings.Split(err.Error(), " ")
	if len(str) > 0 {
		if code, err2 := strconv.Atoi(str[0]); err2 == nil {
			if 200 <= code && code < 399 {
				return nil // Success codes
			} else if 400 <= code && code < 500 {
				return fmt.Errorf("maybe soft bounce: %w", err) // Temporary failure
			} else if 500 <= code && code < 600 {
				return fmt.Errorf("maybe hard bounce: %w", err) // Permanent failure
			}
		}
	}

	return err
}

// IsLocalhost checks if the SMTP server is a local development server
func (c *SMTPClient) IsLocalhost() bool {
	return c.config.Host == "localhost" ||
		c.config.Host == "mailhog" ||
		c.config.Host == "127.0.0.1"
}
