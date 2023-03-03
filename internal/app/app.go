package app

import (
	"log"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"

	"github.com/Albitko/shortener/internal/controller"
	"github.com/Albitko/shortener/internal/entity"
	"github.com/Albitko/shortener/internal/usecase"
	"github.com/Albitko/shortener/internal/usecase/repo"
)

type storage interface {
	AddURL(entity.URLID, entity.OriginalURL)
	GetURLByID(entity.URLID) (entity.OriginalURL, bool)
	Close() error
}

func Run(cfg entity.Config) {
	var repository storage

	repository = repo.NewRepository(cfg.FileStoragePath)
	defer repository.Close()
	uc := usecase.NewURLConverter(repository)
	handler := controller.NewURLHandler(uc, cfg.BaseURL)

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(gzip.Gzip(gzip.DefaultCompression, gzip.WithDecompressFn(gzip.DefaultDecompressHandle)))

	router.POST("/", handler.URLToID)
	router.POST("/api/shorten", handler.URLToIDInJSON)
	router.GET("/:id", handler.GetID)

	err := router.Run(cfg.ServerAddress)
	if err != nil {
		log.Fatal("Couldn't  start server ", err)
		return
	}
}
