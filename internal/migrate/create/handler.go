package create

import (
	"esctl/internal/app"
	"esctl/internal/migrate"
	"esctl/pkg/log"
	"os"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/pflag"
)

type IHandler interface {
	Run(flags *HandlerFlags) error
	ParseCmdFlags(cmdFlags *pflag.FlagSet) (*HandlerFlags, error)
}

type handler struct {
	logHelper log.IHelper
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
	// 生成迁移文件路径
	mgrVer := time.Now().Format("20060102150405")
	mgrUpFilePath := flags.Dir + "/" + mgrVer + "_" + flags.Name + migrate.MIGRATION_UP_FILE_SUFFIX
	mgrDownFilePath := flags.Dir + "/" + mgrVer + "_" + flags.Name + migrate.MIGRATION_DOWN_FILE_SUFFIX

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

func (h *handler) ParseCmdFlags(flags *pflag.FlagSet) (*HandlerFlags, error) {
	handlerFlags := &HandlerFlags{}

	if dir, err := flags.GetString("dir"); err != nil {
		return nil, err
	} else {
		handlerFlags.Dir = dir
	}

	if name, err := flags.GetString("name"); err != nil {
		return nil, err
	} else {
		handlerFlags.Name = name
	}

	if err := h.validateHandlerFlags(handlerFlags); err != nil {
		return nil, err
	}

	return handlerFlags, nil
}

func (h *handler) validateHandlerFlags(flags *HandlerFlags) error {
	mgrDir, err := os.Stat(flags.Dir)
	if err != nil {
		return errors.Wrap(err, "failed to check migrations file directory")
	}
	if !mgrDir.IsDir() {
		return errors.New("--dir is not directory")
	}

	return nil
}
