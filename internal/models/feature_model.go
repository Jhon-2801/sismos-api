package models

import "time"

type Events struct {
	ID        int
	EventID   string
	Magnitude float64
	Place     string
	EventTime time.Time
	URL       string
	Tsunami   bool
	MagType   string
	Title     string
	Longitude float64
	Latitude  float64
	CreatedAt time.Time
}
