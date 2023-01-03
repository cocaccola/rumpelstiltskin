package metrics

import (
	"log"
	"reflect"
	"time"

	"github.com/cocaccola/rumpelstiltskin/internal/config"
	"github.com/cocaccola/rumpelstiltskin/internal/types"
)

func Worker(scrapeInterval time.Duration, metric Metric) {
	currentScheudleIndex := 0
	currentScheduleIntervalsLeft := metric.GetIntervals(currentScheudleIndex)
	schedules := metric.GetScheduleLength()

	for range time.Tick(scrapeInterval) {
		if currentScheduleIntervalsLeft == 0 {
			currentScheudleIndex = (currentScheudleIndex + 1) % schedules
			currentScheduleIntervalsLeft = metric.GetIntervals(currentScheudleIndex)
		}

		metric.SetValues(currentScheudleIndex)

		currentScheduleIntervalsLeft--
	}
}

func SpawnWorkers(scrapeInterval time.Duration, config *config.Config) {
	v := reflect.ValueOf(config)

	for i := 0; i < v.Elem().NumField(); i++ {
		for j := 0; j < v.Elem().Field(i).Len(); j++ {
			metric := v.Elem().Field(i).Index(j).Interface()
			if _, ok := metric.(*types.UnTyped); ok {
				continue
			}
			log.Printf("spawning worker for %s\n", v.Elem().Field(i).Index(j).Elem().FieldByName("Name"))
			go Worker(scrapeInterval, metric.(Metric))
		}
	}
}
