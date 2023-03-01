package server

import (
	ratelimit2 "github.com/go-kratos/aegis/ratelimit"
	prom "github.com/go-kratos/kratos/contrib/metrics/prometheus/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/metrics"
	"github.com/go-kratos/kratos/v2/middleware/ratelimit"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"golang.org/x/time/rate"
	"time"
	pb "user_info/api/user_info"
	"user_info/conf"
	"user_info/internal/service"
)

// Limiter 使用令牌桶实现限流器
type Limiter struct {
	limiter *rate.Limiter
}

func (l *Limiter) Allow() (ratelimit2.DoneFunc, error) {
	if !l.limiter.Allow() {
		return nil, ratelimit2.ErrLimitExceed
	}
	return func(info ratelimit2.DoneInfo) {

	}, nil
}

// NewGRPCServer new a gRPC server.
func NewGRPCServer(c *conf.Server, service *service.UserInfoService, logger log.Logger) *grpc.Server {
	tokenLimiter := Limiter{limiter: rate.NewLimiter(rate.Every(500*time.Microsecond), 1)}

	var opts = []grpc.ServerOption{
		grpc.Middleware(
			ratelimit.Server(ratelimit.WithLimiter(&tokenLimiter)),
			recovery.Recovery(),
			metrics.Server(
				metrics.WithSeconds(prom.NewHistogram(_metricServerSeconds)),
				metrics.WithRequests(prom.NewCounter(_metricServerQPS)),
			),
		),
	}
	if c.Grpc.Network != "" {
		opts = append(opts, grpc.Network(c.Grpc.Network))
	}
	if c.Grpc.Addr != "" {
		opts = append(opts, grpc.Address(c.Grpc.Addr))
		service.IpAddr = c.Grpc.Addr
	}
	if c.Grpc.Timeout != nil {
		opts = append(opts, grpc.Timeout(c.Grpc.Timeout.AsDuration()))
	}
	srv := grpc.NewServer(opts...)
	pb.RegisterUserInfoServer(srv, service)
	return srv
}
