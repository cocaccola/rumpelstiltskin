package types

import "github.com/prometheus/client_golang/prometheus"

type Gauge struct {
	Name     string               `yaml:"name"`
	Help     string               `yaml:"help"`
	Labels   []string             `yaml:"labels"`
	Schedule []GaugeSchedule      `yaml:"schedule"`
	GaugeVec *prometheus.GaugeVec `yaml:"-"`
}

type GaugeSchedule struct {
	Intervals uint64          `yaml:"intervals"`
	Behaviors []GaugeBehavior `yaml:"behaviors"`
}

type GaugeBehavior struct {
	Labels prometheus.Labels `yaml:"labels"`
	Value  float64           `yaml:"value"`
}

func (g *Gauge) Register(reg prometheus.Registerer, prefix string) {
	gauge := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: prefix,
			Name:      g.Name,
			Help:      g.Help,
		},
		g.Labels,
	)

	reg.MustRegister(gauge)
	g.GaugeVec = gauge
}

func (g *Gauge) GetIntervals(index int) uint64 {
	return g.Schedule[index].Intervals
}

func (g *Gauge) GetScheduleLength() int {
	return len(g.Schedule)
}

func (g *Gauge) SetValues(index int) {
	for _, behavior := range g.Schedule[index].Behaviors {
		g.GaugeVec.With(behavior.Labels).Set(behavior.Value)
	}
}
