package model

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
)

type Storage struct{}

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

func (st Storage) Overall() (any, error) {

	disks, err := getStorage()
	if err != nil {
		return nil, fmt.Errorf("error in storage data: %v", err)
	}

	data := struct {
		Storage []Disk `json:"storage"`
	}{
		Storage: disks,
	}

	return data, nil
}

func (st Storage) Extensive() (any, error) {
	return nil, nil
}

func getStorage() ([]Disk, error) {

	wg := sync.WaitGroup{}

	entries, err := os.ReadDir("/sys/block/")
	if err != nil {
		return nil, fmt.Errorf("error in reading /sys/block: %v", err)
	}

	diskCh := make(chan (Disk))
	errCh := make(chan (error))

	for _, e := range entries {

		if strings.HasPrefix(e.Name(), "sd") {

			wg.Add(1)
			go func(e os.DirEntry) {
				defer wg.Done()
				if disk, err := diskData(e); err != nil {
					errCh <- err
				} else {
					diskCh <- disk
				}
			}(e)

		}
	}

	go func() {
		wg.Wait()
		close(diskCh)
		close(errCh)
	}()

	var disks []Disk

	go func() {
		for d := range diskCh {
			disks = append(disks, d)
		}
	}()

	for e := range errCh {
		if e != nil {
			return disks, fmt.Errorf("error in diskData: %v", e)
		}
	}

	return disks, nil
}

func diskData(dir os.DirEntry) (Disk, error) {

	var disk Disk

	disk.Name = dir.Name()

	size, err := os.ReadFile(fmt.Sprintf("/sys/block/%v/size", disk.Name))
	if err != nil {
		return disk, fmt.Errorf("error in reading size: %v", err)
	}

	size = bytes.TrimSpace(size)

	disk.Size, err = strconv.Atoi(string(size))
	if err != nil {
		return disk, fmt.Errorf("error in disk size convertion: %v", err)
	}

	cmd := fmt.Sprintf(`df | grep %v | awk '{print $1,$2,$3,$6}'`, disk.Name)

	data, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		return disk, fmt.Errorf("error in df: %v", err)
	}

	for line := range strings.Lines(string(data)) {

		line = strings.TrimSpace(line)
		args := strings.Split(line, " ")
		name, size, used, mount := args[0], args[1], args[2], args[3]

		s, err := strconv.Atoi(size)
		if err != nil {
			return disk, fmt.Errorf("error in partition size conversion: %v", err)
		}

		u, err := strconv.Atoi(used)
		if err != nil {
			return disk, fmt.Errorf("error in partition used conversion: %v", err)
		}

		disk.Partitions = append(disk.Partitions, Partition{
			Name:       name,
			MountPoint: mount,
			Size:       s,
			Used:       u,
		})
	}

	return disk, nil
}
