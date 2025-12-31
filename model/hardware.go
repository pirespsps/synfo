package model

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"github.com/pirespsps/synfo/utils"
)

type Hardware interface {
	Overall() (Response, error)
	Extensive() (Response, error)
}

type Response struct {
	Data []any
}

func (r Response) Print() {
	fmt.Print("\n")
	fmt.Print(utils.PrintStruct(r.Data))
	fmt.Print("\n")
}

func (r Response) PrintJson() error {
	js, err := json.Marshal(r.Data)
	if err != nil {
		return fmt.Errorf("error in marshal:%v", err)
	}

	fmt.Print(string(js))
	return nil
}

//muda isso tudo pra baixo de lugar

type hinfo struct {
	OS          string     `json:"os"`          //OS name
	OSVersion   string     `json:"osVersion"`   //OS version
	Kernel      string     `json:"kernel"`      //kernel release
	RAM         float64    `json:"ram"`         //total RAM available
	MotherBoard string     `json:"motherboard"` //motheboard name
	GraphicCard string     `json:"graphics"`    //graphic cards name
	CPU         CpuOverall `json:"cpu"`         //CPU info
	Storage     []Storage  `json:"storage"`     //storage unit name+type+size
} //muda o nome

func HardwareOverall() (hinfo, error) { //interface talvez? //MUDAR TUDO PARA SYSFS
	var data hinfo
	var err error
	data.RAM, err = RAMsize()
	if err != nil {
		return data, fmt.Errorf("error in ram size: %v", err)
	}

	data.MotherBoard, err = MotherBoardName()
	if err != nil {
		return data, fmt.Errorf("error in motherboard: %v", err)
	}

	data.CPU, err = CPUOverall()
	if err != nil {
		return data, fmt.Errorf("error in cpu overall: %v", err)
	}

	return data, nil

}

func GraphicsData() (string, error) {

	return "", nil
}

func MotherBoardName() (string, error) {

	data, err := os.ReadFile("/sys/devices/virtual/dmi/id/board_name")
	if err != nil {
		return "", fmt.Errorf("error in reading board_name: %v", err)
	}

	return string(data), nil
}

func RAMsize() (float64, error) {
	memsJ, err := exec.Command("lsmem", "-J").Output()
	if err != nil {
		return 0, fmt.Errorf("error in lsmem cmd: %v", err)
	}

	mems := struct {
		Memory []struct {
			Size string `json:"size"`
		}
	}{}

	err = json.Unmarshal(memsJ, &mems)
	if err != nil {
		return 0, fmt.Errorf("error in unmarshal mem: %v", err)
	}

	var totalSize float64
	r := regexp.MustCompile(`[^0-9,]`)

	for _, v := range mems.Memory {
		str := r.ReplaceAll([]byte(v.Size), []byte(""))
		mem := strings.ReplaceAll(string(str), ",", ".")

		fmem, err := strconv.ParseFloat(mem, 64)
		if err != nil {
			return 0, fmt.Errorf("error in parse value: %v", err)
		}
		totalSize += fmem
	}

	return totalSize, nil
}
