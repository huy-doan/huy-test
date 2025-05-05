package cmd

import (
	"github.com/spf13/cobra"
	Csv2MysqlShell "github.com/vnlab/makeshop-payment/cmd/Csv2MysqlShell"
)

var readers int
var filesLoadPerStream int
var lineOfDataReadePerStream int

var csv2MysqlShell = &cobra.Command{
	Use:   "sync_paypay_csv_to_mysql",
	Short: "run sync_paypay_csv_to_mysql Shell batch job",
	Long:  "run sync_paypay_csv_to_mysql Shell batch job",
	Run: func(cmd *cobra.Command, args []string) {
		Csv2MysqlShell.Execute(readers, filesLoadPerStream, lineOfDataReadePerStream)
	},
}

func init() {
	csv2MysqlShell.Flags().IntVarP(&readers, "readers", "r", 5, "number of concurrent readers")
	csv2MysqlShell.Flags().IntVarP(&filesLoadPerStream, "fileLoadPerStream", "f", 10, "number of files to load per stream")
	csv2MysqlShell.Flags().IntVarP(&lineOfDataReadePerStream, "lineOfDataReadePerStream", "l", 100, "number of files to load per stream")
	rootCmd.AddCommand(csv2MysqlShell)
}
