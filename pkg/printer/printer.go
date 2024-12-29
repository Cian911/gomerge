package printer

import (
	"io"

	"github.com/olekukonko/tablewriter"
)

func NewTable(w io.Writer, headers []string) (t *tablewriter.Table) {
	t = tablewriter.NewWriter(w)
	t.SetHeader(headers)

	return
}

func HeaderStyle(t *tablewriter.Table) *tablewriter.Table {
	t.SetHeaderColor(tablewriter.Colors{tablewriter.Bold},
		tablewriter.Colors{tablewriter.Bold},
		tablewriter.Colors{tablewriter.Bold},
		tablewriter.Colors{tablewriter.Bold},
		tablewriter.Colors{tablewriter.Bold})
	return t
}

func SuccessStyle(t *tablewriter.Table, data []string) *tablewriter.Table {
	t.Rich(data, []tablewriter.Colors{
		{tablewriter.Bold, tablewriter.FgHiCyanColor},
		{tablewriter.Bold, tablewriter.FgHiGreenColor},
		{tablewriter.Bold},
		{tablewriter.Bold},
		{tablewriter.Bold},
	})
	return t
}

func ErrorStyle(t *tablewriter.Table, data []string) *tablewriter.Table {
	t.Rich(
		data,
		[]tablewriter.Colors{
			{tablewriter.Bold, tablewriter.FgHiCyanColor},
			{tablewriter.Bold, tablewriter.FgHiRedColor},
		},
	)
	return t
}

func WaitingStyle(t *tablewriter.Table, data []string) *tablewriter.Table {
	t.Rich(
		data,
		[]tablewriter.Colors{
			{tablewriter.Bold, tablewriter.FgHiCyanColor},
			{tablewriter.Bold, tablewriter.FgHiYellowColor},
		},
	)
	return t
}
