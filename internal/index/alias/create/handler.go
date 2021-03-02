package create

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
	Index string
	Alias string
}

func (h *handler) Run(flags *HandlerFlags) error {
	resp, err := h.esHelper.AliasIndex(flags.Index, flags.Alias)
	if err != nil {
		return err
	}

	fmt.Printf("%v", resp)
	return nil
}

func (h *handler) ParseCmdFlags(cmdFlags *pflag.FlagSet) (*HandlerFlags, error) {
	handlerFlags := &HandlerFlags{}

	if index, err := cmdFlags.GetString("index"); err != nil {
		return nil, err
	} else {
		handlerFlags.Index = index
	}

	if alias, err := cmdFlags.GetString("alias"); err != nil {
		return nil, err
	} else {
		handlerFlags.Alias = alias
	}

	return handlerFlags, nil
}
