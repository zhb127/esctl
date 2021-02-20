package create

import (
	"errors"
	"esctl/internal/index/app"
	"esctl/pkg/es"
	"esctl/pkg/log"
	"fmt"
	"io/ioutil"
	"os"

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
	IndexBody []byte
}

func (h *handler) Handle(flags *HandlerFlags) error {
	indexName := flags.IndexName
	indexBody := flags.IndexBody
	resp, err := h.esHelper.CreateIndex(indexName, indexBody)
	if err != nil {
		return err
	}

	fmt.Printf("%v", resp)
	return nil
}

func (h *handler) ParseCmdFlags(cmdFlags *pflag.FlagSet) (*HandlerFlags, error) {
	handlerFlags := &HandlerFlags{}

	// 处理 --name
	flagName, err := cmdFlags.GetString("name")
	if err != nil {
		return nil, err
	}
	handlerFlags.IndexName = flagName

	// 处理 --body
	flagBody, err := cmdFlags.GetString("body")
	if err != nil {
		return nil, err
	}
	if flagBody != "" {
		handlerFlags.IndexBody = []byte(flagBody)
	}

	// 处理 --file
	if flagBody == "" {
		flagFile, err := cmdFlags.GetString("file")
		if err != nil {
			return nil, err
		}
		if flagFile == "" {
			return nil, errors.New("oneof --body, --file is required")
		}

		bodyFile, err := os.Open(flagFile)
		if err != nil {
			return nil, err
		}
		defer bodyFile.Close()

		bodyBytes, _ := ioutil.ReadAll(bodyFile)
		handlerFlags.IndexBody = bodyBytes
	}

	return handlerFlags, nil
}
