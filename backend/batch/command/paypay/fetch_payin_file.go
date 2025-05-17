package command

import (
	"time"

	application "github.com/huydq/test/batch/application/paypay/fetch_payin_file"
	"github.com/spf13/cobra"
)

var workers int
var fileLoadSizePerStream int
var targetDate string

var fetchPayinFile = &cobra.Command{
	Use:   "paypay_fetch_payin_file",
	Short: "run paypay_fetch_payin_file Shell batch job",
	Long:  "run paypay_fetch_payin_file Shell batch job for fetching paypay payin data from Paypay Server to Storage",
	Run: func(batch *cobra.Command, args []string) {
		application.Execute(workers, fileLoadSizePerStream, targetDate)
	},
}

func InitUploadPaypayCSVToS3Batch(rootBatch *cobra.Command) {
	fetchPayinFile.Flags().IntVarP(&workers, "workers", "w", 5, "number of concurrent workers")
	fetchPayinFile.Flags().IntVarP(&fileLoadSizePerStream, "fileLoadSizePerStream", "f", 10, "number of files to load per stream")
	fetchPayinFile.Flags().StringVarP(&targetDate, "targetDate", "t", time.Now().Format("20060102"), "target date in format: yyyymmdd")

	rootBatch.AddCommand(fetchPayinFile)
}
