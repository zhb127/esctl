package es

import (
	"esctl/pkg/config/dotenv"
	"esctl/pkg/es"
	"esctl/pkg/log"
)

func MockHelper(logHelper log.IHelper) es.IHelper {
	config := es.HelperConfig{}
	if err := dotenv.Decode(&config); err != nil {
		panic(err)
	}
	inst, err := es.NewHelper(config, logHelper)
	if err != nil {
		panic(err)
	}
	return inst
}
