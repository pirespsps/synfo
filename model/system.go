package model

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type System struct{}

type SystemData struct {
	OS            OS     `json:"os"`
	Bios          Bios   `json:"bios"`
	GUI           GUI    `json:"gui"`
	KernelVersion string `json:"kernel"`
	Graphics      string `json:"graphics"` //x11 or wayland(can be called xorg)
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

	sd.OS, err = osData()
	if err != nil {
		return Response{}, fmt.Errorf("error in os: %v", err)
	}

	sd.KernelVersion, err = kernelData()
	if err != nil {
		return Response{}, fmt.Errorf("error in kernel version: %v", err)
	}

	sd.Architecture, err = architectureData()
	if err != nil {
		return Response{}, fmt.Errorf("error in architecture: %v", err)
	}

	sd.GUI, err = guiData()
	if err != nil {
		return Response{}, fmt.Errorf("error in gui: %v", err)
	}

	sd.Graphics = graphicsData()

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

func osData() (OS, error) {

	data, err := os.ReadFile("/etc/os-release")
	if err != nil {
		return OS{}, fmt.Errorf("error in reading /etc/os-release: %v", err)
	}

	var os OS

	for v := range strings.Lines(string(data)) {
		v = strings.TrimSpace(v)
		if strings.Contains(v, "PRETTY_NAME=") {
			s := strings.ReplaceAll(v, "PRETTY_NAME=", "")
			os.Name = strings.TrimSpace(s)
		} else if strings.Contains(v, "VERSION=") {
			s := strings.ReplaceAll(v, "VERSION=", "")
			os.Release = strings.TrimSpace(s)
		}
	}

	return os, nil
}

func kernelData() (string, error) {

	name, err := exec.Command("uname", "-r").Output()
	if err != nil {
		return "", fmt.Errorf("error in kernel name cmd: %v", err)
	}

	return strings.TrimSpace(string(name)), nil
}

func architectureData() (string, error) {
	arc, err := exec.Command("uname", "-m").Output()
	if err != nil {
		return "", fmt.Errorf("error in architecture: %v", arc)
	}

	return strings.TrimSpace(string(arc)), nil
}

func guiData() (GUI, error) {

	desk := os.Getenv("XDG_CURRENT_DESKTOP")
	//wm or de....

	return GUI{
		Name: strings.TrimSpace(string(desk)),
	}, nil
}

func graphicsData() string {
	return strings.TrimSpace(os.Getenv("XDG_SESSION_TYPE"))
}
