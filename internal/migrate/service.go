package migrate

import (
	"esctl/internal/app"
	indexCreate "esctl/internal/index/create"
	indexMove "esctl/internal/index/move"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

type IService interface {
	ParseMigrationFile(file string) (*Migration, error)
	ExecMigration(migration *Migration) error
}

type service struct {
	app      app.IApp
	handlers *handlers
}

type handlers struct {
	IndexCreate indexCreate.IHandler
	IndexMove   indexMove.IHandler
}

func NewService(app app.IApp) IService {
	h := &service{
		app:      app,
		handlers: &handlers{},
	}

	h.handlers.IndexCreate = indexCreate.NewHandler(app)

	return h
}

func (s *service) ParseMigrationFile(file string) (*Migration, error) {
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

func (s *service) ExecMigration(migration *Migration) error {
	for _, v := range migration.CMDs {
		switch v.CMD {
		case "index create":
			flags := &indexCreate.HandlerFlags{
				Name: v.Flags["name"].(string),
				Body: []byte(v.Flags["body"].(string)),
			}
			if err := s.handlers.IndexCreate.Run(flags); err != nil {
				return err
			}
		case "index move":
			flags := &indexMove.HandlerFlags{
				Src:   v.Flags["src"].(string),
				Dest:  v.Flags["dest"].(string),
				Purge: v.Flags["purge"].(bool),
			}
			if err := s.handlers.IndexMove.Run(flags); err != nil {
				return err
			}
		}
	}

	return nil
}
