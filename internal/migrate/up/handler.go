package up

import (
	"esctl/internal/app"
	"esctl/internal/migrate"
	"esctl/pkg/es"
	"esctl/pkg/log"
	"os"

	"github.com/spf13/pflag"
)

type IHandler interface {
	Run(flags *HandlerFlags) error
	ParseCmdFlags(cmdFlags *pflag.FlagSet) (*HandlerFlags, error)
}

type handler struct {
	app            app.IApp
	logHelper      log.IHelper
	esHelper       es.IHelper
	migrateService migrate.IService
}

func NewHandler(a app.IApp) IHandler {
	return &handler{
		app:            a,
		logHelper:      a.LogHelper(),
		esHelper:       a.ESHelper(),
		migrateService: migrate.NewService(a),
	}
}

type HandlerFlags struct {
	Dir string
}

func (h *handler) Run(flags *HandlerFlags) error {
	f, err := os.Open(flags.Dir)
	if err != nil {
		return err
	}
	defer f.Close()

	files, err := f.Readdir(-1)
	if err != nil {
		return nil
	}

	for _, file := range files {
		migrationFile := flags.Dir + "/" + file.Name()

		migration, err := h.migrateService.ParseMigrationFile(migrationFile)
		if err != nil {
			return err
		}

		if err := h.migrateService.ExecMigration(migration); err != nil {
			return err
		}
	}

	return nil
}

func (h *handler) ParseCmdFlags(cmdFlags *pflag.FlagSet) (*HandlerFlags, error) {
	handlerFlags := &HandlerFlags{}

	flagDir, err := cmdFlags.GetString("dir")
	if err != nil {
		return nil, err
	}
	handlerFlags.Dir = flagDir

	return handlerFlags, nil
}
