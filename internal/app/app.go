package app

import (
	"flag"
	"log"

	"github.com/caarlos0/env/v6"
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

func Run() {
	var cfg entity.Config
	var repository storage

	flag.StringVar(&cfg.ServerAddress, "a", ":8080", "port to listen on")
	flag.StringVar(&cfg.BaseURL, "b", "http://localhost:8080", "http://HOST:PORT")
	flag.StringVar(&cfg.FileStoragePath, "f", "", "File that stores URL -> ID")
	flag.Parse()

	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}

	repository = repo.NewRepository(cfg.FileStoragePath)
	defer repository.Close()
	uc := usecase.NewURLConverter(repository)
	handler := controller.NewURLHandler(uc, cfg.BaseURL)

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(gzip.Gzip(gzip.DefaultCompression))

	router.POST("/", handler.URLToID)
	router.POST("/api/shorten", handler.URLToIDInJSON)
	router.GET("/:id", handler.GetID)

	err = router.Run(cfg.ServerAddress)
	if err != nil {
		log.Fatal("Couldn't  start server ", err)
		return
	}
}
