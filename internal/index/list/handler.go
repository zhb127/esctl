package list

import (
	"esctl/internal/index/app"
	"esctl/pkg/es"
	"esctl/pkg/log"
	"os"
	"text/template"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/pflag"
)

type IHandler interface {
	ParseCmdFlags(cmdFlags *pflag.FlagSet) (*HandlerFlags, error)
	Handle(flags *HandlerFlags, args []string) error
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

type HandlerFlags struct {
	Format string
}

func (h *handler) Handle(flags *HandlerFlags, indexNameWithWildcards []string) error {
	resp, err := h.esHelper.CatIndices(indexNameWithWildcards...)
	if err != nil {
		return err
	}

	if err := h.printf(flags.Format, resp); err != nil {
		return err
	}

	return nil
}

func (h *handler) ParseCmdFlags(cmdFlags *pflag.FlagSet) (*HandlerFlags, error) {
	handlerFlags := &HandlerFlags{}

	flagFormat, err := cmdFlags.GetString("format")
	if err != nil {
		return nil, err
	}
	handlerFlags.Format = flagFormat

	return handlerFlags, nil
}

func (*handler) printf(format string, resp *es.CatIndicesResp) error {
	// 按指定格式打印
	if format != "" {
		t := template.Must(template.New("specifiedFormat").Parse(format + "\n"))
		for _, item := range resp.Items {
			if err := t.Execute(os.Stdout, item); err != nil {
				return err
			}
		}
		return nil
	}

	// 按默认格式打印
	columns := [][]string{}
	for _, item := range resp.Items {
		columns = append(columns, []string{
			item.ID,
			item.Name,
			item.Health,
			item.Status,
			item.Pri,
			item.Rep,
			item.DocsCount,
			item.DocsDeleted,
			item.StoreSize,
			item.PriStoreSize,
		})
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "Name", "Health", "Status", "Pri", "Rep", "DocsCount", "DocsDeleted", "StoreSize", "PriStoreSize"})
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
