package model

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

type System struct{}

type SystemData struct {
	OS            OS     `json:"os"`
	Bios          Bios   `json:"bios"`
	GUI           GUI    `json:"gui"`
	KernelVersion string `json:"kernel"`
	Graphics      string `json:"graphics"` //x11 or wayland
	Architecture  string `json:"architecture"`
	Terminal      string `json:"terminal"`
}

type GUI struct {
	Name string `json:"name"`
	Type string `json:"type"` //either DE or WM
}

type OS struct {
	Name    string `json:"name"`
	Release string `json:"release"`
}

type Bios struct {
	Vendor  string `json:"vendor"`
	Version string `json:"version"`
	Release string `json:"release"`
	Date    string `json:"date"`
}

func (s System) Overall() (Response, error) {

	var sd SystemData
	var err error

	sd.Bios, err = biosData()
	if err != nil {
		return Response{}, fmt.Errorf("error in bios: %v", err)
	}

	var r Response
	r.Data = append(r.Data, sd)

	return r, nil
}

func (s System) Extensive() (Response, error) {

	return Response{}, nil
}

func biosData() (Bios, error) {

	dir := "/sys/class/dmi/id/"

	ven, err := os.ReadFile(dir + "bios_vendor")
	if err != nil {
		return Bios{}, fmt.Errorf("error in vendor: %v", err)
	}

	ver, err := os.ReadFile(dir + "bios_version")
	if err != nil {
		return Bios{}, fmt.Errorf("error in version: %v", err)
	}

	rel, err := os.ReadFile(dir + "bios_release")
	if err != nil {
		return Bios{}, fmt.Errorf("error in release: %v", err)
	}
	dt, err := os.ReadFile(dir + "bios_date")
	if err != nil {
		return Bios{}, fmt.Errorf("error in date: %v", err)
	}

	return Bios{
		Vendor:  strings.TrimSpace(string(ven)),
		Version: strings.TrimSpace(string(ver)),
		Release: strings.TrimSpace(string(rel)),
		Date:    strings.TrimSpace(string(dt)),
	}, nil
}

func osData() (name string, release string, err error) {

	data, err := os.ReadFile("/etc/os-release")
	if err != nil {
		return "", "", fmt.Errorf("error in reading /etc/os-release: %v", err)
	}

	rn := regexp.MustCompile(`PRETTY_NAME=`)
	rv := regexp.MustCompile(`VERSION=`)

	for v := range bytes.Lines(data) {
		if rn.Match(v) {
			name = string(rn.ReplaceAll(v, []byte("")))

		} else if rv.Match(v) {
			release = string(rv.ReplaceAll(v, []byte("")))
		}
	}

	return name, release, nil
}

func kernelData() (string, error) {

	name, err := exec.Command("uname", "-r").Output()
	if err != nil {
		return "", fmt.Errorf("error in kernel name cmd: %v", err)
	}

	return string(name), nil
}
