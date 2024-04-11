package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Jhon-2801/courses-meta/meta"
	"github.com/Jhon-2801/sismos-api/internal/models"
	"github.com/Jhon-2801/sismos-api/internal/services"
	"github.com/gin-gonic/gin"
)

type (
	Controller func(c *gin.Context)
	Endpoints  struct {
		GetFeactures  Controller
		PostComment   Controller
		UpdateFeature Controller
	}
	CommentReq struct {
		Body string `form:"body"`
	}
)

func MakeEndPoints(s services.Service) Endpoints {
	return Endpoints{
		GetFeactures:  makeGetFeactures(s),
		UpdateFeature: makeUpdateFeacture(s),
		PostComment:   makePostComment(s),
	}
}

func makeGetFeactures(s services.Service) Controller {
	return func(c *gin.Context) {
		// Parsear los parámetros de consulta
		pageStr := c.Query("page")
		perPageStr := c.Query("per_page")
		magTypes := c.QueryArray("mag_type[]")

		// Convertir los valores de los parámetros de consulta a enteros
		page, _ := strconv.Atoi(pageStr)
		perPage, _ := strconv.Atoi(perPageStr)

		total, err := s.Count()
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"status": 500, "message": err})
			return
		}
		meta, err := meta.New(page, perPage, total, "10")

		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"status": 500, "message": err})
			return
		}

		data, err := s.GetFeactures(meta.PerPage, meta.Page, magTypes)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"status": 500, "message": err})
			return
		}

		c.IndentedJSON(http.StatusOK, gin.H{
			"data": data,
			"pagination": gin.H{
				"curren_page": meta.Page,
				"total":       total,
				"perPage":     meta.PerPage,
			}})

	}
}

func makeUpdateFeacture(s services.Service) Controller {
	return func(c *gin.Context) {
		var req models.Events
		err := c.ShouldBind(&req)
		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"status": 400, "message": err})
			return
		}
		idStr := c.Param("id")

		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"status": 400, "message": err})
			return
		}

		_, err = s.GetFeactureById(id)

		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"status": 400, "message": "feature_id not found"})
			return
		}

		req.ID = id

		requiredFields := []struct {
			Name  string
			Value string
		}{
			{"EventID", req.EventID},
			{"Magnitude", fmt.Sprintf("%f", req.Magnitude)},
			{"Place", req.Place},
			{"URL", req.URL},
			{"MagType", req.MagType},
			{"Title", req.Title},
			{"Longitude", fmt.Sprintf("%f", req.Longitude)},
			{"Latitude", fmt.Sprintf("%f", req.Latitude)},
		}
		for _, field := range requiredFields {
			if field.Value == "" {
				c.IndentedJSON(http.StatusBadRequest, gin.H{"status": 400, "message": fmt.Sprintf("%s is required", field.Name)})
				return
			}
		}

		err = s.UpdateFeactureById(&req)

		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"status": 500, "message": err})
			return
		}
	}
}

func makePostComment(s services.Service) Controller {
	return func(c *gin.Context) {
		var req CommentReq
		c.ShouldBind(&req)

		if len(req.Body) <= 0 {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"status": 400, "message": "comment is required"})
			return
		}

		idStr := c.Param("id")

		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"status": 400, "message": err})
			return
		}

		_, err = s.GetFeactureById(id)

		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"status": 400, "message": "feature_id not found"})
			return
		}

		err = s.PostComment(id, req.Body)

		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"status": 500, "message": err})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"status": 201, "message": "create"})
	}
}
