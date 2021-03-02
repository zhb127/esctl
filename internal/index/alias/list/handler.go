package list

import (
	"esctl/internal/app"
	"esctl/pkg/es"
	"esctl/pkg/log"
	"html/template"
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/pflag"
)

type IHandler interface {
	Run(flags *HandlerFlags) error
	ParseCmdFlags(cmdFlags *pflag.FlagSet) (*HandlerFlags, error)
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

func (h *handler) Run(flags *HandlerFlags) error {
	resp, err := h.esHelper.ListAliases()
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

	if format, err := cmdFlags.GetString("format"); err != nil {
		return nil, err
	} else {
		handlerFlags.Format = format
	}

	return handlerFlags, nil
}

func (*handler) printf(format string, resp *es.ListAliasesResp) error {
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
			item.Index,
			item.Alias,
			item.Filter,
			item.RoutingIndex,
			item.RoutingSearch,
		})
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Index", "Alias", "Filter", "RoutingIndex", "RoutingSearch"})
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
