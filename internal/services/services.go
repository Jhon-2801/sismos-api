package services

import (
	"github.com/Jhon-2801/sismos-api/internal/models"
	"github.com/Jhon-2801/sismos-api/internal/repository"
)

type (
	Service interface {
		GetFeactures(limit, offset int) ([]*models.Feature, error)
		Count() (int, error)
	}

	service struct {
		repo repository.Repository
	}
)

func NewService(repo repository.Repository) Service {
	return &service{repo: repo}
}

// count implements Service.
func (s *service) Count() (int, error) {
	return s.repo.HttpCount()

}

func (s *service) GetFeactures(limit, offset int) ([]*models.Feature, error) {
	var featuresResponse []*models.Feature

	features, err := s.repo.HttpGet(limit, offset)
	if err != nil {
		return nil, err
	}

	err = s.repo.PostFectures(features)
	if err != nil {
		return nil, err
	}
	featuresDB, err := s.repo.GetFeactures(offset, limit)
	if err != nil {
		return nil, err
	}
	for _, f := range featuresDB {
		// Crear las coordenadas
		coordinates := models.Coordinates{
			Longitude: f.Longitude,
			Latitude:  f.Latitude,
		}

		// Crear el modelo Feature
		feature := &models.Feature{
			ID:   f.ID,
			Type: "feature",
			Attributes: models.FeatureAttributes{
				ExternalID:  f.EventID,
				Magnitude:   f.Magnitude,
				Place:       f.Place,
				Time:        f.EventTime,
				Tsunami:     f.Tsunami,
				MagType:     f.MagType,
				Title:       f.Title,
				Coordinates: coordinates,
			},
			Links: struct {
				ExternalURL string `json:"external_url"`
			}{
				ExternalURL: f.URL,
			},
		}
		featuresResponse = append(featuresResponse, feature)
	}

	return featuresResponse, nil
}
