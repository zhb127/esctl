package move

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
	Src   string
	Dest  string
	Purge bool
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
	// find src index, dest index
	listIndicesResp, err := h.esHelper.ListIndices(flags.Src, flags.Dest)
	if err != nil {
		return errors.Wrap(err, "find srcIndex and destIndex")
	}
	if len(listIndicesResp.Items) < 2 {
		return errors.Wrap(err, "srcIndex or destIndex not found")
	}

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

	h.logHelper.Info("find srcIndex aliases", map[string]interface{}{
		"aliases": srcAliases,
	})

	// reindex
	reindexResp, err := h.esHelper.Reindex(flags.Src, flags.Dest)
	if err != nil {
		return errors.Wrap(err, "reindex srcIndex to destIndex")
	}

	h.logHelper.Info("reindex srcIndex to DestIndex", map[string]interface{}{
		"result": reindexResp,
	})

	if len(srcAliases) > 0 {
		// copy src aliases to dest
		for _, v := range srcAliases {
			aliasIndexResp, err := h.esHelper.AliasIndex(flags.Dest, v)
			if err != nil {
				return errors.Wrap(err, "copy alias from srcIndex to destIndex")
			}
			h.logHelper.Info("copy alias from srcIndex to destIndex", map[string]interface{}{
				"result": aliasIndexResp,
			})
		}

		// delete src aliases
		unaliasIndexResp, err := h.esHelper.UnaliasIndex(flags.Src, srcAliases)
		if err != nil {
			return errors.Wrap(err, "delete aliases from srcIndex")
		}
		h.logHelper.Info("delete aliases from srcIndex", map[string]interface{}{
			"result": unaliasIndexResp,
		})
	}

	// delete src index
	if flags.Purge {
		deleteIndexResp, err := h.esHelper.DeleteIndices(flags.Src)
		if err != nil {
			return errors.Wrap(err, "delete srcIndex")
		}
		h.logHelper.Info("delete srcIndex", map[string]interface{}{
			"result": deleteIndexResp,
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

	flagPurge, err := cmdFlags.GetBool("purge")
	if err != nil {
		return nil, err
	}
	result.Purge = flagPurge

	return result, nil
}
