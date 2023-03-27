package app

import (
	"context"
	"log"

	"github.com/gin-contrib/gzip"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"

	"github.com/Albitko/shortener/internal/controller"
	"github.com/Albitko/shortener/internal/entity"
	"github.com/Albitko/shortener/internal/repo"
	"github.com/Albitko/shortener/internal/usecase"
)

func Run(cfg entity.Config) {
	var db *repo.DB
	repository := repo.NewRepository(cfg.FileStoragePath)
	defer repository.Close()
	userRepository := repo.NewUserRepo()
	uc := usecase.NewURLConverter(repository, userRepository, db)
	if cfg.DatabaseDSN != "" {
		db = repo.NewPostgres(context.Background(), cfg.DatabaseDSN)
		defer db.Close()
		uc = usecase.NewURLConverter(db, db, db)
	}

	handler := controller.NewURLHandler(uc, cfg.BaseURL)
	store := cookie.NewStore([]byte(cfg.CookiesStorageSecret))

	router := gin.New()
	router.Use(sessions.Sessions("session", store))
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(gzip.Gzip(gzip.DefaultCompression, gzip.WithDecompressFn(gzip.DefaultDecompressHandle)))

	router.POST("/", handler.URLToID)
	router.POST("/api/shorten", handler.URLToIDInJSON)
	router.POST("/api/shorten/batch", handler.BatchURLToIDInJSON)
	router.GET("/:id", handler.GetID)
	router.GET("/api/user/urls", handler.GetIDForUser)
	router.GET("/ping", handler.CheckDBConnection)

	err := router.Run(cfg.ServerAddress)
	if err != nil {
		log.Fatal("Couldn't  start server ", err)
		return
	}
}
