package create

import (
	"esctl/internal/app"
	"esctl/internal/migrate"
	"esctl/pkg/log"
	"os"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/pflag"

	indexAliasCreate "esctl/internal/index/alias/create"
	indexCreate "esctl/internal/index/create"
	indexMove "esctl/internal/index/move"
)

type IHandler interface {
	Run(flags *HandlerFlags) error
	ParseCmdFlags(cmdFlags *pflag.FlagSet) (*HandlerFlags, error)
}

type handler struct {
	logHelper log.IHelper
}

type handlerSubHandlers struct {
	IndexCreate      indexCreate.IHandler
	IndexMove        indexMove.IHandler
	IndexAliasCreate indexAliasCreate.IHandler
}

func NewHandler(a app.IApp) IHandler {
	return &handler{
		logHelper: a.LogHelper(),
	}
}

type HandlerFlags struct {
	Dir  string
	Name string
}

func (h *handler) Run(flags *HandlerFlags) error {
	// 验证选项 --dir
	mgrDir, err := os.Stat(flags.Dir)
	if err != nil {
		return errors.Wrap(err, "failed to check migrations file directory")
	}
	if !mgrDir.IsDir() {
		return errors.New("--dir is not directory")
	}

	// 生成迁移文件路径
	mgrVer := time.Now().Format("20060102150405")
	mgrUpFilePath := flags.Dir + "/" + mgrVer + "_" + flags.Name + migrate.MIGRATION_UP_FILE_EXT
	mgrDownFilePath := flags.Dir + "/" + mgrVer + "_" + flags.Name + migrate.MIGRATION_DOWN_FILE_EXT

	// 创建向上迁移文件
	mgrUpFile, err := os.Create(mgrUpFilePath)
	if err != nil {
		return errors.Wrap(err, "failed to create migration up file")
	}
	defer mgrUpFile.Close()
	h.logHelper.Info("create migration up file", map[string]interface{}{"path": mgrUpFilePath})

	// 创建向下迁移文件
	mgrDownFile, err := os.Create(mgrDownFilePath)
	if err != nil {
		return errors.Wrap(err, "failed to create migration down file")
	}
	defer mgrDownFile.Close()
	h.logHelper.Info("create migration down file", map[string]interface{}{"path": mgrDownFilePath})

	return nil
}

func (h *handler) ParseCmdFlags(cmdFlags *pflag.FlagSet) (*HandlerFlags, error) {
	handlerFlags := &HandlerFlags{}

	flagDir, err := cmdFlags.GetString("dir")
	if err != nil {
		return nil, err
	}
	handlerFlags.Dir = flagDir

	flagName, err := cmdFlags.GetString("name")
	if err != nil {
		return nil, err
	}
	handlerFlags.Name = flagName

	return handlerFlags, nil
}
