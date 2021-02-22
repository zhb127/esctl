package delete

import (
	"esctl/internal/app"
	"esctl/pkg/es"
	"esctl/pkg/log"
	"fmt"

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
	IndexName string
	AliasName string
}

func (h *handler) Run(flags *HandlerFlags) error {
	resp, err := h.esHelper.UnaliasIndex(flags.IndexName, []string{flags.AliasName})
	if err != nil {
		return err
	}

	fmt.Printf("%v", resp)
	return nil
}

func (h *handler) ParseCmdFlags(cmdFlags *pflag.FlagSet) (*HandlerFlags, error) {
	handlerFlags := &HandlerFlags{}

	flagIndex, err := cmdFlags.GetString("index")
	if err != nil {
		return nil, err
	}
	handlerFlags.IndexName = flagIndex

	flagAlias, err := cmdFlags.GetString("alias")
	if err != nil {
		return nil, err
	}
	handlerFlags.AliasName = flagAlias

	return handlerFlags, nil
}
