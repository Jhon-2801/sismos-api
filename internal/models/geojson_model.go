package models

type GeoJSON struct {
	Features []struct {
		ID         string `json:"id"`
		Properties struct {
			Mag     float64    `json:"mag"`
			Place   string     `json:"place"`
			Time    CustomTime `json:"time"`
			URL     string     `json:"url"`
			Tsunami int        `json:"tsunami"`
			MagType string     `json:"magType"`
			Title   string     `json:"title"`
		} `json:"properties"`
		Geometry struct {
			Coordinates []float64 `json:"coordinates"`
		} `json:"geometry"`
	} `json:"features"`
}
