package delete

import (
	"esctl/internal/index/app"
	"esctl/pkg/es"
	"esctl/pkg/log"
	"fmt"
)

type IHandler interface {
	Run(indexNames []string) error
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

func (h *handler) Run(indexNames []string) error {
	if _, err := h.esHelper.DeleteIndices(indexNames...); err != nil {
		return err
	}

	fmt.Println("success")
	return nil
}
