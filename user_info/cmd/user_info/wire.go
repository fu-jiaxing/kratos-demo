//go:build wireinject
// +build wireinject

//
package main

import (
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	consul_api "github.com/hashicorp/consul/api"
	"user_info/conf"
	"user_info/internal/biz"
	"user_info/internal/data"
	"user_info/internal/server"
	"user_info/internal/service"
)

func WireApp(*conf.Server, *conf.Data, log.Logger, *consul_api.Client) (*kratos.App, func(), error) {
	panic(wire.Build(data.ProviderSet, biz.ProviderSet, service.ProviderSet, server.ProviderSet, newApp))
}
