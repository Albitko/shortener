package app

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

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
	"github.com/Albitko/shortener/internal/utils"
	"github.com/Albitko/shortener/internal/workers"
)

type rep interface {
	BatchDeleteShortURLs(context.Context, []entity.ModelURLForDelete) error
}

// Run main application func that runs App
func Run(cfg entity.Config) {
	var db *postgres.DB
	var r rep
	var err error

	ctx, cancel := context.WithCancel(context.Background())
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	repository := memstorage.New(cfg.FileStoragePath)

	userRepository := usermemstorage.New()
	uc := usecase.New(repository, userRepository, db)
	r = repository
	if cfg.DatabaseDSN != "" {
		db = postgres.New(ctx, cfg.DatabaseDSN)
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

	srv := &http.Server{
		Addr:    cfg.ServerAddress,
		Handler: router,
	}
	go func() {
		<-signalChan

		log.Println("Shutting down...")
		cancel()
		if err = srv.Shutdown(ctx); err != nil {
			log.Println("error shutting down the server: ", err)
		}

		err = repository.Close()
		if err != nil {
			log.Println("error closing database: ", err)
			return
		}
		if cfg.DatabaseDSN != "" {
			db.Close()
		}
	}()

	if cfg.EnableHTTPS {
		certPath, keyPath, errCertCreate := utils.CreateCertAndKeyFiles()
		if errCertCreate != nil {
			log.Print("error crete crt anf key files for HTTPS ", errCertCreate)
		}
		if err = srv.ListenAndServeTLS(certPath, keyPath); err != nil {
			log.Println("Couldn't  start server: ", err)
		}
	} else {
		if err = srv.ListenAndServe(); err != nil {
			log.Println("Couldn't  start server: ", err)
		}
	}
}
