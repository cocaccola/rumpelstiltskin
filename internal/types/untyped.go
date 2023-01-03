package types

import (
	"sync/atomic"

	"github.com/prometheus/client_golang/prometheus"
)

type UnTyped struct {
	Name         string            `yaml:"name"`
	Help         string            `yaml:"help"`
	StaticLabels prometheus.Labels `yaml:"labels"`
	Schedule     []UnTypedSchedule `yaml:"schedule"`
}

type UnTypedSchedule struct {
	Scrapes uint64  `yaml:"scrapes"`
	Value   float64 `yaml:"value"`
}

func (u *UnTyped) Register(reg prometheus.Registerer, prefix string) {
	currentScheudleIndex := &atomic.Uint64{}
	currentScheduleScrapesLeft := &atomic.Uint64{}
	currentScheduleScrapesLeft.Store(uint64(u.Schedule[currentScheudleIndex.Load()].Scrapes))
	schedules := uint64(len(u.Schedule))

	untypedFunc := prometheus.NewUntypedFunc(
		prometheus.UntypedOpts{
			Namespace:   prefix,
			Name:        u.Name,
			Help:        u.Help,
			ConstLabels: u.StaticLabels,
		},
		func() float64 {
			if currentScheduleScrapesLeft.Load() == 0 {
				currentScheudleIndex.Store((currentScheudleIndex.Load() + 1) % schedules)
				currentScheduleScrapesLeft.Store(u.Schedule[currentScheudleIndex.Load()].Scrapes)
			}
			currentScheduleScrapesLeft.Store(currentScheduleScrapesLeft.Load() - 1)
			return u.Schedule[currentScheudleIndex.Load()].Value
		},
	)

	reg.MustRegister(untypedFunc)
}
