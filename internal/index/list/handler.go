package list

import (
	"esctl/internal/index/app"
	"esctl/pkg/es"
	"esctl/pkg/log"
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/pflag"
)

type IHandler interface {
	Handle(flags *pflag.FlagSet, args []string) error
}

type handler struct {
	logHelper log.IHelper
	esHelper  es.IHelper
}

func NewHandler(a app.IApp) IHandler {
	return &handler{
		logHelper: a.LogHelper(),
		esHelper:  a.ESHelper(),
	}
}

func (h *handler) Handle(flags *pflag.FlagSet, args []string) error {
	resp, err := h.esHelper.CatIndices(args...)
	if err != nil {
		return err
	}

	columns := [][]string{}
	for k := range resp.Items {
		item := resp.Items[k]
		columns = append(columns, []string{
			item.Health,
			item.Status,
			item.Index,
			item.Uuid,
			item.Pri,
			item.Rep,
			item.DocsCount,
			item.DocsDeleted,
			item.StoreSize,
			item.PriStoreSize,
		})
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Health", "Status", "Index", "UUID", "Pri", "Rep", "DocsCount", "DocsDeleted", "StoreSize", "PriStoreSize"})
	table.SetAutoFormatHeaders(false)
	table.SetAutoWrapText(false)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	table.SetRowSeparator("")
	table.SetHeaderLine(false)
	table.SetBorder(false)
	table.SetTablePadding("\t")
	table.SetNoWhiteSpace(true)
	table.AppendBulk(columns)
	table.Render()

	return nil
}
