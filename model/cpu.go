package model

import (
	"bytes"
	"encoding/json"
	"errors"
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
	if err := c.Load(); err != nil {
		return nil, fmt.Errorf("error in load: %v", err)
	}

	bt, err := json.Marshal(c.Data)
	if err != nil {
		return nil, fmt.Errorf("error marsheling: %v", err)
	}

	return bt, nil
}

func (c *CPU) Extensive() ([]byte, error) {
	if err := c.Load(); err != nil {
		return nil, fmt.Errorf("error in load: %v", err)
	}
	return nil, nil
}

func (c *CPU) Load() error {
	if c == nil {
		return errors.New("instance is nil")
	}
	var err error

	c.Data.Arch = getArchitecture()
	c.Data.Threads = getThreads()
	c.Data.Name, err = getCPUmodel()
	if err != nil {
		return fmt.Errorf("error in cpu model: %v", err)
	}

	c.Data.Cores, err = getCoreData()
	if err != nil {
		return fmt.Errorf("error in cores: %v", err)
	}

	c.Data.Cache, err = getCache()
	if err != nil {
		return fmt.Errorf("error in cache: %v", err)
	}

	c.Data.Frequency, err = getFrequency()
	if err != nil {
		return fmt.Errorf("error in frequency: %v", err)
	}

	c.Data.VendorId, err = getVendor()
	if err != nil {
		return fmt.Errorf("error in vendor: %v", err)
	}

	return nil
}

func getCoreData() ([]CoreData, error) {

	data, err := os.ReadFile("/proc/cpuinfo")
	if err != nil {
		return nil, fmt.Errorf("error in reading /proc/cpuinfo: %v", err)
	}

	var cores []CoreData

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

func getFrequency() (string, error) {
	data, err := os.ReadFile("/sys/devices/system/cpu/cpu0/cpufreq/cpuinfo_max_freq")
	if err != nil {
		return "", fmt.Errorf("error in reading file: %v", err)
	}

	data = bytes.TrimSpace(data)

	freq, err := strconv.ParseFloat(string(data), 64)
	if err != nil {
		return "", fmt.Errorf("error in converting to int: %v", err)
	}

	freq = freq / 1000_000 //kb to gb

	freqS := fmt.Sprintf("%.1f GHz", freq)

	return freqS, nil
}

func getVendor() (string, error) {

	data, err := os.ReadFile("/proc/cpuinfo")
	if err != nil {
		return "", fmt.Errorf("error opening file: %v", err)
	}

	for line := range strings.Lines(string(data)) {
		if strings.Contains(line, "vendor_id") {
			str := strings.ReplaceAll(line, "vendor_id", "")
			str = strings.ReplaceAll(str, ":", "")
			str = strings.TrimSpace(str)

			return str, nil
		}
	}

	return "", errors.New("vendor not found")
}
