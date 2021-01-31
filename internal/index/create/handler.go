package create

import (
	"esctl/internal/index/app"
	"esctl/pkg/es"
	"esctl/pkg/log"
	"fmt"

	"github.com/spf13/pflag"
)

type IHandler interface {
	ParseCmdFlags(cmdFlags *pflag.FlagSet) (*HandlerFlags, error)
	Handle(flags *HandlerFlags) error
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
	IndexBody string
}

func (h *handler) Handle(flags *HandlerFlags) error {
	indexName := flags.IndexName
	indexBody := []byte(flags.IndexBody)
	resp, err := h.esHelper.CreateIndex(indexName, indexBody)
	if err != nil {
		return err
	}

	fmt.Printf("%v", resp)
	return nil
}

func (h *handler) ParseCmdFlags(cmdFlags *pflag.FlagSet) (*HandlerFlags, error) {
	handlerFlags := &HandlerFlags{}

	flagName, err := cmdFlags.GetString("name")
	if err != nil {
		return nil, err
	}
	handlerFlags.IndexName = flagName

	flagBody, err := cmdFlags.GetString("body")
	if err != nil {
		return nil, err
	}
	handlerFlags.IndexBody = flagBody

	return handlerFlags, nil
}
