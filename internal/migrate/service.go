package migrate

import (
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

func ParseMigrationFile(file string) (*Migration, error) {
	viper.SetConfigFile(file)
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		return nil, nil
	}

	result := &Migration{}
	if err := viper.Unmarshal(result); err != nil {
		return nil, errors.Wrap(err, "Unmarshal migration file content")
	}

	return result, nil
}
