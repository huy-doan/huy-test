package cmd

import (
	"github.com/spf13/cobra"
	"github.com/vnlab/makeshop-payment/cmd/Server2S3Shell"
)

var workers int
var fileLoadSizePerStream int

var server2S3Shell = &cobra.Command{
	Use:   "sync_paypay_csv_to_s3",
	Short: "run sync_paypay_csv_to_s3 Shell batch job",
	Long:  "run sync_paypay_csv_to_s3 Shell batch job",
	Run: func(cmd *cobra.Command, args []string) {
		Server2S3Shell.Execute(workers, fileLoadSizePerStream)
	},
}

func init() {
	server2S3Shell.Flags().IntVarP(&workers, "workers", "w", 5, "number of concurrent workers")
	server2S3Shell.Flags().IntVarP(&fileLoadSizePerStream, "fileLoadSizePerStream", "f", 100, "number of files to load per stream")
	rootCmd.AddCommand(server2S3Shell)
}
