package batch

import (
	"os"

	"github.com/spf13/cobra"
)

// rootbatch represents the base command when called without any subcommands
var rootBatch = &cobra.Command{
	Use:   "backend-job",
	Short: "backend-job is a job runner for frnc-backend",
	Long:  `backend-job is a job runner for frnc-backend`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootbatch.
func Execute() {
	err := rootBatch.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootbatch.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.frnc-backend.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootbatch.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
