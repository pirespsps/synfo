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

type Storage struct {
	Data []Disk
}

type Disk struct {
	Name       string      `json:"name"`
	Model      string      `json:"model"`
	Type       string      `json:"type"`
	Size       int         `json:"size"`
	Used       int         `json:"used"`
	Partitions []Partition `json:"partitions"`
}

type Partition struct {
	Name       string `json:"name"`
	Type       string `json:"type"`
	Size       int    `json:"size"`
	Used       int    `json:"used"`
	MountPoint string `json:"mountpoint"`
}

func (st Storage) Overall() (Response, error) {

	var r Response

	st, err := st.getData()
	if err != nil {
		return Response{}, fmt.Errorf("error in storage data: %v", err)
	}

	for _, v := range st.Data {
		ov := struct {
			Name           string `json:"name"`
			Type           string `json:"type"`
			Size           int    `json:"size"`
			UsedPercentage int    `json:"usedPercentage"`
		}{
			Name:           v.Name,
			Type:           v.Type,
			Size:           v.Size,
			UsedPercentage: int(float64(v.Used) / float64(v.Size) * 100),
		}
		r.Data = append(r.Data, ov)
	}

	return r, nil
}

func (st Storage) Extensive() (Response, error) {
	var r Response

	st, err := st.getData()
	if err != nil {
		return Response{}, fmt.Errorf("error in storage data: %v", err)
	}

	for _, v := range st.Data {
		r.Data = append(r.Data, v)
	}

	return r, nil
}

func (st *Storage) getData() (Storage, error) {

	wg := sync.WaitGroup{}

	entries, err := os.ReadDir("/sys/block/")
	if err != nil {
		return Storage{}, fmt.Errorf("error in reading /sys/block: %v", err)
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

	for d := range diskCh {
		disks = append(disks, d)
	}

	for e := range errCh {
		if e != nil {
			return Storage{}, fmt.Errorf("error in diskData: %v", e)
		}
	}

	return Storage{
		Data: disks,
	}, nil
}

func diskData(dir os.DirEntry) (Disk, error) {

	var disk Disk

	disk.Name = dir.Name()

	model, err := os.ReadFile(fmt.Sprintf("/sys/block/%v/device/model", disk.Name))
	if err != nil {
		return disk, fmt.Errorf("error in reading model: %v", err)
	}

	disk.Model = string(bytes.TrimSpace(model))

	isRotational, err := os.ReadFile(fmt.Sprintf("/sys/block/%v/queue/rotational", disk.Name))
	if err != nil {
		return disk, fmt.Errorf("error in reading rotational: %v", err)
	}

	isRotational = bytes.TrimSpace(isRotational)

	if string(isRotational) == "0" {
		disk.Type = "SSD"
	} else {
		disk.Type = "HDD"
	}

	size, err := os.ReadFile(fmt.Sprintf("/sys/block/%v/size", disk.Name))
	if err != nil {
		return disk, fmt.Errorf("error in reading size: %v", err)
	}

	size = bytes.TrimSpace(size)

	s, err := strconv.Atoi(string(size))
	if err != nil {
		return disk, fmt.Errorf("error in disk size convertion: %v", err)
	}

	disk.Size = s * 512 / 1000 //each sector has 512 bytes -> bytes to KBs

	disk.Partitions, err = partitionData(disk)
	if err != nil {
		return disk, fmt.Errorf("error in partitions: %v", err)
	}

	var usedSpace = 0

	for _, v := range disk.Partitions {
		usedSpace += v.Used
	}

	disk.Used = usedSpace

	return disk, nil
}

func partitionData(disk Disk) ([]Partition, error) { //just consider linux partitions
	cmd := fmt.Sprintf(`df -kT| grep %v | awk '{print $1,$2,$3,$4,$7}'`, disk.Name)
	//filtrar em go ao inves de awk (talvez)
	//não usar bash -c (ver como fazer funcionar com os pipes)
	//usar lsblk -b -o NAME,FSTYPE,SIZE,MOUNTPOINT -J | grep %v

	//cmd2 := fmt.Sprintf("lsblk -b -o NAME,FSTYPE,SIZE,MOUNTPOINT -J | grep %v", disk.Name)
	//
	//data2, err := exec.Command(cmd2).Output()
	//if err != nil {
	//	return nil, fmt.Errorf("error in cmd2: %v", err)
	//}
	//fmt.Printf("data: %v", data2)

	data, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		return nil, fmt.Errorf("error in df: %v", err)
	}

	var partitions []Partition

	for line := range strings.Lines(string(data)) {

		line = strings.TrimSpace(line)
		args := strings.Split(line, " ")

		name, systype, size, used, mount := args[0], args[1], args[2], args[3], args[4]

		s, err := strconv.Atoi(size)
		if err != nil {
			return nil, fmt.Errorf("error in partition size conversion: %v", err)
		}

		u, err := strconv.Atoi(used)
		if err != nil {
			return nil, fmt.Errorf("error in partition used conversion: %v", err)
		}

		p := Partition{
			Name:       name,
			Type:       systype,
			MountPoint: mount,
			Size:       s,
			Used:       u,
		}

		partitions = append(partitions, p)
	}

	return partitions, nil
}
