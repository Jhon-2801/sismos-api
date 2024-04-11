package models

import (
	"encoding/json"
	"time"
)

// Feature representa un evento sísmico
// Coordinates representa las coordenadas geográficas
type Coordinates struct {
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
}

// FeatureAttributes representa los atributos de un evento sísmico
type FeatureAttributes struct {
	ExternalID  string      `json:"event_id"`
	Magnitude   float64     `json:"magnitude"`
	Place       string      `json:"place"`
	Time        time.Time   `json:"time"`
	Tsunami     bool        `json:"tsunami"`
	MagType     string      `json:"mag_type"`
	Title       string      `json:"title"`
	Coordinates Coordinates `json:"coordinates"`
}

// Feature representa un evento sísmico
type Feature struct {
	ID         int               `json:"id"`
	Type       string            `json:"type"`
	Attributes FeatureAttributes `json:"attributes"`
	Links      struct {
		ExternalURL string `json:"external_url"`
	} `json:"links"`
}

type CustomTime struct {
	time.Time
}

func (ct *CustomTime) UnmarshalJSON(b []byte) error {
	var timestamp int64
	err := json.Unmarshal(b, &timestamp)
	if err != nil {
		return err
	}
	ct.Time = time.Unix(0, timestamp*int64(time.Millisecond))
	return nil
}
