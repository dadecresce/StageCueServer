// internal/metrics/metrics.go
package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	PeersOnline = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "stagecue_peers_online",
			Help: "Current number of connected peers",
		},
	)
	TracksTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "stagecue_tracks_total",
			Help: "Total published tracks since start",
		},
	)
)

func MustRegisterDefault() {
	prometheus.MustRegister(PeersOnline, TracksTotal)
}
