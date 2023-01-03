package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/cocaccola/rumpelstiltskin/internal/config"
	"github.com/cocaccola/rumpelstiltskin/internal/metrics"
)

var (
	version        string
	versionFlag    bool
	scrapeInterval time.Duration
	metricPrefix   string
	port           uint64
	configFile     string
)

func main() {
	flag.DurationVar(&scrapeInterval, "interval", 10*time.Second, "interval that metrics will be updated")
	flag.StringVar(&metricPrefix, "prefix", "", "prefix for the metric names")
	flag.Uint64Var(&port, "port", 9090, "port to listen on")
	flag.StringVar(&configFile, "config", "config.yaml", "path to config file")
	flag.BoolVar(&versionFlag, "version", false, "print version info")
	flag.Parse()

	if versionFlag {
		fmt.Println(version)
		os.Exit(0)
	}

	c, err := config.ParseConfig(configFile)
	if err != nil {
		panic(err)
	}

	metrics.Register(c, metricPrefix)

	metrics.SpawnWorkers(scrapeInterval, c)

	err = metrics.StartPromServer(port)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(2)
	}
}
