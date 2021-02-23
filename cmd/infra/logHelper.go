package infra

import (
	"esctl/pkg/log"
)

var LogHelper log.IHelper

func InitLogHelper(config log.HelperConfig) error {
	res, err := log.NewHelper(config)
	if err != nil {
		return err
	}

	LogHelper = res
	return nil
}
