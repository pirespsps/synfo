package parser

import (
	"errors"

	"github.com/pirespsps/synfo/model"
)

func GetResponse(comp string, option string) (model.Response, error) {

	switch comp {

	case "all":

	case "cpu":

	case "ram":

	case "storage":
		return storageInfo(option)

	case "process": //fazer programa separado

	case "hardware":
		return hardwareInfo()

	case "network":

	case "system":
		return systemInfo()

	default:
		return model.Response{}, errors.New("option doesnt exist")
	}

	return model.Response{}, nil
}

func systemInfo() (model.Response, error) {
	var sys model.System
	return sys.Overall()
}

func hardwareInfo() (model.Response, error) {

	//call from hardware

	//return model.CPUData()
	return model.Response{}, nil
}

func storageInfo(option string) (model.Response, error) {

	var storage model.Storage

	if option == "extensive" {
		return storage.Extensive()
	}

	return storage.Overall()
}
