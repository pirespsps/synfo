package model

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Storage struct {
	Disks []Disk `json:"disks"`
}

type Disk struct {
	Name       string      `json:"name"`
	Type       string      `json:"type"`
	Size       int         `json:"size"`
	Used       int         `json:"used"`
	Partitions []Partition `json:"partitions"`
}

type Partition struct {
	Name       string `json:"name"`
	Size       int    `json:"size"`
	Used       int    `json:"used"`
	MountPoint string `json:"mountpoint"`
}

func (st Storage) Overall() OverallData {

	return OverallData{}
}

func (st Storage) Extensive() ExtensiveData {

	return ExtensiveData{}
}

func getStorage() ([]Disk, error) {

	var disks []Disk

	entries, err := os.ReadDir("/sys/block/")
	if err != nil {
		return nil, fmt.Errorf("error in reading /sys/block: %v", err)
	}

	for _, e := range entries {
		if strings.HasPrefix(e.Name(), "sd") {
			go getDiskData(e) //...
		}
	}

	//waitgroup...

	return nil, nil
}

func getDiskData(dir os.DirEntry) (Disk, error) {
	var disk Disk

	disk.Name = dir.Name()

	size, err := os.ReadFile(fmt.Sprint("/sys/block/%v/size", dir.Name()))
	if err != nil {
		return disk, fmt.Errorf("error in reading size: %v", err)
	}

	disk.Size, err = strconv.Atoi(string(size))
	if err != nil {
		return disk, fmt.Errorf("error in disk size convertion: %v", err)
	}

	return disk, nil
}
