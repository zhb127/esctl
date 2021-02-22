package up

import (
	"esctl/internal/app"
	"esctl/pkg/es"
	"esctl/pkg/log"
	"fmt"
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
	Dir string
}

func (h *handler) Run(flags *HandlerFlags) error {
	f, err := os.Open(flags.Dir)
	if err != nil {
		return err
	}
	defer f.Close()

	files, err := f.Readdir(-1)
	if err != nil {
		return nil
	}

	for _, file := range files {
		fmt.Println(file.Name())
	}

	return nil
}

func (h *handler) ParseCmdFlags(cmdFlags *pflag.FlagSet) (*HandlerFlags, error) {
	handlerFlags := &HandlerFlags{}

	flagDir, err := cmdFlags.GetString("dir")
	if err != nil {
		return nil, err
	}
	handlerFlags.Dir = flagDir

	return handlerFlags, nil
}
