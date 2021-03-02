package move

import (
	"esctl/internal/app"
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
		return errors.Wrap(err, "failed to find srcIndex/destIndex")
	}
	if len(listIndicesResp.Items) < 2 {
		return errors.Wrap(err, "srcIndex or destIndex not found")
	}

	// list src aliases
	listAliasesResp, err := h.esHelper.ListAliases()
	if err != nil {
		return errors.Wrap(err, "failed to find srcIndex aliases")
	}

	srcIndexAliases := []string{}
	for _, v := range listAliasesResp.Items {
		if v.Index == flags.Src {
			srcIndexAliases = append(srcIndexAliases, v.Alias)
		}
	}

	h.logHelper.Info("find srcIndex aliases", map[string]interface{}{
		"result": srcIndexAliases,
	})

	// reindex
	reindexResp, err := h.esHelper.Reindex(flags.Src, flags.Dest)
	if err != nil {
		return errors.Wrap(err, "failed to reindex srcIndex to destIndex")
	}

	h.logHelper.Info("reindex srcIndex to DestIndex", map[string]interface{}{
		"result": reindexResp,
	})

	if len(srcIndexAliases) > 0 {
		// copy srcIndex aliases to destIndex
		for _, v := range srcIndexAliases {
			aliasIndexResp, err := h.esHelper.AliasIndex(flags.Dest, v)
			if err != nil {
				return errors.Wrap(err, "failed to copy alias from srcIndex to destIndex")
			}
			h.logHelper.Info("copy alias from srcIndex to destIndex", map[string]interface{}{
				"result": aliasIndexResp,
			})
		}

		// delete src aliases
		unaliasIndexResp, err := h.esHelper.UnaliasIndex(flags.Src, srcIndexAliases)
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
			return errors.Wrap(err, "failed to purge srcIndex")
		}
		h.logHelper.Info("purge srcIndex", map[string]interface{}{
			"result": deleteIndexResp,
		})
	}

	return nil
}

func (h *handler) ParseCmdFlags(cmdFlags *pflag.FlagSet) (*HandlerFlags, error) {
	result := &HandlerFlags{}

	if src, err := cmdFlags.GetString("src"); err != nil {
		return nil, err
	} else {
		result.Src = src
	}

	if dest, err := cmdFlags.GetString("dest"); err != nil {
		return nil, err
	} else {
		result.Dest = dest
	}

	if purge, err := cmdFlags.GetBool("purge"); err != nil {
		return nil, err
	} else {
		result.Purge = purge
	}

	return result, nil
}
