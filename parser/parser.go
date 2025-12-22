package parser

import (
	"errors"

	"github.com/pirespsps/synfo/model"
	"github.com/shirou/gopsutil/v4/cpu"
)

func FetchData(option string, isJSON bool) (any, error) {

	//get option from main -> fetch -> parse to obj -> return

	switch option {

	case "all":

	case "cpu":
		return cpuInfo(isJSON)
	case "ram":

	case "storage":

	case "process": //fazer programa separado

	case "hardware":
		return hardwareInfo(isJSON)

	case "network":

	case "system":

	default:
		return nil, errors.New("option doesnt exist")
	}

	return nil, nil
}

func cpuInfo(isJSON bool) (any, error) {

	cores, err := cpu.Counts(false)
	if err != nil {
		return nil, errors.New("error in core count")
	}

	used, err := cpu.Percent(1000, false)
	if err != nil {
		return nil, errors.New("error in percentage use")
	}

	data := struct {
		Used  float64
		Cores int
	}{
		Used:  used[0],
		Cores: cores,
	}

	return data, nil
}

func hardwareInfo(isJson bool) (any, error) {

	//call from hardware
	return model.GetCPUmodel()
	return nil, nil
}
