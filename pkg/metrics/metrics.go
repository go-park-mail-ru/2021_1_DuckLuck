package metrics

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
)

type Metrics struct {
	TotalHits  prometheus.Counter
	AccessHits *prometheus.CounterVec
	Errors     *prometheus.CounterVec
	Durations  *prometheus.HistogramVec
}

func CreateNewMetrics(name string) (*Metrics, error) {
	metrics := &Metrics{}
	metrics.TotalHits = prometheus.NewCounter(prometheus.CounterOpts{
		Name: fmt.Sprintf("%s_total", name),
	})
	if err := prometheus.Register(metrics.TotalHits); err != nil {
		return nil, err
	}

	metrics.AccessHits = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: fmt.Sprintf("%s_access_hits", name),
	}, []string{"status_code", "method", "path"})
	if err := prometheus.Register(metrics.AccessHits); err != nil {
		return nil, err
	}

	metrics.Errors = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: fmt.Sprintf("%s_errors", name),
	}, []string{"status_code", "method", "path"})
	if err := prometheus.Register(metrics.Errors); err != nil {
		return nil, err
	}

	metrics.Durations = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name: fmt.Sprintf("%s_durations", name),
	}, []string{"status_code", "method", "path"})
	if err := prometheus.Register(metrics.Durations); err != nil {
		return nil, err
	}

	return metrics, nil
}
