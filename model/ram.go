package model

import (
	"errors"
	"fmt"
)

type RAM struct {
	Data []RAMData
}

type RAMData struct {
	Name         string  `json:"name"`
	CapacityG    float64 `json:"capacityG"`
	DDR          string  `json:"ddr"`
	FrequencyMHz int     `json:"frequencyMHz"`
}

func (r *RAM) Overall() ([]byte, error) {
	if err := r.load(); err != nil {
		return nil, fmt.Errorf("error in load:%v", err)
	}

	return nil, nil
}

func (r *RAM) Extensive() ([]byte, error) {
	if err := r.load(); err != nil {
		return nil, fmt.Errorf("error in load:%v", err)
	}

	return nil, nil
}

func (r *RAM) load() error {
	if r == nil {
		return errors.New("instance is nil")
	}

	return nil
}
