package model

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strings"
)

//structs

func GetCPUmodel() (string, error) {

	data, err := os.ReadFile("/proc/cpuinfo")
	if err != nil {
		return "", fmt.Errorf("error in reading /proc/cpuinfo: %v", err)
	}

	r := regexp.MustCompile(`model name ?[\s]+:\s`)
	var model string

	for v := range strings.Lines(string(data)) {
		if r.Match([]byte(v)) {
			model = string(r.ReplaceAll([]byte(v), []byte("")))
			break
		}
	}

	return model, nil
}

func GetCache() (map[string]string, error) {
	cmd := `LC_ALL=C lscpu | grep cache | awk {'print $1 ":" $3 " " $4'}`
	data, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		return nil, fmt.Errorf("error in lscpu grep cache: %v", err)
	}

	var cmap = make(map[string]string)

	for v := range strings.Lines(string(data)) {
		b, a, _ := strings.Cut(v, ":")
		cmap[b] = a
	}

	return cmap, nil
}

func GetCores() int {
	return runtime.NumCPU()
}

func GetArchitecture() string {
	return runtime.GOARCH
}
