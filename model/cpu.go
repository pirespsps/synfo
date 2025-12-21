package model

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
)

//structs

func GetCPUmodel() (string, error) { //ler com reader
	cmd := "cat /proc/cpuinfo | egrep '^model name' | uniq | awk '{print substr($0, index($0,$4))}'"
	data, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		return "", fmt.Errorf("cat model name failed: %s", cmd)
	}
	return string(data), nil
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

func GetArchitecture() (string, error) {
	cmd := `LC_ALL=C lscpu | grep Architecture | awk {'print $2}`
	data, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		return "", fmt.Errorf("error in lscpu grep Architecture: %v", err)
	}

	return string(data), nil

	//return runtime.GOARCH,nil ?
}
