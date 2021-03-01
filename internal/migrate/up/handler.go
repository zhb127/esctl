package up

import (
	"esctl/internal/app"
	"esctl/internal/migrate"
	"esctl/pkg/es"
	"esctl/pkg/log"
	"os"
	"sort"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	indexAliasCreate "esctl/internal/index/alias/create"
	indexAliasDelete "esctl/internal/index/alias/delete"
	indexCreate "esctl/internal/index/create"
	indexDelete "esctl/internal/index/delete"
	indexMove "esctl/internal/index/move"
)

type IHandler interface {
	Run(flags *HandlerFlags) error
	ParseCmdFlags(cmdFlags *pflag.FlagSet) (*HandlerFlags, error)
}

type handler struct {
	logHelper   log.IHelper
	esHelper    es.IHelper
	subHandlers *handlerSubHandlers

	svcUp IService
}

type handlerSubHandlers struct {
	IndexCreate      indexCreate.IHandler
	indexDelete      indexDelete.IHandler
	IndexMove        indexMove.IHandler
	IndexAliasCreate indexAliasCreate.IHandler
	indexAliasDelete indexAliasDelete.IHandler
}

func NewHandler(a app.IApp) IHandler {
	h := &handler{
		logHelper: a.LogHelper(),
		esHelper:  a.ESHelper(),
		subHandlers: &handlerSubHandlers{
			IndexCreate:      indexCreate.NewHandler(a),
			indexDelete:      indexDelete.NewHandler(a),
			IndexMove:        indexMove.NewHandler(a),
			IndexAliasCreate: indexAliasCreate.NewHandler(a),
			indexAliasDelete: indexAliasDelete.NewHandler(a),
		},
	}

	h.svcUp = NewService(h.logHelper, h.esHelper)

	return h
}

type HandlerFlags struct {
	Dir     string
	From    string
	To      string
	Reverse bool
}

func (h *handler) Run(flags *HandlerFlags) error {
	if err := h.svcUp.InitMigrationHistoryRepo(); err != nil {
		h.logHelper.Error("failed to init migration history repo", map[string]interface{}{
			"error": err.Error(),
		})
	}

	if flags.Reverse {
		return h.runMigrateDown(flags)
	}

	if err := h.runMigrateUp(flags); err != nil {
		return err
	}

	return nil
}

func (h *handler) runMigrateUp(flags *HandlerFlags) error {
	// 列出迁移文件名称
	mgrFileNames, err := h.listMigrateUpFileNames(flags.Dir)
	if err != nil {
		return err
	}

	if len(mgrFileNames) == 0 {
		return errors.Errorf("there is no migration files")
	}

	h.logHelper.Debug("list migration files", map[string]interface{}{
		"count": len(mgrFileNames),
	})

	// 获取最后执行的迁移名称
	mgrNameLastExecuted, err := h.svcUp.GetUpMigrationNameLastExecuted()
	if err != nil {
		return err
	}

	h.logHelper.Debug("get migration name last executed", map[string]interface{}{
		"name": mgrNameLastExecuted,
	})

	isStartMigration := false
	for _, fileName := range mgrFileNames {
		mgrFilePath := flags.Dir + "/" + fileName
		mgrName := strings.TrimSuffix(fileName, migrate.UP_MIGRATION_FILE_SUFFIX)

		// 判断是否始执行迁移
		if !isStartMigration {
			if mgrNameLastExecuted != "" {
				if mgrNameLastExecuted != mgrName {
					h.logHelper.Debug("file did not match, skipped", map[string]interface{}{
						"file":          fileName,
						"last_executed": mgrNameLastExecuted,
					})
					continue
				} else {
					isStartMigration = true
					continue
				}
			}
		}

		// 解析迁移文件
		mgr, err := h.parseMigrationFile(mgrFilePath)
		if err != nil {
			return err
		}

		// 执行迁移
		if err := h.execMigration(mgr); err != nil {
			return err
		}

		// 保存迁移历史条目
		if err := h.svcUp.SaveMigrationHistoryEntry(mgrName); err != nil {
			return err
		}

		if flags.To != "" {
			if flags.To == mgrName {
				h.logHelper.Debug("file match --to, end", map[string]interface{}{
					"file": fileName,
				})
				break
			}
		}
	}

	return nil
}

func (h *handler) runMigrateDown(flags *HandlerFlags) error {
	// 获取最后迁移版本
	lastMgrName, err := h.svcUp.GetUpMigrationNameLastExecuted()
	if err != nil {
		return err
	}
	h.logHelper.Debug("get last migration name", map[string]interface{}{
		"name": lastMgrName,
	})
	if lastMgrName == "" {
		return errors.New("last migration name is empty")
	}

	// 获取最后迁移文件
	filePath := flags.Dir + "/" + lastMgrName + migrate.DOWN_MIGRATION_FILE_SUFFIX

	migration, err := h.parseMigrationFile(filePath)
	if err != nil {
		return err
	}

	if err := h.execMigration(migration); err != nil {
		return err
	}

	if err := h.svcUp.DeleteMigrationHistoryEntry(lastMgrName); err != nil {
		return err
	}

	return nil
}

func (h *handler) ParseCmdFlags(flags *pflag.FlagSet) (*HandlerFlags, error) {
	handlerFlags := &HandlerFlags{}

	if dir, err := flags.GetString("dir"); err != nil {
		return nil, err
	} else {
		handlerFlags.Dir = dir
	}

	if from, err := flags.GetString("from"); err != nil {
		return nil, err
	} else {
		handlerFlags.From = from
	}

	if to, err := flags.GetString("to"); err != nil {
		return nil, err
	} else {
		handlerFlags.To = to
	}

	if reverse, err := flags.GetBool("reverse"); err != nil {
		return nil, err
	} else {
		handlerFlags.Reverse = reverse
	}

	return handlerFlags, nil
}

func (h *handler) listMigrateUpFileNames(dir string) ([]string, error) {
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
		fName := file.Name()
		if strings.HasSuffix(fName, migrate.UP_MIGRATION_FILE_SUFFIX) {
			res = append(res, fName)
		}
	}

	sort.Strings(res)
	return res, nil
}

func (h *handler) parseMigrationFile(file string) (*migrate.Migration, error) {
	h.logHelper.Debug("start parse migration file", map[string]interface{}{
		"file": file,
	})

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
	h.logHelper.Debug("start exec migration", map[string]interface{}{
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
		case "index-delete":
			if err := h.subHandlers.indexDelete.Run(v.Args); err != nil {
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
		case "index-alias-delete":
			flags := &indexAliasDelete.HandlerFlags{
				Index: v.Flags["index"].(string),
				Alias: v.Flags["alias"].(string),
			}
			if err := h.subHandlers.indexAliasDelete.Run(flags); err != nil {
				return err
			}
		default:
			return errors.Errorf("cmd=%v not supported", v.CMD)
		}
	}

	return nil
}
