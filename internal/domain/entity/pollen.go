package entity

import "time"

type Pollen []struct {
	Datetime      time.Time
	PollenCount   int
	WindDirection string
	WindSpeed     int
	Temperature   int
	Rainfall      int
}

type PollenRepository interface {
	FetchPollen(areaName string, observatoryID int64, from, to time.Time) error
}
