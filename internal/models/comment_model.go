package models

import (
	"time"
)

// Comment representa un comentario asociado a un evento sísmico
type Comment struct {
	ID        int       `json:"id"`
	FeatureID int       `json:"feature_id"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"created_at"`
}
