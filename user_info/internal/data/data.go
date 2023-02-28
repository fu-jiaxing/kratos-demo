package data

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-redis/redis"
	"github.com/google/wire"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"user_info/conf"
)

var ProviderSet = wire.NewSet(NewData, NewUserRepo, NewRedis)

type Data struct {
	db    *gorm.DB
	redis *redis.Client
}

func NewData(c *conf.Data, logger log.Logger) (*Data, func(), error) {
	db, err := gorm.Open(mysql.Open(c.Database.Source), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	rdb := redis.NewClient(&redis.Options{
		Addr:         c.Redis.Addr,
		ReadTimeout:  c.Redis.ReadTimeout.AsDuration(),
		WriteTimeout: c.Redis.WriteTimeout.AsDuration(),
	})
	data := Data{
		db:    db,
		redis: rdb,
	}

	cleanup := func() {
		log.NewHelper(logger).Info("closing the data resources")
		rdb.Close()
	}

	return &data, cleanup, nil
}
