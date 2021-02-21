package app

import (
	"esctl/pkg/config/dotenv"
	"esctl/pkg/es"
	"esctl/pkg/log"
)

type appConfig struct {
	logHelper log.HelperConfig
	esHelper  es.HelperConfig
}

func newAppConfig() *appConfig {
	c := &appConfig{}

	if err := dotenv.Decode(c); err != nil {
		panic(err)
	}

	return c
}
