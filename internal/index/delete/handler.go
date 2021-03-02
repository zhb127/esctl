package delete

import (
	"esctl/internal/app"
	"esctl/pkg/es"
	"esctl/pkg/log"
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
	_, err := h.esHelper.DeleteIndices(indexNames...)
	if err != nil {
		return err
	}

	return nil
}
