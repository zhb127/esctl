package list

import (
	"esctl/internal/app"
	"esctl/pkg/es"
	"esctl/pkg/log"
	"os"
	"strings"
	"text/template"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/pflag"
)

type IHandler interface {
	Run(flags *HandlerFlags, args []string) error
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
	All    bool
}

func (h *handler) Run(flags *HandlerFlags, indexNameWildcardExps []string) error {
	resp, err := h.esHelper.ListIndices(indexNameWildcardExps...)
	if err != nil {
		return err
	}

	// 剔除隐藏索引
	if !flags.All {
		respNew := &es.ListIndicesResp{}
		for k := range resp.Items {
			item := resp.Items[k]
			if !strings.HasPrefix(item.Name, ".") {
				respNew.Items = append(respNew.Items, item)
			}
		}
		resp = respNew
	}

	if err := h.printf(flags.Format, resp); err != nil {
		return err
	}

	return nil
}

func (h *handler) ParseCmdFlags(cmdFlags *pflag.FlagSet) (*HandlerFlags, error) {
	handlerFlags := &HandlerFlags{}

	format, err := cmdFlags.GetString("format")
	if err != nil {
		return nil, err
	}
	handlerFlags.Format = format

	all, err := cmdFlags.GetBool("all")
	if err != nil {
		return nil, err
	}
	handlerFlags.All = all

	return handlerFlags, nil
}

func (*handler) printf(format string, resp *es.ListIndicesResp) error {
	// 按自定义格式打印
	if format != "" {
		t := template.Must(template.New("custom").Parse(format + "\n"))
		for k := range resp.Items {
			item := resp.Items[k]
			if err := t.Execute(os.Stdout, item); err != nil {
				return err
			}
		}
		return nil
	}

	// 按默认格式打印
	columns := [][]string{}
	for k := range resp.Items {
		item := resp.Items[k]
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
