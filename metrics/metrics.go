package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var ModelPredictions = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "model_predictions",
		Help: "Model predictions",
	},
	[]string{"model"},
)

func init() {
	// Register the custom counter metric with Prometheus
	prometheus.MustRegister(ModelPredictions)
}
