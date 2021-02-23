package up

import (
	"esctl/internal/app"
	"esctl/internal/migrate"
	"esctl/pkg/log"
	"os"
	"path"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	indexCreate "esctl/internal/index/create"
	indexMove "esctl/internal/index/move"
)

type IHandler interface {
	Run(flags *HandlerFlags) error
	ParseCmdFlags(cmdFlags *pflag.FlagSet) (*HandlerFlags, error)
}

type handler struct {
	logHelper   log.IHelper
	subHandlers *handlerSubHandlers
}

type handlerSubHandlers struct {
	IndexCreate indexCreate.IHandler
	IndexMove   indexMove.IHandler
}

func NewHandler(a app.IApp) IHandler {
	return &handler{
		logHelper: a.LogHelper(),
		subHandlers: &handlerSubHandlers{
			IndexCreate: indexCreate.NewHandler(a),
			IndexMove:   indexMove.NewHandler(a),
		},
	}
}

type HandlerFlags struct {
	Dir  string
	From string
	To   string
}

func (h *handler) Run(flags *HandlerFlags) error {
	mgrFiles, err := h.ListMigrationFiles(flags.Dir)
	if err != nil {
		return err
	}

	for _, mgrFile := range mgrFiles {
		mgrFileName := mgrFile.Name()
		mgrFilePath := flags.Dir + "/" + mgrFileName
		mgrFileExt := path.Ext(mgrFilePath)
		mgrFileNameWithoutExt := strings.TrimSuffix(mgrFileName, mgrFileExt)

		// 后缀名不一致，则跳过
		if mgrFileExt != "yaml" && mgrFileExt != "yml" {
			h.logHelper.Debug("ignore file because of the file ext is not (.yaml|.yml)", map[string]interface{}{
				"file": mgrFileName,
			})
			continue
		}

		// 判断是否开始
		if flags.From != "" {
			if flags.From != mgrFileNameWithoutExt {
				continue
			}
		}

		migration, err := h.ParseMigrationFile(mgrFilePath)
		if err != nil {
			return err
		}

		if err := h.ExecMigration(migration); err != nil {
			return err
		}

		// 判断是否结束
		if flags.To != "" {
			if flags.To == mgrFileNameWithoutExt {
				break
			}
		}
	}

	return nil
}

func (h *handler) ParseCmdFlags(cmdFlags *pflag.FlagSet) (*HandlerFlags, error) {
	handlerFlags := &HandlerFlags{}

	// --dir
	flagDir, err := cmdFlags.GetString("dir")
	if err != nil {
		return nil, err
	}
	handlerFlags.Dir = flagDir

	// --from
	flagFrom, err := cmdFlags.GetString("from")
	if err != nil {
		return nil, err
	}
	handlerFlags.From = flagFrom

	// --to
	flagTo, err := cmdFlags.GetString("to")
	if err != nil {
		return nil, err
	}
	handlerFlags.To = flagTo

	return handlerFlags, nil
}

func (h *handler) ListMigrationFiles(dir string) ([]os.FileInfo, error) {
	fd, err := os.Open(dir)
	if err != nil {
		return nil, err
	}
	defer fd.Close()

	res, err := fd.Readdir(-1)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (h *handler) ParseMigrationFile(file string) (*migrate.Migration, error) {
	viper.SetConfigFile(file)
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		return nil, errors.Wrap(err, "read migration file")
	}

	res := &migrate.Migration{}
	if err := viper.Unmarshal(res); err != nil {
		return nil, errors.Wrap(err, "parse migration file content")
	}

	return res, nil
}

func (h *handler) ExecMigration(migration *migrate.Migration) error {
	for _, v := range migration.CMDs {
		switch v.CMD {
		case "index-create":
			flags := &indexCreate.HandlerFlags{
				Name: v.Flags["name"].(string),
				Body: []byte(v.Flags["body"].(string)),
			}
			if err := h.subHandlers.IndexCreate.Run(flags); err != nil {
				return err
			}
		case "index-move":
			flags := &indexMove.HandlerFlags{
				Src:   v.Flags["src"].(string),
				Dest:  v.Flags["dest"].(string),
				Purge: v.Flags["purge"].(bool),
			}
			if err := h.subHandlers.IndexMove.Run(flags); err != nil {
				return err
			}
		}
	}

	return nil
}
