package version

import (
	"esctl/internal/version/app"
	"esctl/pkg/es"
	"esctl/pkg/log"
)

type IHandler interface {
	Run() error
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

func (h *handler) Run() error {
	return nil
}
