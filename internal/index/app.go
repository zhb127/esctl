package index

import (
	"esctl/pkg/config/dotenv"
	"esctl/pkg/es"
	"esctl/pkg/log"
)

type IApp interface {
	LogHelper() log.IHelper
	ESHelper() es.IHelper
}

type app struct {
	logHelper log.IHelper
	esHelper  es.IHelper
}

func NewApp() IApp {
	a := &app{}
	a.initLogHelper()
	a.initESHelper()
	return a
}

func (a *app) initLogHelper() {
	config := log.HelperConfig{}
	if err := dotenv.Decode(&config); err != nil {
		panic(err)
	}
	inst, err := log.NewHelper(config)
	if err != nil {
		panic(err)
	}
	a.logHelper = inst
}

func (a *app) initESHelper() {
	config := es.HelperConfig{}
	if err := dotenv.Decode(&config); err != nil {
		panic(err)
	}
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
