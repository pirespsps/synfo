package model

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

type Hinfo struct {
	OS          string  `json:"os"`          //OS name
	OSVersion   string  `json:"osVersion"`   //OS version
	Kernel      string  `json:"kernel"`      //kernel release
	RAM         float64 `json:"ram"`         //total RAM available
	MotherBoard string  `json:"motherboard"` //motheboard name
	CPU         struct {
		Name  string             `json:"name"`  //name of the CPU
		Cores int                `json:"cores"` //physical cores count
		Arq   string             `json:"arq"`   //architecture
		Cache map[string]float32 `json:"cache"` //caches sizes
	}
	Storage struct {
		Name     string  `json:"name"`     //disk name
		Type     string  `json:"type"`     //disk type
		Capacity float32 `json:"capacity"` //disk size
	}
}

func GetAll() {

}

func OsData() (name string, release string, err error) {

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

	for _, v := range mems.Memory {
		mem := strings.ReplaceAll(v.Size, "G", "")
		mem = strings.ReplaceAll(mem, ",", ".")

		fmem, err := strconv.ParseFloat(mem, 64)
		if err != nil {
			return 0, fmt.Errorf("error in parse value: %v", err)
		}
		totalSize += fmem
	}

	return totalSize, nil
}
