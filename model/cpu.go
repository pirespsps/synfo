package model

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strings"
)

type CpuOverall struct {
	Name  string            `json:"name"`         //name of the CPU
	Cores int               `json:"cores"`        //physical cores count
	Arch  string            `json:"architecture"` //architecture
	Cache map[string]string `json:"cache"`        //caches sizes
}

type CpuData struct {
	Name         string            `json:"name"`
	Producer     string            `json:"producer"`
	Cores        int               `json:"cores"`
	Threads      int               `json:"threads"`
	Frequency    string            `json:"frequency"`
	Arch         string            `json:"architecture"`
	Cache        map[string]string `json:"cache"`
	VendorId     string            `json:"vendor"`
	UsagePerCent int               `json:"usage"`
	TemperatureC int               `json:"temperatureC"`
}

func CPUOverall() (CpuOverall, error) {

	var o CpuOverall
	var err error

	o.Arch = getArchitecture()
	o.Cores = getCores()

	o.Name, err = getCPUmodel()
	if err != nil {
		return o, fmt.Errorf("error in cpu model: %v", err)
	}
	o.Cache, err = getCache()
	if err != nil {
		return o, fmt.Errorf("error in cache: %v", err)
	}

	return o, nil
}

func CPUData() (CpuData, error) {

	var d CpuData

	return d, nil
}

func getCPUmodel() (string, error) {

	data, err := os.ReadFile("/proc/cpuinfo")
	if err != nil {
		return "", fmt.Errorf("error in reading /proc/cpuinfo: %v", err)
	}

	r := regexp.MustCompile(`model name ?[\s]+:\s`)
	var model string

	for v := range bytes.Lines(data) {
		if r.Match([]byte(v)) {
			model = string(r.ReplaceAll([]byte(v), []byte("")))
			break
		}
	}

	return model, nil
}

func getCache() (map[string]string, error) {
	cmd := `LC_ALL=C lscpu | grep cache | awk {'print $1 ":" $3 " " $4'}`
	data, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		return nil, fmt.Errorf("error in lscpu grep cache: %v", err)
	}

	var cmap = make(map[string]string)

	for v := range bytes.Lines(data) {
		b, a, _ := strings.Cut(string(v), ":")
		cmap[b] = a
	}

	return cmap, nil
}

func getCores() int {
	return runtime.NumCPU()
}

func getArchitecture() string {
	return runtime.GOARCH
}
