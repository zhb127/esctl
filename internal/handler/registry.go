// +build wireinject

package handler

import (
	"esctl/internal/app"
	indexList "esctl/internal/index/list"
	"esctl/internal/version"
	"esctl/pkg/es"
	"esctl/pkg/log"

	"github.com/google/wire"
)

type Registry struct {
	Version   version.IHandler
	IndexList indexList.IHandler
}

var providerSet = wire.NewSet(version.NewHandler, indexList.NewHandler)

func NewRegistry(a app.IApp, logHelper log.IHelper, esHelper es.IHelper) *Registry {
	panic(wire.Build(providerSet, wire.Struct(new(Registry), "*")))
}
