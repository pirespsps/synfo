package model

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strconv"
	"strings"
)

type CPU struct {
	Data   CpuData
	Manage CpuManageData
}

type CpuManageData struct {
	UsagePerCent int `json:"usage"`
	TemperatureC int `json:"temperatureC"` //temperature in Celsius
	//etc....
}

type CpuData struct {
	Name      string            `json:"name"` //model
	Producer  string            `json:"producer"`
	Cores     []CoreData        `json:"cores"`
	Threads   int               `json:"threads"`
	Frequency string            `json:"frequency"`
	Arch      string            `json:"architecture"`
	Cache     map[string]string `json:"cache"`
	VendorId  string            `json:"vendor"`
}

type CoreData struct {
	Id           int     `json:"coreID"`    //core id (repeat in physical and virtual)
	Processor    int     `json:"processor"` //processor id (can't repeat)
	FrequencyMHz float64 `json:"frequency"`
	VendorId     string  `json:"vendor"`
}

func (c *CPU) Overall() ([]byte, error) {

	return nil, nil
}

func (c *CPU) Extensive() ([]byte, error) {

	return nil, nil
}

func (c *CPU) Load() error {
	return nil
}

func getCoreData(cores []CoreData) ([]CoreData, error) { //mudar para sysfs

	data, err := os.ReadFile("/proc/cpuinfo")
	if err != nil {
		return nil, fmt.Errorf("error in reading /proc/cpuinfo: %v", err)
	}

	core := CoreData{}

	for v := range strings.Lines(string(data)) {

		v = strings.TrimSpace(v)

		if strings.HasPrefix(v, "core id") {
			_, af, _ := strings.Cut(v, ":")
			af = strings.TrimSpace(af)

			core.Id, err = strconv.Atoi(af)
			if err != nil {
				return nil, fmt.Errorf("error in parsing core id: %v", err)
			}

		} else if strings.HasPrefix(v, "cpu MHz") {
			_, af, _ := strings.Cut(v, ":")
			af = strings.TrimSpace(af)

			core.FrequencyMHz, err = strconv.ParseFloat(af, 64)
			if err != nil {
				return nil, fmt.Errorf("error in parsing frequency: %v", err)
			}

		} else if strings.HasPrefix(v, "vendor_id") {
			_, af, _ := strings.Cut(v, ":")
			core.VendorId = af

		} else if strings.HasPrefix(v, "processor") {
			_, af, _ := strings.Cut(v, ":")
			af = strings.TrimSpace(af)

			core.Processor, err = strconv.Atoi(af)
			if err != nil {
				return nil, fmt.Errorf("error in parsing processor: %v", err)
			}
		}

		if v == "" {
			cores = append(cores, core)

			core = CoreData{}
			continue
		}
	}

	return cores, nil
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

func getThreads() int {
	return runtime.NumCPU()
}

func getArchitecture() string {
	return runtime.GOARCH
}
