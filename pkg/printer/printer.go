package printer

import (
	"io"

	"github.com/olekukonko/tablewriter"
)

func NewTable(w io.Writer, headers []string) (t *tablewriter.Table) {
	t = tablewriter.NewWriter(w)
	t.SetHeader(headers)
	applyStyle(t)

	return
}

func SuccessStyle(t *tablewriter.Table, data []string) *tablewriter.Table {
	t.Rich(data, []tablewriter.Colors{tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiGreenColor}, tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiGreenColor}})
	return t
}

func ErrorStyle(t *tablewriter.Table, data []string) *tablewriter.Table {
	t.Rich(data, []tablewriter.Colors{tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiRedColor}, tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiRedColor}})
	return t
}

func applyStyle(t *tablewriter.Table) {
	t.SetAutoWrapText(false)
	t.SetAutoFormatHeaders(true)
	t.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	t.SetAlignment(tablewriter.ALIGN_LEFT)
	t.SetCenterSeparator("")
	t.SetColumnSeparator("")
	t.SetRowSeparator("")
	t.SetHeaderLine(false)
	t.SetBorder(false)
	t.SetTablePadding("\t")
	t.SetNoWhiteSpace(true)
}
