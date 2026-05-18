package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	RegisterUserTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "user_service_register_total",
		Help: "Total number of user registrations",
	})

	LoginUserTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "user_service_login_total",
		Help: "Total number of login attempts",
	})

	LoginUserFailedTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "user_service_login_failed_total",
		Help: "Total number of failed login attempts",
	})

	DatabaseQueryDuration = promauto.NewHistogram(prometheus.HistogramOpts{
		Name:    "user_service_db_query_duration_seconds",
		Help:    "Database query duration in seconds",
		Buckets: prometheus.DefBuckets,
	})
)
