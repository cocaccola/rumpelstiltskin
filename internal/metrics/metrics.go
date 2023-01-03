package metrics

import (
	"fmt"
	"net/http"

	"github.com/cocaccola/rumpelstiltskin/internal/config"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var registry *prometheus.Registry

type Metric interface {
	GetIntervals(index int) uint64
	GetScheduleLength() int
	SetValues(index int)
}

func init() {
	registry = prometheus.NewRegistry()
}

func Register(config *config.Config, prefix string) {
	config.Register(registry, prefix)
}

func StartPromServer(port uint64) error {
	if port <= 0 || port >= 65535 {
		return fmt.Errorf("port out of range")
	}

	handler := promhttp.HandlerFor(registry, promhttp.HandlerOpts{})

	http.Handle("/metrics", handler)
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	return nil
}
