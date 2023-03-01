package server

import (
	"context"
	ratelimit2 "github.com/go-kratos/aegis/ratelimit"
	prom "github.com/go-kratos/kratos/contrib/metrics/prometheus/v2"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/metrics"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"golang.org/x/time/rate"
	"math/rand"
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
		return nil, errors.New(504, "rate limit exceeded", "testsettetsetstst")
	}
	return func(info ratelimit2.DoneInfo) {

	}, nil
}

func RandomError500(handler middleware.Handler) middleware.Handler {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		if rand.Intn(10) < 5 {
			return nil, errors.New(500, "random error", "random error")
		}
		reply, err := handler(ctx, req)
		return reply, err
	}
}

// NewGRPCServer new a gRPC server.
func NewGRPCServer(c *conf.Server, service *service.UserInfoService, logger log.Logger) *grpc.Server {
	// tokenLimiter := Limiter{limiter: rate.NewLimiter(rate.Every(100*time.Millisecond), 1)}

	var opts = []grpc.ServerOption{
		grpc.Middleware(
			// ratelimit.Server(ratelimit.WithLimiter(&tokenLimiter)),
			recovery.Recovery(),
			// RandomError500,
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
