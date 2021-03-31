package delete

import (
	"esctl/internal/app"
	"esctl/pkg/es"
	"esctl/pkg/log"

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
	Index string
	Query string
}

func (h *handler) Run(flags *HandlerFlags) error {
	if err := h.esHelper.DeleteDocsByQuery(flags.Index, []byte(flags.Query)); err != nil {
		return err
	}
	return nil
}

func (h *handler) ParseCmdFlags(cmdFlags *pflag.FlagSet) (*HandlerFlags, error) {
	handlerFlags := &HandlerFlags{}

	index, err := cmdFlags.GetString("index")
	if err != nil {
		return nil, err
	}
	handlerFlags.Index = index

	query, err := cmdFlags.GetString("query")
	if err != nil {
		return nil, err
	}
	handlerFlags.Query = query

	return handlerFlags, nil
}
