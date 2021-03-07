package create

import (
	"errors"
	"esctl/internal/app"
	"esctl/pkg/es"
	"esctl/pkg/log"
	"io/ioutil"
	"os"

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
	Name string
	Body []byte
}

func (h *handler) Run(flags *HandlerFlags) error {
	resp, err := h.esHelper.CreateIndex(flags.Name, flags.Body)
	if err != nil {
		return err
	}

	h.logHelper.Info("success", map[string]interface{}{
		"result": resp,
	})

	return nil
}

func (h *handler) ParseCmdFlags(cmdFlags *pflag.FlagSet) (*HandlerFlags, error) {
	handlerFlags := &HandlerFlags{}

	name, err := cmdFlags.GetString("name")
	if err != nil {
		return nil, err
	}
	handlerFlags.Name = name

	body, err := cmdFlags.GetString("body")
	if err != nil {
		return nil, err
	}
	if body != "" {
		handlerFlags.Body = []byte(body)
	}

	if handlerFlags.Body == nil {
		file, err := cmdFlags.GetString("file")
		if err != nil {
			return nil, err
		}

		if file == "" {
			return nil, errors.New("oneof --body, --file is required")
		}

		fd, err := os.Open(file)
		if err != nil {
			return nil, err
		}
		defer fd.Close()

		byteArr, _ := ioutil.ReadAll(fd)
		handlerFlags.Body = byteArr
	}

	return handlerFlags, nil
}
