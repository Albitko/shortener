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
	"github.com/Albitko/shortener/internal/repo/memstorage"
	"github.com/Albitko/shortener/internal/repo/postgres"
	"github.com/Albitko/shortener/internal/repo/usermemstorage"
	"github.com/Albitko/shortener/internal/usecase"
	"github.com/Albitko/shortener/internal/workers"
)

type rep interface {
	BatchDeleteShortURLs(context.Context, []entity.ModelURLForDelete) error
}

// Run main application func that runs App
func Run(cfg entity.Config) {
	var db *postgres.DB
	var r rep

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	repository := memstorage.New(cfg.FileStoragePath)
	defer repository.Close()
	userRepository := usermemstorage.New()
	uc := usecase.New(repository, userRepository, db)
	r = repository
	if cfg.DatabaseDSN != "" {
		db = postgres.New(ctx, cfg.DatabaseDSN)
		defer db.Close()
		uc = usecase.New(db, db, db)
		r = db
	}

	queue := workers.Init(ctx, r)
	handler := controller.New(uc, cfg.BaseURL, queue)
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
	router.DELETE("/api/user/urls", handler.DeleteURL)

	err := router.Run(cfg.ServerAddress)
	if err != nil {
		log.Print("Couldn't  start server ", err)
		return
	}
}
