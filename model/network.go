package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
)

type Network struct {
	Data NetworkData
}

type NetworkData struct {
	SystemHost string `json:"system"`
	Hosts      []Host `json:"hosts"`
}

type Host struct {
	IP     string `json:"ip"`
	Domain string `json:"domain"`
	Name   string `json:"name"`
}

func (n *Network) Overall() ([]byte, error) {

	if err := n.Load(); err != nil {
		return nil, fmt.Errorf("error in load: %v", err)
	}

	by, err := json.Marshal(n)
	if err != nil {
		return nil, fmt.Errorf("error in marshal: %v", by)
	}

	return by, nil

}

func (n *Network) Extensive() ([]byte, error) {

	if err := n.Load(); err != nil {
		return nil, fmt.Errorf("error in load: %v", err)
	}

	by, err := json.Marshal(n)
	if err != nil {
		return nil, fmt.Errorf("error in marshal: %v", by)
	}

	return by, nil

}

func (n *Network) Load() error {

	if n == nil {
		return errors.New("instance is nil")
	}

	var err error

	n.Data.Hosts, err = getHosts()
	if err != nil {
		return fmt.Errorf("error in hosts: %v", err)
	}

	return nil
}

func getHosts() ([]Host, error) {

	var hosts []Host

	data, err := os.ReadFile("/etc/hosts")
	if err != nil {
		return nil, fmt.Errorf("error opening file:%v", err)
	}

	for line := range strings.Lines(string(data)) {

		line = strings.TrimSpace(line)

		if strings.HasPrefix(line, "#") {
			continue
		} else if len(line) < 1 {
			continue
		}

		val := strings.Fields(line)

		if len(val) >= 3 {
			hosts = append(hosts, Host{
				IP:     val[0],
				Domain: val[1],
				Name:   strings.Join(val[2:], " "),
			})
		}
	}

	return hosts, nil
}
