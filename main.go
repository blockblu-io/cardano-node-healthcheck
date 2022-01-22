package main

import (
	"flag"
	"fmt"
	"github.com/blockblu-io/cardano-node-healthcheck/health"
	log "github.com/sirupsen/logrus"
	"os"
	"time"
)

// printUsage prints a simple usage message for this application.
func printUsage(args []string) {
	var appName string
	if len(args) == 0 || args[0] == "" {
		appName = "healthcheck"
	} else {
		appName = args[0]
	}
	fmt.Printf("\nUsage:\n\t%s [flags]\nFlags:\n", appName)
	flag.PrintDefaults()
}

// main is the entry point for this application.
func main() {
	configFilePath := flag.String("config-file", "", "path to the configuration file of cardano-node")
	genesisFilePath := flag.String("genesis-file", "", "path to the genesis file of cardano-node")
	maxTimeSinceLastBlock := flag.Duration("max-time-since-last-block", 10*time.Minute,
		"threshold for duration between now and the creation date of the most recently received block")
	flag.Parse()

	if *configFilePath == "" {
		log.Error("configuration file hasn't been specified")
		printUsage(os.Args)
		os.Exit(1)
	}
	if *genesisFilePath == "" {
		log.Error("genesis file hasn't been specified")
		printUsage(os.Args)
		os.Exit(1)
	}

	timeSettings, err := getTimeSettings(*genesisFilePath)
	if err != nil {
		log.Errorf("error: %s", err.Error())
		os.Exit(1)
	}
	prometheusURL, err := getPrometheusURL(*configFilePath)
	if err != nil {
		log.Errorf("error: %s", err.Error())
		os.Exit(1)
	}
	cfg := health.Config{
		PrometheusURL:         prometheusURL,
		TimeSettings:          *timeSettings,
		MaxTimeSinceLastBlock: *maxTimeSinceLastBlock,
	}
	healthy, err := health.Check(cfg)
	if err == nil {
		if *healthy {
			log.Info("node is healthy")
			os.Exit(0)
		} else {
			log.Info("node isn't healthy")
			os.Exit(1)
		}
	} else {
		log.Errorf("error: %s", err.Error())
		os.Exit(1)
	}
}
