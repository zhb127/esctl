package delete

import (
	"esctl/internal/index/app"
	"esctl/pkg/es"
	"esctl/pkg/log"
	"fmt"

	"github.com/spf13/pflag"
)

type IHandler interface {
	Handle(flags *pflag.FlagSet, args []string) error
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

func (h *handler) Handle(flags *pflag.FlagSet, args []string) error {
	if _, err := h.esHelper.DeleteIndices(args...); err != nil {
		return err
	}

	fmt.Println("success")
	return nil
}
