package up

import (
	"esctl/internal/app"
	"esctl/internal/migrate"
	"esctl/pkg/es"
	"esctl/pkg/log"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/pflag"

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
	logHelper log.IHelper
	esHelper  es.IHelper

	cmds *handlerCMDs

	svcUp IService
}

type handlerCMDs struct {
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
		cmds: &handlerCMDs{
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
		return h.migrateDown(flags)
	}

	if err := h.migrateUp(flags); err != nil {
		return err
	}

	return nil
}

func (h *handler) migrateUp(flags *HandlerFlags) error {
	upMgrFileNames, err := h.svcUp.ListUpMigrationFileNames(flags.Dir)
	if err != nil {
		return err
	}

	upMgrFileNamesCount := len(upMgrFileNames)
	if upMgrFileNamesCount == 0 {
		return errors.Errorf("there is no up migration files")
	}
	h.logHelper.Debug("list up migration files", map[string]interface{}{
		"count": upMgrFileNamesCount,
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
	for _, fileName := range upMgrFileNames {
		mgrFilePath := flags.Dir + "/" + fileName
		mgrName := strings.TrimSuffix(fileName, migrate.UpMigrationFileSuffix)

		// 判断是否开始迁移
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
		mgr, err := h.svcUp.ParseMigrationFile(mgrFilePath)
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

func (h *handler) migrateDown(flags *HandlerFlags) error {
	mgrNameLastExecuted, err := h.svcUp.GetUpMigrationNameLastExecuted()
	if err != nil {
		return err
	}
	h.logHelper.Debug("get migration name last executed", map[string]interface{}{
		"name": mgrNameLastExecuted,
	})
	if mgrNameLastExecuted == "" {
		return errors.New("migration name last executed is empty")
	}

	downMgrFilePath := flags.Dir + "/" + mgrNameLastExecuted + migrate.DownMigrationFileSuffix

	mgr, err := h.svcUp.ParseMigrationFile(downMgrFilePath)
	if err != nil {
		return err
	}

	if err := h.execMigration(mgr); err != nil {
		return err
	}

	if err := h.svcUp.DeleteMigrationHistoryEntry(mgrNameLastExecuted); err != nil {
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

func (h *handler) execMigration(migration *migrate.Migration) error {
	for _, v := range migration.CMDs {
		switch v.CMD {
		case "index-create":
			flags := &indexCreate.HandlerFlags{
				Name: v.Flags["name"].(string),
				Body: []byte(v.Flags["body"].(string)),
			}
			if err := h.cmds.IndexCreate.Run(flags); err != nil {
				return err
			}
		case "index-delete":
			if err := h.cmds.indexDelete.Run(v.Args); err != nil {
				return err
			}
		case "index-move":
			flags := &indexMove.HandlerFlags{
				Src:   v.Flags["src"].(string),
				Dest:  v.Flags["dest"].(string),
				Purge: v.Flags["purge"].(bool),
			}
			if err := h.cmds.IndexMove.Run(flags); err != nil {
				return err
			}
		case "index-alias-create":
			flags := &indexAliasCreate.HandlerFlags{
				Index: v.Flags["index"].(string),
				Alias: v.Flags["alias"].(string),
			}
			if err := h.cmds.IndexAliasCreate.Run(flags); err != nil {
				return err
			}
		case "index-alias-delete":
			flags := &indexAliasDelete.HandlerFlags{
				Index: v.Flags["index"].(string),
				Alias: v.Flags["alias"].(string),
			}
			if err := h.cmds.indexAliasDelete.Run(flags); err != nil {
				return err
			}
		default:
			return errors.Errorf("cmd=%s not supported", v.CMD)
		}
	}

	return nil
}
