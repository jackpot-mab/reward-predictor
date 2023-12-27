package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var ModelPredictions = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "reward_predictor_model_predictions",
		Help: "Model predictions",
	},
	[]string{"model"},
)

var ModelUpdated = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "reward_predictor_model_updated",
		Help: "Ticks when the model changed and is updated.",
	},
	[]string{"model"},
)

func init() {
	// Register the custom counter metric with Prometheus
	prometheus.MustRegister(ModelPredictions, ModelUpdated)
}
