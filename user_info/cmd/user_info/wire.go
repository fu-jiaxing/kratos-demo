//go:build wireinject
// +build wireinject
package main

import (
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"user_info/conf"
	"user_info/internal/biz"
	"user_info/internal/server"
	"user_info/internal/service"
)

func WireApp(*conf.Server, *conf.Data, log.Logger) (*kratos.App, func(), error) {
	panic(wire.Build(server.ProviderSet, biz.ProviderSet, service.ProviderSet, newApp))
}
