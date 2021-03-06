package entity

import "time"

type Pollen struct {
	Datetime      time.Time
	PollenCount   *int64
	WindDirection *string
	WindSpeed     *int64
	Temperature   *float64
	Rainfall      *int64
}

type PollenRepository interface {
	FetchPollen(area Area, observatory Observatory, from, to time.Time) ([]Pollen, error)
}
