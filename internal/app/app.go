package app

import (
	"log"

	"github.com/caarlos0/env/v6"
	"github.com/gin-gonic/gin"

	"github.com/Albitko/shortener/internal/controller"
	"github.com/Albitko/shortener/internal/entity"
	"github.com/Albitko/shortener/internal/usecase"
	"github.com/Albitko/shortener/internal/usecase/repo"
)

func Run() {
	serverAddress := ":8080"
	var cfg entity.Config
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}

	repository := repo.NewRepository()
	uc := usecase.NewURLConverter(repository)
	handler := controller.NewURLHandler(uc, cfg.BaseURL)

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	router.POST("/", handler.URLToID)
	router.POST("/api/shorten", handler.URLToIDInJSON)
	router.GET("/:id", handler.GetID)

	if cfg.ServerAddress != "" {
		serverAddress = cfg.ServerAddress
	}

	err = router.Run(serverAddress)
	if err != nil {
		log.Fatal("Couldn't  start server ", err)
		return
	}
}
