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
		tablewriter.Colors{tablewriter.Bold},
		tablewriter.Colors{tablewriter.Bold})
	return t
}

func SuccessStyle(t *tablewriter.Table, data []string) *tablewriter.Table {
	mergeable := tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiWhiteColor}

	if data[len(data)-1] == "FALSE" {
		mergeable = tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiRedColor}
	} else {
		mergeable = tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiGreenColor}
	}
	t.Rich(data, []tablewriter.Colors{
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiCyanColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiGreenColor},
		tablewriter.Colors{tablewriter.Bold},
		tablewriter.Colors{tablewriter.Bold},
		tablewriter.Colors{tablewriter.Bold},
		mergeable,
	})
	return t
}

func ErrorStyle(t *tablewriter.Table, data []string) *tablewriter.Table {
	t.Rich(data, []tablewriter.Colors{
		tablewriter.Colors{tablewriter.Bold,
			tablewriter.FgHiRedColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiRedColor}})
	return t
}
