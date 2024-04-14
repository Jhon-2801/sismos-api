package services

import (
	"sync"
	"time"

	"github.com/Jhon-2801/sismos-api/internal/models"
	"github.com/Jhon-2801/sismos-api/internal/repository"
)

type (
	Service interface {
		GetFeactures(limit, offset int, filters []string) ([]*models.Feature, error)
		GetFeactureById(id int) (models.Events, *models.Feature, error)
		UpdateFeactureById(feature *models.Events) (models.Feature, error)
		PostComment(id int, comment string) error
		GetComment(id int) ([]models.Comment, error)
		Count() (int, error)
	}

	service struct {
		repo repository.Repository
	}
)

func NewService(repo repository.Repository) Service {
	return &service{repo: repo}
}

func (s *service) GetFeactures(limit, offset int, filters []string) ([]*models.Feature, error) {

	geoJson, err := s.repo.HttpGet(limit, offset)
	if err != nil {
		return nil, err
	}

	//Persistir en la base de datos
	features := parcerGeoJsonToEvents(geoJson)
	err = s.repo.PostFeatures(features)
	if err != nil {
		return nil, err
	}

	//Obtener de la base de datos
	featuresModels, err := s.repo.GetFeatures(offset, limit, filters)
	if err != nil {
		return nil, err
	}

	//Pasear para mandar model de respuesta
	featuresResponse := parcerModelArrayToResponse(featuresModels)

	return featuresResponse, nil
}

// GetFeactureById implements Service.
func (s *service) GetFeactureById(id int) (models.Events, *models.Feature, error) {
	resp, err := s.repo.GetFeatureById(id)
	if err != nil {
		return resp, nil, err
	}
	featuresResponse := parcerModelToResponse(&resp)
	return resp, featuresResponse, nil
}

// UpdateFeactureById implements Service.
func (s *service) UpdateFeactureById(feature *models.Events) (models.Feature, error) {
	resp, _, _ := s.GetFeactureById(feature.ID)
	feature.CreatedAt = resp.CreatedAt
	err := s.repo.UpdateFeature(feature)
	featuresResponse := parcerModelToResponse(feature)
	return *featuresResponse, err
}

// PostComment implements Service.
func (s *service) PostComment(id int, comment string) error {
	commentModel := models.Comment{
		FeatureID: id,
		Body:      comment,
	}
	return s.repo.PostComment(&commentModel)
}

// GetComment implements Service.
func (s *service) GetComment(idFeature int) ([]models.Comment, error) {
	return s.repo.GetComment(idFeature)
}

// count implements Service.
func (s *service) Count() (int, error) {
	return s.repo.HttpCount()

}

func parcerGeoJsonToEvents(geoJson *models.GeoJSON) []*models.Events {
	var (
		features []*models.Events
		mutex    sync.Mutex
		wg       sync.WaitGroup
	)
	for _, f := range geoJson.Features {
		wg.Add(1)
		go func() {
			defer wg.Done()
			feature := &models.Events{
				EventID:   f.ID,
				Magnitude: f.Properties.Mag,
				Place:     f.Properties.Place,
				EventTime: f.Properties.Time.Time,
				URL:       f.Properties.URL,
				Tsunami:   f.Properties.Tsunami != 0,
				MagType:   f.Properties.MagType,
				Title:     f.Properties.Title,
				Longitude: f.Geometry.Coordinates[0],
				Latitude:  f.Geometry.Coordinates[1],
				CreatedAt: time.Now(),
			}
			mutex.Lock()
			features = append(features, feature)
			mutex.Unlock()
		}()
	}
	wg.Wait()
	return features
}

func parcerModelArrayToResponse(model []models.Events) []*models.Feature {
	var (
		featuresResponse []*models.Feature
		mutex            sync.Mutex
		wg               sync.WaitGroup
	)
	for _, f := range model {
		wg.Add(1)
		go func() {
			defer wg.Done()
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
			mutex.Lock()
			featuresResponse = append(featuresResponse, feature)
			mutex.Unlock()
		}()
	}
	wg.Wait()
	return featuresResponse
}

func parcerModelToResponse(model *models.Events) *models.Feature {
	// Crear las coordenadas
	coordinates := models.Coordinates{
		Longitude: model.Longitude,
		Latitude:  model.Latitude,
	}

	// Crear el modelo Feature
	feature := &models.Feature{
		ID:   model.ID,
		Type: "feature",
		Attributes: models.FeatureAttributes{
			ExternalID:  model.EventID,
			Magnitude:   model.Magnitude,
			Place:       model.Place,
			Time:        model.EventTime,
			Tsunami:     model.Tsunami,
			MagType:     model.MagType,
			Title:       model.Title,
			Coordinates: coordinates,
		},
		Links: struct {
			ExternalURL string `json:"external_url"`
		}{
			ExternalURL: model.URL,
		},
	}

	return feature
}
