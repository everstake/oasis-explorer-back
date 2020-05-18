package services

import (
	"fmt"
	"io/ioutil"
)

const metricsPath = "./resources/metrics.json"

func (s ServiceFacade) GetMetricsConfig() ([]byte, error) {
	data, err := ioutil.ReadFile(metricsPath)
	if err != nil {
		return nil, fmt.Errorf("ioutil.ReadAll: %s", err.Error())
	}
	return data, nil
}
