package app

import (
	"esctl/pkg/es"
	"esctl/pkg/log"
)

type IApp interface {
	LogHelper() log.IHelper
	ESHelper() es.IHelper
}

type app struct {
	config    *appConfig
	logHelper log.IHelper
	esHelper  es.IHelper
}

func New() IApp {
	a := &app{
		config: newAppConfig(),
	}

	a.initLogHelper()
	a.initESHelper()

	return a
}

func (a *app) initLogHelper() {
	inst, err := log.NewHelper(a.config.LogHelper)
	if err != nil {
		panic(err)
	}
	a.logHelper = inst
}

func (a *app) initESHelper() {
	inst, err := es.NewHelper(a.config.ESHelper, a.logHelper)
	if err != nil {
		panic(err)
	}
	a.esHelper = inst
}

func (a *app) LogHelper() log.IHelper {
	return a.logHelper
}

func (a *app) ESHelper() es.IHelper {
	return a.esHelper
}
