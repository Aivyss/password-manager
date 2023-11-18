package view

import (
	"github.com/olekukonko/tablewriter"
	"os"
)

func StdoutTableView(headerCsvLine []string, onlyDataCsvLines [][]string) {
	viewTable := tablewriter.NewWriter(os.Stdout)
	viewTable.SetHeader(headerCsvLine)
	viewTable.SetBorder(false)
	viewTable.AppendBulk(onlyDataCsvLines)
	viewTable.Render()
}
