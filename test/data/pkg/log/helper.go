package log

import (
	"esctl/pkg/config/dotenv"
	"esctl/pkg/log"
)

func MockHelper() log.IHelper {
	config := log.HelperConfig{}
	if err := dotenv.Decode(&config); err != nil {
		panic(err)
	}
	inst, err := log.NewHelper(config)
	if err != nil {
		panic(err)
	}
	return inst
}
