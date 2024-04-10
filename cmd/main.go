package main

import (
	"log"

	"github.com/Jhon-2801/sismos-api/database"
	"github.com/Jhon-2801/sismos-api/internal/api"
	"github.com/Jhon-2801/sismos-api/internal/repository"
	"github.com/Jhon-2801/sismos-api/internal/services"
	"github.com/gin-gonic/gin"
)

func main() {
	db, err := database.ConnectionDB()

	if err != nil {
		log.Fatalf(err.Error())
	}

	sismoRepo := repository.NewRepo(db)
	sismoServ := services.NewService(sismoRepo)
	sismoEnd := api.MakeEndPoints(sismoServ)

	gin.SetMode(gin.ReleaseMode)

	router := gin.Default()
	router.GET("/api/events", gin.HandlerFunc(sismoEnd.GetFeactures))

	router.Run(":8080")
}
