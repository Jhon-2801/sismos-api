package repository

import (
	"crypto/tls"
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
		GetFeatures(offsite, limit int, filters []string) ([]models.Events, error)
		GetFeatureById(id int) (models.Events, error)
		UpdateFeature(feature *models.Events) error
		PostFeatures(features []*models.Events) error
		PostComment(comment *models.Comment) error
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
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	path := fmt.Sprintf("https://earthquake.usgs.gov/fdsnws/event/1/count?starttime=%s&endtime=%s", endtimef, starttimeF)
	resp, err := client.Get(path)
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
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	path := fmt.Sprintf("https://earthquake.usgs.gov/fdsnws/event/1/query?format=geojson&starttime=%s&endtime=%s&limit=%d&offset=%d", endtimef, starttimeF, limit, offset)
	resp, err := client.Get(path)
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
func (repo *repo) GetFeatures(offsite, limit int, filters []string) ([]models.Events, error) {
	var c []models.Events
	tx := repo.db.Model(&c)
	tx = applyFilters(tx, filters)
	tx = tx.Limit(limit).Offset(offsite)
	result := tx.Order("created_at desc").Find(&c)
	if result.Error != nil {
		return nil, result.Error
	}
	return c, nil
}

// GetFeatureById implements Repository.
func (repo *repo) GetFeatureById(id int) (models.Events, error) {
	feacture := models.Events{}
	err := repo.db.Where("id = ?", id).First(&feacture).Error
	return feacture, err
}

// UpdateFeature implements Repository.
func (repo *repo) UpdateFeature(feature *models.Events) error {
	return repo.db.Save(feature).Error
}

// PostFectures implements Repository.
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

// PostComment implements Repository.
func (repo *repo) PostComment(comment *models.Comment) error {
	if err := repo.db.Create(comment).Error; err != nil {
		return err
	}
	return nil
}

func applyFilters(tx *gorm.DB, filters []string) *gorm.DB {
	if len(filters) == 0 {
		return tx
	}
	tx.Where("mag_type IN (?)", filters)
	return tx

}
