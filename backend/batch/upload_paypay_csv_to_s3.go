package batch

import (
	UploadPaypayCSVToS3 "github.com/huydq/test/batch/UploadPaypayCSVToS3"
	"github.com/spf13/cobra"
)

var workers int
var fileLoadSizePerStream int

var UploadPaypayCSVToS3Batch = &cobra.Command{
	Use:   "upload_paypay_csv_to_s3",
	Short: "run upload_paypay_csv_to_s3 Shell batch job",
	Long:  "run upload_paypay_csv_to_s3 Shell batch job",
	Run: func(batch *cobra.Command, args []string) {
		UploadPaypayCSVToS3.Execute(workers, fileLoadSizePerStream)
	},
}

func init() {
	UploadPaypayCSVToS3Batch.Flags().IntVarP(&workers, "workers", "w", 5, "number of concurrent workers")
	UploadPaypayCSVToS3Batch.Flags().IntVarP(&fileLoadSizePerStream, "fileLoadSizePerStream", "f", 10, "number of files to load per stream")
	rootBatch.AddCommand(UploadPaypayCSVToS3Batch)
}
