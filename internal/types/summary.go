package types

import (
	"strconv"
	"strings"
	"time"

	"github.com/cocaccola/rumpelstiltskin/internal/helpers"
	"github.com/prometheus/client_golang/prometheus"
)

type Summary struct {
	Name       string                 `yaml:"name"`
	Help       string                 `yaml:"help"`
	Labels     []string               `yaml:"labels"`
	Quantiles  []string               `yaml:"quantiles"`
	MaxAge     time.Duration          `yaml:"maxAge,omitempty"`
	Schedule   []SummarySchedule      `yaml:"schedule"`
	SummaryVec *prometheus.SummaryVec `yaml:"-"`
}

type SummarySchedule struct {
	Intervals uint64            `yaml:"intervals"`
	Behaviors []SummaryBehavior `yaml:"behaviors"`
}

type SummaryBehavior struct {
	Labels prometheus.Labels `yaml:"labels"`
	Values []string          `yaml:"values"`
}

func (s *Summary) Register(reg prometheus.Registerer, prefix string) {
	summary := prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace:  prefix,
			Name:       s.Name,
			Help:       s.Help,
			Objectives: s.Objectives(),
			MaxAge:     s.MaxAge,
		},
		s.Labels,
	)

	reg.MustRegister(summary)
	s.SummaryVec = summary
}

func (s *Summary) Objectives() map[float64]float64 {
	objs := make(map[float64]float64, len(s.Quantiles))
	for _, quantile := range s.Quantiles {
		pair := strings.Split(quantile, "|")

		q, err := strconv.ParseFloat(pair[0], 64)
		helpers.PanicWithConfigError(err)

		e, err := strconv.ParseFloat(pair[1], 64)
		helpers.PanicWithConfigError(err)

		objs[q] = e
	}
	return objs
}

func (s *Summary) GetIntervals(index int) uint64 {
	return s.Schedule[index].Intervals
}

func (s *Summary) GetScheduleLength() int {
	return len(s.Schedule)
}

func (s *Summary) SetValues(index int) {
	for _, behavior := range s.Schedule[index].Behaviors {
		for _, value := range behavior.Values {
			pair := strings.Split(value, "|")
			value, err := strconv.ParseFloat(pair[0], 64)
			helpers.PanicWithConfigError(err)

			iterations, err := strconv.ParseUint(pair[1], 10, 64)
			helpers.PanicWithConfigError(err)

			for i := iterations; i > 0; i-- {
				s.SummaryVec.With(behavior.Labels).Observe(value)
			}
		}
	}
}
