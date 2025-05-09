package cmd

import (
	"github.com/huydq/test/cmd/ExampleShell"
	"github.com/spf13/cobra"
)

// exampleShell represents the ExampleShell command.
// To run this command on local, use the following command:
// $ make shell "go run main.go example"
var exampleShell = &cobra.Command{
	Use:   "example",
	Short: "run exampleShell batch job",
	Long:  "run exampleShell batch job",
	Run: func(cmd *cobra.Command, args []string) {
		ExampleShell.Execute()
	},
}

func init() {
	rootCmd.AddCommand(exampleShell)
}
