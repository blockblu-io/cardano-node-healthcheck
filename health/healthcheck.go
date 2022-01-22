package health

import (
	"bufio"
	"errors"
	"fmt"
	ctime "github.com/godano/cardano-lib/time"
	"io"
	"math/big"
	"net/http"
	"strings"
	"time"
)

// Config contains the required information for performing a health check for a certain cardano-node instance.
type Config struct {
	PrometheusURL         string
	TimeSettings          ctime.TimeSettings
	MaxTimeSinceLastBlock time.Duration
	MinPeerConnections    int
}

// buildMap is creating a map of keys and values from the body given by a request to
// the Prometheus endpoint of a cardano-node. An error will be returned, if the body
// couldn't be parsed properly. Otherwise, the key and value map will be returned.
func buildMap(r io.Reader) (map[string]string, error) {
	pMap := make(map[string]string)
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		lineArr := strings.Split(line, " ")
		if len(lineArr) == 2 {
			pMap[lineArr[0]] = lineArr[1]
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, errors.New(fmt.Sprintf("error reading response from Prometheus: '%s'", scanner.Err().Error()))
	}
	return pMap, nil
}

// checkMaxTimeSinceLastBlock checks the value pair stating the slot number
// of the most recent received block, and computes the current epoch and slot
// from it. This function checks whether the slot is behind more than the given
// 'maxTimeSinceLastBlock' duration, and returns false if this is the case. True
// is returned, if the most recent block was in the time frame. An epoch will be
// returned if the corresponding value pair cannot be parsed correctly.
func checkMaxTimeSinceLastBlock(maxTimeSinceLastBlock time.Duration, pMap map[string]string,
	settings ctime.TimeSettings) (*bool, error) {
	slotNumString, foundSlotNum := pMap["cardano_node_metrics_slotNum_int"]
	if !foundSlotNum {
		return nil, errors.New("could not find the correct information (slotNum) from the Prometheus endpoint")
	}
	slotNum, validSlotNum := new(big.Int).SetString(slotNumString, 10)
	if !validSlotNum {
		return nil, errors.New("slotNum in Prometheus endpoint isn't a valid integer")
	}
	epoch := new(big.Int).Div(slotNum, settings.SlotsPerEpoch)
	slot := new(big.Int).Mod(slotNum, settings.SlotsPerEpoch)
	slotDate, err := ctime.FullSlotDateFrom(epoch, slot, settings)
	if err != nil {
		panic(fmt.Sprintf("epoch/slot date does not match blockchain details: %s", err.Error()))
	}
	healthy := time.Now().Sub(slotDate.GetEndDateTime()) <= maxTimeSinceLastBlock
	return &healthy, nil
}

// Check checks whether the node with the given Config is healthy. An error will be returned,
// if the request to the Prometheus endpoint failed or the request couldn't be parsed.
// Otherwise, true will be returned, if the node is healthy, and false, if the node is
// unhealthy.
func Check(config Config) (*bool, error) {
	client := http.Client{
		Timeout: 5 * time.Second,
	}
	response, err := client.Get(config.PrometheusURL)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("not able to reach prometheus endpoint at '%s':%s",
			config.PrometheusURL, err.Error()))
	}
	if 200 <= response.StatusCode && response.StatusCode < 300 {
		var pMap, err = buildMap(response.Body)
		if err != nil {
			return nil, err
		}
		return checkMaxTimeSinceLastBlock(config.MaxTimeSinceLastBlock, pMap, config.TimeSettings)
	} else {
		return nil, errors.New(fmt.Sprintf("prometheus endpoint reported status code '%s'", response.Status))
	}
}
