package repository

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/Jhon-2801/sismos-api/internal/models"
	"gorm.io/gorm"
)

type (
	Repository interface {
		HttpGet(limit, offset int) ([]*models.Events, error)
		HttpCount() (int, error)
		GetFeactures(offsite, limit int) ([]models.Events, error)
		PostFectures(features []*models.Events) error
	}

	repo struct {
		db *gorm.DB
	}
)

func NewRepo(db *gorm.DB) Repository {
	return &repo{db: db}
}

// PostFectures implements Repository.
// Todo: Agregar canal para notificar el error
func (repo *repo) PostFectures(features []*models.Events) error {
	var wg sync.WaitGroup
	for _, v := range features {
		wg.Add(1)
		go func(feature *models.Events) {
			defer wg.Done()
			var count int64

			if err := repo.db.Model(&models.Events{}).Where("event_id = ?", feature.EventID).Count(&count).Error; err != nil {
				return
			}

			// El evento no existe, proceder con la inserción
			if err := repo.db.Create(feature).Error; err != nil {
				return
			}
		}(v)
	}
	wg.Wait()
	return nil
}

// GetFeactures implements Repository.
func (repo *repo) GetFeactures(offsite, limit int) ([]models.Events, error) {
	var c []models.Events
	tx := repo.db.Model(&c)
	// tx = applyFilters(tx, filters)
	tx = tx.Limit(limit).Offset(offsite)
	result := tx.Order("created_at desc").Find(&c)
	if result.Error != nil {
		return nil, result.Error
	}
	return c, nil
}

func (repo *repo) HttpCount() (int, error) {
	var count int64
	tx := repo.db.Model(models.Events{})
	if err := tx.Count(&count).Error; err != nil {
		return 0, nil
	}
	return int(count), nil
}

// httpGet implements Repository.
func (repo *repo) HttpGet(limit, offset int) ([]*models.Events, error) {

	var (
		features []*models.Events
		mutex    sync.Mutex
		wg       sync.WaitGroup
	)
	starttime := time.Now()
	endtime := starttime.AddDate(0, 0, -30)
	starttimeF := starttime.Format("2006-01-02")
	endtimef := endtime.Format("2006-01-02")

	path := fmt.Sprintf("https://earthquake.usgs.gov/fdsnws/event/1/query?format=geojson&starttime=%s&endtime=%s&limit=%d&offset=%d", endtimef, starttimeF, limit, offset)
	resp, err := http.Get(path)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	// Decodificar la respuesta JSON en una estructura GeoJSON
	var geoJSON struct {
		Features []struct {
			ID         string `json:"id"`
			Properties struct {
				Mag     float64           `json:"mag"`
				Place   string            `json:"place"`
				Time    models.CustomTime `json:"time"`
				URL     string            `json:"url"`
				Tsunami int               `json:"tsunami"`
				MagType string            `json:"magType"`
				Title   string            `json:"title"`
			} `json:"properties"`
			Geometry struct {
				Coordinates []float64 `json:"coordinates"`
			} `json:"geometry"`
		} `json:"features"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&geoJSON); err != nil {
		fmt.Println(err)
		return nil, err
	}

	for _, f := range geoJSON.Features {
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
	return features, nil
}
