package metric

import "github.com/prometheus/client_golang/prometheus"

var AccessHits *prometheus.CounterVec

// RegisterAccessHitsMetric registers new access metric
func RegisterAccessHitsMetric(namespace string) {
	AccessHits = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: namespace,
		Name:      "hits_by_status_code",
		Help:      "Total hits sorted by status codes",
	},
		[]string{"status_code", "path", "method"},
	)
}
