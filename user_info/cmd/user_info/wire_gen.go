// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"user_info/conf"
	"user_info/internal/biz"
	"user_info/internal/data"
	"user_info/internal/server"
	"user_info/internal/service"
)

// Injectors from wire.go:

func WireApp(confServer *conf.Server, confData *conf.Data, logger log.Logger) (*kratos.App, func(), error) {
	dataData, cleanup, err := data.NewData(confData, logger)
	if err != nil {
		return nil, nil, err
	}
	userRepo := data.NewUserRepo(dataData, logger)
	redisRepo := data.NewRedis(dataData)
	userBizHandler := biz.NewUserBizHandler(userRepo, redisRepo)
	userInfoService := service.NewUserInfoService(userBizHandler)
	grpcServer := server.NewGRPCServer(confServer, userInfoService, logger)
	httpServer := server.NewHTTPServer(confServer, userInfoService, logger)
	app := newApp(logger, grpcServer, httpServer)
	return app, func() {
		cleanup()
	}, nil
}
