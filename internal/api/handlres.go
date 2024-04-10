package api

import (
	"net/http"
	"strconv"

	"github.com/Jhon-2801/courses-meta/meta"
	"github.com/Jhon-2801/sismos-api/internal/services"
	"github.com/gin-gonic/gin"
)

type (
	Controller func(c *gin.Context)
	Endpoints  struct {
		GetFeactures Controller
	}
)

func MakeEndPoints(s services.Service) Endpoints {
	return Endpoints{
		GetFeactures: makeGetFeactures(s),
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
