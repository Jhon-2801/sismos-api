package models

import "time"

type Events struct {
	ID        int
	EventID   string    `form:"event_id" json:"event_id"`
	Magnitude float64   `form:"magnitude" json:"magnitude"`
	Place     string    `form:"place" json:"place"`
	EventTime time.Time `form:"time" json:"time"  time_format:"2006-01-02T15:04:05Z07:00"`
	URL       string    `form:"external_url" json:"external_url"`
	Tsunami   bool      `form:"tsunami" json:"tsunami"`
	MagType   string    `form:"mag_type" json:"mag_type"`
	Title     string    `form:"title" json:"title"`
	Longitude float64   `form:"longitude" json:"longitude"`
	Latitude  float64   `form:"latitude" json:"latitude"`
	CreatedAt time.Time
}
