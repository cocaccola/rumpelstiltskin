package config

import (
	"io"
	"log"
	"os"
	"reflect"

	"github.com/cocaccola/rumpelstiltskin/internal/types"
	"github.com/prometheus/client_golang/prometheus"
	yaml "gopkg.in/yaml.v3"
)

type Registerer interface {
	Register(prometheus.Registerer, string)
}

type Config struct {
	Counters   []*types.Counter   `yaml:"counters"`
	Gauges     []*types.Gauge     `yaml:"gauges"`
	Histograms []*types.Histogram `yaml:"histograms"`
	Summaries  []*types.Summary   `yaml:"summaries"`
	UnTyped    []*types.UnTyped   `yaml:"untyped"`
}

func (c *Config) Register(reg prometheus.Registerer, prefix string) {
	v := reflect.ValueOf(c)

	for i := 0; i < v.Elem().NumField(); i++ {
		for j := 0; j < v.Elem().Field(i).Len(); j++ {
			v.Elem().Field(i).Index(j).Interface().(Registerer).Register(reg, prefix)
			log.Printf("registered %s\n", v.Elem().Field(i).Index(j).Elem().FieldByName("Name"))
		}
	}
}

func ParseConfig(file string) (*Config, error) {
	f, err := os.OpenFile(file, os.O_RDONLY, 0)
	if err != nil {
		return nil, err
	}

	data, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	config := &Config{}

	err = yaml.Unmarshal(data, config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
