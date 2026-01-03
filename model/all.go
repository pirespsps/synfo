package model

import (
	"fmt"
	"os"
)

//instead of hardware, make all overall and extensive

func motherBoardData() (string, error) {

	data, err := os.ReadFile("/sys/devices/virtual/dmi/id/board_name")
	if err != nil {
		return "", fmt.Errorf("error in reading board_name: %v", err)
	}

	return string(data), nil
}
