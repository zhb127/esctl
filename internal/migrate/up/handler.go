package up

import (
	"esctl/internal/app"
	"esctl/internal/migrate"
	"esctl/pkg/log"
	"os"
	"path"
	"sort"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	indexAliasCreate "esctl/internal/index/alias/create"
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
	IndexCreate      indexCreate.IHandler
	IndexMove        indexMove.IHandler
	IndexAliasCreate indexAliasCreate.IHandler
}

func NewHandler(a app.IApp) IHandler {
	return &handler{
		logHelper: a.LogHelper(),
		subHandlers: &handlerSubHandlers{
			IndexCreate:      indexCreate.NewHandler(a),
			IndexMove:        indexMove.NewHandler(a),
			IndexAliasCreate: indexAliasCreate.NewHandler(a),
		},
	}
}

type HandlerFlags struct {
	Dir  string
	From string
	To   string
}

func (h *handler) Run(flags *HandlerFlags) error {
	mgrFileNames, err := h.listMigrationFileNames(flags.Dir)
	if err != nil {
		return err
	}

	h.logHelper.Debug("list migration files", map[string]interface{}{
		"count": len(mgrFileNames),
	})

	for _, mgrFileName := range mgrFileNames {
		mgrFilePath := flags.Dir + "/" + mgrFileName
		mgrFileExt := path.Ext(mgrFilePath)
		mgrFileNameWithoutExt := strings.TrimSuffix(mgrFileName, mgrFileExt)

		// 后缀名不一致，则跳过
		if mgrFileExt != ".yaml" && mgrFileExt != ".yml" {
			h.logHelper.Warn("file ext is not (.yaml|.yml)", map[string]interface{}{
				"file": mgrFileName,
				"ext":  mgrFileExt,
			})
			continue
		}

		// 判断是否开始
		if flags.From != "" {
			if flags.From != mgrFileNameWithoutExt {
				h.logHelper.Debug("file did not match --from, not start", map[string]interface{}{
					"file": mgrFileName,
				})
				continue
			}
		}

		h.logHelper.Debug("start parse migration file", map[string]interface{}{
			"file": mgrFilePath,
		})

		migration, err := h.parseMigrationFile(mgrFilePath)
		if err != nil {
			return err
		}

		h.logHelper.Debug("start exec migration file", map[string]interface{}{
			"file": mgrFilePath,
		})

		if err := h.execMigration(migration); err != nil {
			return err
		}

		// 判断是否结束
		if flags.To != "" {
			if flags.To == mgrFileNameWithoutExt {
				h.logHelper.Debug("file match --to, done", map[string]interface{}{
					"file": mgrFileName,
				})
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

func (h *handler) listMigrationFileNames(dir string) ([]string, error) {
	fd, err := os.Open(dir)
	if err != nil {
		return nil, err
	}
	defer fd.Close()

	files, err := fd.Readdir(-1)
	if err != nil {
		return nil, err
	}

	var res []string
	for _, file := range files {
		res = append(res, file.Name())
	}

	sort.Strings(res)

	return res, nil
}

func (h *handler) parseMigrationFile(file string) (*migrate.Migration, error) {
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

func (h *handler) execMigration(migration *migrate.Migration) error {
	h.logHelper.Debug("exec migration", map[string]interface{}{
		"migration": migration,
	})

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
		case "index-alias-create":
			flags := &indexAliasCreate.HandlerFlags{
				Index: v.Flags["index"].(string),
				Alias: v.Flags["alias"].(string),
			}
			if err := h.subHandlers.IndexAliasCreate.Run(flags); err != nil {
				return err
			}
		default:
			return errors.Errorf("cmd=%v not supported", v.CMD)
		}
	}

	return nil
}
