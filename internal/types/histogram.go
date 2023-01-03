package types

import (
	"strconv"
	"strings"

	"github.com/cocaccola/rumpelstiltskin/internal/helpers"
	"github.com/prometheus/client_golang/prometheus"
)

type Histogram struct {
	Name         string                   `yaml:"name"`
	Help         string                   `yaml:"help"`
	Labels       []string                 `yaml:"labels"`
	Buckets      []float64                `yaml:"buckets"`
	Schedule     []HistogramSchedule      `yaml:"schedule"`
	HistogramVec *prometheus.HistogramVec `yaml:"-"`
}

type HistogramSchedule struct {
	Intervals uint64              `yaml:"intervals"`
	Behaviors []HistogramBehavior `yaml:"behaviors"`
}

type HistogramBehavior struct {
	Labels prometheus.Labels `yaml:"labels"`
	Values []string          `yaml:"values"`
}

func (h *Histogram) Register(reg prometheus.Registerer, prefix string) {
	histogram := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: prefix,
			Name:      h.Name,
			Help:      h.Help,
			Buckets:   h.Buckets,
		},
		h.Labels,
	)

	reg.MustRegister(histogram)
	h.HistogramVec = histogram
}

func (h *Histogram) GetIntervals(index int) uint64 {
	return h.Schedule[index].Intervals
}

func (h *Histogram) GetScheduleLength() int {
	return len(h.Schedule)
}

func (h *Histogram) SetValues(index int) {
	for _, behavior := range h.Schedule[index].Behaviors {
		for _, valuePair := range behavior.Values {
			pair := strings.Split(valuePair, "|")

			obs, err := strconv.ParseUint(pair[1], 10, 64)
			helpers.PanicWithConfigError(err)

			for i := obs; i > 0; i-- {
				value, err := strconv.ParseFloat(pair[0], 64)
				helpers.PanicWithConfigError(err)

				h.HistogramVec.With(behavior.Labels).Observe(value - 1)
			}
		}
	}
}
