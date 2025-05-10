package batch

import (
	ImportPaypayPayinData "github.com/huydq/test/batch/ImportPaypayPayinData"
	"github.com/spf13/cobra"
)

var readers int
var filesLoadPerStream int
var lineOfDataReadPerStream int

var ImportPaypayPayinDataBatch = &cobra.Command{
	Use:   "import_paypay_payin_data",
	Short: "run import_paypay_payin_data Shell batch job",
	Long:  "run import_paypay_payin_data Shell batch job",
	Run: func(batch *cobra.Command, args []string) {
		ImportPaypayPayinData.Execute(readers, filesLoadPerStream, lineOfDataReadPerStream)
	},
}

func init() {
	ImportPaypayPayinDataBatch.Flags().IntVarP(&readers, "readers", "r", 5, "number of concurrent readers")
	ImportPaypayPayinDataBatch.Flags().IntVarP(&filesLoadPerStream, "fileLoadPerStream", "f", 10, "number of files to load per stream")
	ImportPaypayPayinDataBatch.Flags().IntVarP(&lineOfDataReadPerStream, "lineOfDataReadePerStream", "l", 100, "number of files to load per stream")
	rootBatch.AddCommand(ImportPaypayPayinDataBatch)
}
