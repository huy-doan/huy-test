package command

import (
	application "github.com/huydq/test/batch/application/paypay/import_payin_file"
	"github.com/spf13/cobra"
)

var readers int
var filesLoadPerStream int
var lineOfDataReadPerStream int

var importPaypayPayinFile = &cobra.Command{
	Use:   "paypay_import_payin_file",
	Short: "run paypay_import_payin_file Shell batch job",
	Long:  "run paypay_import_payin_file Shell batch job for importing paypay payin data from Zip file to DB",
	Run: func(batch *cobra.Command, args []string) {
		application.Execute(readers, filesLoadPerStream, lineOfDataReadPerStream)
	},
}

func InitImportPaypayPayinDataBatch(rootBatch *cobra.Command) {
	importPaypayPayinFile.Flags().IntVarP(&readers, "readers", "r", 5, "number of concurrent readers")
	importPaypayPayinFile.Flags().IntVarP(&filesLoadPerStream, "fileLoadPerStream", "f", 10, "number of files to load per stream")
	importPaypayPayinFile.Flags().IntVarP(&lineOfDataReadPerStream, "lineOfDataReadePerStream", "l", 100, "number of files to load per stream")

	rootBatch.AddCommand(importPaypayPayinFile)
}
