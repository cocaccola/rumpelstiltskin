package types

import "github.com/prometheus/client_golang/prometheus"

type Counter struct {
	Name       string                 `yaml:"name"`
	Help       string                 `yaml:"help"`
	Labels     []string               `yaml:"labels"`
	Schedule   []CounterSchedule      `yaml:"schedule"`
	CounterVec *prometheus.CounterVec `yaml:"-"`
}

type CounterSchedule struct {
	Intervals uint64            `yaml:"intervals"`
	Behaviors []CounterBehavior `yaml:"behaviors"`
}

type CounterBehavior struct {
	Labels prometheus.Labels `yaml:"labels"`
	Add    float64           `yaml:"add"`
}

func (c *Counter) Register(reg prometheus.Registerer, prefix string) {
	counter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: prefix,
			Name:      c.Name,
			Help:      c.Help,
		},
		c.Labels,
	)

	reg.MustRegister(counter)
	c.CounterVec = counter
}

func (c *Counter) GetIntervals(index int) uint64 {
	return c.Schedule[index].Intervals
}

func (c *Counter) GetScheduleLength() int {
	return len(c.Schedule)
}

func (c *Counter) SetValues(index int) {
	for _, behavior := range c.Schedule[index].Behaviors {
		c.CounterVec.With(behavior.Labels).Add(behavior.Add)
	}
}
