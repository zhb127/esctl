package migrate

import (
	"esctl/internal/index/app"
	"esctl/pkg/es"
	"esctl/pkg/log"

	"github.com/pkg/errors"
	"github.com/spf13/pflag"
)

type IHandler interface {
	Run(flags *HandlerFlags) error
	ParseCmdFlags(cmdFlags *pflag.FlagSet) (*HandlerFlags, error)
}

type HandlerFlags struct {
	Src  string
	Dest string
}

type handler struct {
	logHelper log.IHelper
	esHelper  es.IHelper
}

func NewHandler(a app.IApp) IHandler {
	return &handler{
		logHelper: a.LogHelper(),
		esHelper:  a.ESHelper(),
	}
}

func (h *handler) Run(flags *HandlerFlags) error {
	// get src alias
	listAliasesResp, err := h.esHelper.ListAliases()
	if err != nil {
		return errors.Wrap(err, "find srcIndex aliases")
	}

	srcAliases := []string{}
	for _, v := range listAliasesResp.Items {
		if v.Index == flags.Src {
			srcAliases = append(srcAliases, v.Alias)
		}
	}

	h.logHelper.Debug("find srcIndex aliases", map[string]interface{}{
		"aliases": srcAliases,
	})

	// reindex
	reindexResp, err := h.esHelper.Reindex(flags.Src, flags.Dest)
	if err != nil {
		return errors.Wrapf(err, "reindex srcIndex to destIndex")
	}

	h.logHelper.Debug("reindex srcIndex to DestIndex", map[string]interface{}{
		"result": reindexResp,
	})

	// copy alias
	for _, v := range srcAliases {
		aliasIndexResp, err := h.esHelper.AliasIndex(flags.Src, v)
		if err != nil {
			return errors.Wrapf(err, "copy alias=%s from srcIndex=%s to destIndex=%s", v, flags.Src, flags.Dest)
		}
		h.logHelper.Debug("copy alias=%s from srcIndex=%s to destIndex=%s", map[string]interface{}{
			"result": aliasIndexResp,
		})
	}

	return nil
}

func (h *handler) ParseCmdFlags(cmdFlags *pflag.FlagSet) (*HandlerFlags, error) {
	result := &HandlerFlags{}

	flagSrc, err := cmdFlags.GetString("src")
	if err != nil {
		return nil, err
	}
	result.Src = flagSrc

	flagDest, err := cmdFlags.GetString("dest")
	if err != nil {
		return nil, err
	}
	result.Dest = flagDest

	return result, nil
}
