package server

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	Name                 = "bbs_user_info"
	_metricServerSeconds = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace:   Name + "_server",
		Subsystem:   "requests",
		Name:        "duration_sec",
		Help:        "server requests duration(sec).",
		ConstLabels: nil,
		Buckets:     []float64{0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1},
	}, []string{"kind", "operation"})

	_metricServerQPS = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace:   Name + "_server",
		Subsystem:   "requests",
		Name:        "throughput",
		Help:        "throughput",
		ConstLabels: nil,
	}, []string{"kind", "operation", "code", "reason"})
)

func init() {
	prometheus.MustRegister(_metricServerQPS, _metricServerSeconds)
}
