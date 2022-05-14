// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"github.com/newpurr/stock/internal/biz"
	"github.com/newpurr/stock/internal/conf"
	"github.com/newpurr/stock/internal/data"
	"github.com/newpurr/stock/internal/data/repositories"
	"github.com/newpurr/stock/internal/event"
	"github.com/newpurr/stock/internal/server"
	"github.com/newpurr/stock/internal/service"
)

// initApp init kratos application.
func initApp(*conf.Server, *conf.Data, log.Logger) (*kratos.App, func(), error) {
	panic(wire.Build(server.ProviderSet, data.ProviderSet, biz.ProviderSet, service.ProviderSet, event.ProviderSet, repositories.ProviderSet, newApp))
}
