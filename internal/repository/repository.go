package repository

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/Jhon-2801/sismos-api/internal/models"
	"gorm.io/gorm"
)

type (
	Repository interface {
		HttpGet(limit, offset int) (*models.GeoJSON, error)
		HttpCount() (int, error)
		GetFeatures(offsite, limit int) ([]models.Events, error)
		PostFeatures(features []*models.Events) error
	}

	repo struct {
		db *gorm.DB
	}
)

func NewRepo(db *gorm.DB) Repository {
	return &repo{db: db}
}

// httpGet implements Repository.
func (repo *repo) HttpCount() (int, error) {
	starttime := time.Now()
	endtime := starttime.AddDate(0, 0, -30)
	starttimeF := starttime.Format("2006-01-02")
	endtimef := endtime.Format("2006-01-02")

	path := fmt.Sprintf("https://earthquake.usgs.gov/fdsnws/event/1/count?starttime=%s&endtime=%s", endtimef, starttimeF)
	resp, err := http.Get(path)
	if err != nil {
		return 0, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	count, err := strconv.Atoi(string(body))
	if err != nil {
		return 0, err
	}

	return count, nil
}

// httpGet implements Repository.
func (repo *repo) HttpGet(limit, offset int) (*models.GeoJSON, error) {

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

	// Crear una instancia de GeoJSON
	geoJSON := &models.GeoJSON{}

	// Decodificar la respuesta JSON en la instancia GeoJSON
	if err := json.NewDecoder(resp.Body).Decode(geoJSON); err != nil {
		fmt.Println(err)
		return nil, err
	}

	return geoJSON, nil
}

// GetFeactures implements Repository.
func (repo *repo) GetFeatures(offsite, limit int) ([]models.Events, error) {
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

// PostFectures implements Repository.
// Todo: Agregar canal para notificar el error
func (repo *repo) PostFeatures(features []*models.Events) error {
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
