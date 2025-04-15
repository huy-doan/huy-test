// cmd/TestEmail/main.go

package TestEmail

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/vnlab/makeshop-payment/src/infrastructure/email"
	"github.com/vnlab/makeshop-payment/src/infrastructure/logger"
)

// Execute runs the email test utility
func Execute() {
	log.Println("======= Starting Email Test =======")
	defer log.Println("======= Email Test Complete =======")

	appLogger := logger.GetLogger()

	if os.Getenv("SMTP_HOST") == "" {
		log.Fatal("Email configuration is incomplete. Please set SMTP_HOST and other required environment variables.")
	}

	// Initialize mail service
	mailService, err := email.NewMailService(appLogger)
	if err != nil {
		log.Fatalf("Failed to initialize mail service: %v", err)
	}

	if err := mailService.TestConnection(); err != nil {
		log.Fatalf("SMTP connection test failed: %v", err)
	}
	log.Println("SMTP connection test succeeded")

	log.Println("Sending password changed email...")
	
	emailData := email.EmailData{
		To:             "test@example.com",
		ToName:         "Test User",
		Subject:        "パスワード変更完了のお知らせ",
		TemplateFile:   "password_changed.tmpl",
		TemplateFolder: "auth",
		Data: map[string]interface{}{
			"UserID":    1,
			"UserEmail": "test@example.com",
			"UserName":  "Test User",
			"ToName":    "Test User",
			"Timestamp": time.Now().Format("2006年01月02日 15:04"),
		},
	}
	
	if err := mailService.SendEmail(emailData); err != nil {
		log.Printf("Failed to send password changed email: %v", err)
	} else {
		log.Println("Password changed email sent successfully")
	}

	fmt.Println("\nEmail test completed. Check MailHog UI at http://localhost:8025 to view the sent email.")
}
