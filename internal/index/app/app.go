package app

import (
	"esctl/pkg/es"
	"esctl/pkg/log"
)

type IApp interface {
	LogHelper() log.IHelper
	ESHelper() es.IHelper
}

type AppConfig struct {
	LogHelper log.HelperConfig
	ESHelper  es.HelperConfig
}

type app struct {
	config    AppConfig
	logHelper log.IHelper
	esHelper  es.IHelper
}

func New(config AppConfig) IApp {
	a := &app{
		config: config,
	}
	a.initLogHelper(config.LogHelper)
	a.initESHelper(config.ESHelper)
	return a
}

func (a *app) initLogHelper(config log.HelperConfig) {
	inst, err := log.NewHelper(config)
	if err != nil {
		panic(err)
	}
	a.logHelper = inst
}

func (a *app) initESHelper(config es.HelperConfig) {
	inst, err := es.NewHelper(config, a.logHelper)
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
