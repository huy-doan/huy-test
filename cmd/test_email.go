package cmd

import (
	"github.com/spf13/cobra"
	"github.com/vnlab/makeshop-payment/cmd/TestEmail"
)

// testEmailCmd represents the testEmail command
var testEmailCmd = &cobra.Command{
	Use:   "test-email",
	Short: "Test email sending functionality",
	Long:  "Test email sending functionality by sending test emails to verify SMTP configuration",
	Run: func(cmd *cobra.Command, args []string) {
		TestEmail.Execute()
	},
}

func init() {
	rootCmd.AddCommand(testEmailCmd)
}
