package main

import (
	"errors"
	"fmt"
	"github.com/blockblu-io/cardano-node-healthcheck/config"
	ctime "github.com/godano/cardano-lib/time"
	"math/big"
	"time"
)

// getTimeSettings is extracting the blockchain time settings from the given genesis
// file. If the genesis file couldn't be read and parsed, then an error will be
// returned. Otherwise, the extracted time settings will be returned.
func getTimeSettings(genesisFilePath string) (*ctime.TimeSettings, error) {
	genesis, err := config.ParseGenesis(genesisFilePath)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("the genesis at '%s' cannot be parsed: %s\n",
			genesisFilePath, err.Error()))
	}
	timeSettings := &ctime.TimeSettings{
		GenesisBlockDateTime: genesis.GenesisBlockCreationTime,
		SlotsPerEpoch:        new(big.Int).SetUint64(genesis.SlotsPerEpoch),
		SlotDuration:         time.Duration(genesis.SlotDurationInS) * time.Second,
	}
	return timeSettings, nil
}

// getPrometheusURL is assembling the URL for the Prometheus endpoint of the cardano-node
// by reading the given configuration file of the node. If the configuration file couldn't
// be parsed, then an error will be returned. Otherwise, the assembled Prometheus URL will
// be returned.
func getPrometheusURL(configFilePath string) (string, error) {
	nodeConfig, err := config.ParseNodeConfig(configFilePath)
	if err != nil {
		return "", errors.New(fmt.Sprintf("the node configuration at '%s' cannot be parsed: %s\n", configFilePath,
			err.Error()))
	}
	return fmt.Sprintf("http://%v:%v/metrics", nodeConfig.Prometheus[0], nodeConfig.Prometheus[1]), nil
}
