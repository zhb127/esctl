package app

import (
	"esctl/pkg/config/dotenv"
	"esctl/pkg/es"
	"esctl/pkg/log"
)

type appConfig struct {
	LogHelper log.HelperConfig
	ESHelper  es.HelperConfig
}

func newAppConfig() *appConfig {
	res := &appConfig{}

	if err := dotenv.Decode(res); err != nil {
		panic(err)
	}

	return res
}
