package model

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

type StorageUnit struct {
	Name string `json:"name"` //disk name
	Type string `json:"type"` //disk type
	Size string `json:"size"` //disk size
}

type CpuInfo struct {
	Name  string             `json:"name"`  //name of the CPU
	Cores int                `json:"cores"` //physical cores count
	Arq   string             `json:"arq"`   //architecture
	Cache map[string]float32 `json:"cache"` //caches sizes
}

type hinfo struct {
	OS          string        `json:"os"`          //OS name
	OSVersion   string        `json:"osVersion"`   //OS version
	Kernel      string        `json:"kernel"`      //kernel release
	RAM         float64       `json:"ram"`         //total RAM available
	MotherBoard string        `json:"motherboard"` //motheboard name
	GraphicCard string        `json:"graphics"`    //graphic cards name
	CPU         CpuInfo       `json:"cpu"`
	Storage     []StorageUnit `json:"storage"`
}

func GetAll() (hinfo, error) {
	var data hinfo
	var err error

	data.OS, data.OSVersion, err = OsData()
	if err != nil {
		return data, fmt.Errorf("error in osData: %v", err)
	}

	data.Kernel, err = KernelName()
	if err != nil {
		return data, fmt.Errorf("error in kernel name: %v", err)
	}

	data.RAM, err = RAMsize()
	if err != nil {
		return data, fmt.Errorf("error in ram size: %v", err)
	}

	data.MotherBoard, err = MotherBoardName()
	if err != nil {
		return data, fmt.Errorf("error in motherboard: %v", err)
	}

	data.Storage, err = StorageData()
	if err != nil {
		return data, fmt.Errorf("error in storage data: %v", err)
	}

	return data, nil

}

func GraphicsData() (string, error) {

	return "", nil
}

func StorageData() ([]StorageUnit, error) {

	disks := struct {
		BlockDevices []StorageUnit `json:"blockdevices"`
	}{}

	data, err := exec.Command("lsblk", "-J").Output()
	if err != nil {
		return nil, fmt.Errorf("error in lsblk: %v", err)
	}

	err = json.Unmarshal(data, &disks)
	if err != nil {
		return nil, fmt.Errorf("error in unmarshal: %v", err)
	}

	return disks.BlockDevices, nil
}

func MotherBoardName() (string, error) { //ler com reader

	data, err := os.ReadFile("/sys/devices/virtual/dmi/id/board_name")
	if err != nil {
		return "", fmt.Errorf("error in reading board_name: %v", err)
	}

	return string(data), nil
}

func OsData() (name string, release string, err error) { //ler com reader

	data, err := exec.Command("lsb_release", "-a").Output()
	if err != nil {
		return "", "", fmt.Errorf("error in lsb_release: %v", err)
	}

	for v := range strings.Lines(string(data)) {

		if strings.Contains(v, "Release") {

			r := regexp.MustCompile(`[a-zA-Z]+:[\s]+`)
			release = string(r.ReplaceAll([]byte(v), []byte("")))
			continue

		} else if strings.Contains(v, "Description") {

			r := regexp.MustCompile(`[a-zA-Z]+:[\s]+`)
			name = string(r.ReplaceAll([]byte(v), []byte("")))

		}

	}

	return name, release, nil
}

func KernelName() (string, error) {

	name, err := exec.Command("uname", "-r").Output()
	if err != nil {
		return "", fmt.Errorf("error in kernel name cmd: %v", err)
	}

	return string(name), nil
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
