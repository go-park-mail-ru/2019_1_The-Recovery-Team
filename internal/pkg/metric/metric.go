package metric

import "github.com/prometheus/client_golang/prometheus"

var (
	AccessHits   *prometheus.CounterVec
	TotalRooms   prometheus.Gauge
	TotalPlayers prometheus.Gauge
)

// RegisterAccessHitsMetric registers new access metric
func RegisterAccessHitsMetric(namespace string) {
	AccessHits = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: namespace,
		Name:      "hits_by_status_code",
		Help:      "Total hits sorted by status codes",
	},
		[]string{"status_code", "path", "method"},
	)
	prometheus.MustRegister(AccessHits)
}

// RegisterTotalRoomsMetric registers new total rooms metric
func RegisterTotalRoomsMetric(namespace string) {
	TotalRooms = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "total_rooms",
		Help:      "Total number of rooms",
	})
	prometheus.MustRegister(TotalRooms)
}

// RegisterTotalPlayersMetric registers new total players metric
func RegisterTotalPlayersMetric(namespace string) {
	TotalPlayers = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "total_players",
		Help:      "Total number of players",
	})
	prometheus.MustRegister(TotalPlayers)
}
