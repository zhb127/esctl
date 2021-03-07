package reindex

import (
	"encoding/json"
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
	SrcIndexName  string
	DestIndexName string
}

func (h *handler) Run(flags *HandlerFlags) error {
	resp, err := h.esHelper.Reindex(flags.SrcIndexName, flags.DestIndexName)
	if err != nil {
		return err
	}

	jsonBytes, _ := json.Marshal(resp)
	fmt.Println(string(jsonBytes))

	return nil
}

func (h *handler) ParseCmdFlags(cmdFlags *pflag.FlagSet) (*HandlerFlags, error) {
	handlerFlags := &HandlerFlags{}

	src, err := cmdFlags.GetString("src")
	if err != nil {
		return nil, err
	}
	handlerFlags.SrcIndexName = src

	dest, err := cmdFlags.GetString("dest")
	if err != nil {
		return nil, err
	}
	handlerFlags.DestIndexName = dest

	return handlerFlags, nil
}
