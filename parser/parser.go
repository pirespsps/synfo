package parser

import (
	"errors"

	"github.com/pirespsps/synfo/model"
)

func FetchData(option string) (any, error) {

	//get option from main -> fetch -> parse to obj -> return

	switch option {

	case "all":

	case "cpu":

	case "ram":

	case "storage":
		return storageInfo()

	case "process": //fazer programa separado

	case "hardware":
		return hardwareInfo()

	case "network":

	case "system":

	default:
		return nil, errors.New("option doesnt exist")
	}

	return nil, nil
}

func hardwareInfo() (any, error) {

	//call from hardware

	return model.CPUData()
}

func storageInfo() (any, error) {
	var storage model.Storage

	return storage.Overall()
}
