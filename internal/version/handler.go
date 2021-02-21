package version

import (
	"bytes"
	"encoding/json"
	"esctl/internal/version/app"
	"esctl/pkg/es"
	"esctl/pkg/log"
	"fmt"
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
	resp, err := h.esHelper.Info()
	if err != nil {
		return err
	}

	jsonBytes, _ := json.Marshal(resp)

	var resBuff bytes.Buffer
	if err := json.Indent(&resBuff, jsonBytes, "", "  "); err != nil {
		return err
	}

	fmt.Printf("%s\n", resBuff.Bytes())
	return nil
}
