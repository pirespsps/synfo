package parser

import (
	"errors"

	"github.com/pirespsps/synfo/model"
)

func GetResponse(comp string, option string) ([]byte, error) {

	switch comp {

	case "all":

	case "cpu":

	case "ram":

	case "storage":
		return storageInfo(option)

	case "process": //fazer programa separado

	case "network":

	case "system":
		return systemInfo(option)

	default:
		return nil, errors.New("option doesnt exist")
	}

	return nil, nil
}

func systemInfo(option string) ([]byte, error) {
	var sys model.System

	if option == "extensive" {
		return sys.Extensive()
	}

	return sys.Overall()
}

func storageInfo(option string) ([]byte, error) {

	var storage model.Storage

	if option == "extensive" {
		return storage.Extensive()
	}

	return storage.Overall()
}
