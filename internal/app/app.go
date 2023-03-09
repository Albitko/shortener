package app

import (
	"github.com/gin-contrib/gzip"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"log"

	"github.com/Albitko/shortener/internal/controller"
	"github.com/Albitko/shortener/internal/entity"
	"github.com/Albitko/shortener/internal/usecase"
	"github.com/Albitko/shortener/internal/usecase/repo"
)

func Run(cfg entity.Config) {
	repository := repo.NewRepository(cfg.FileStoragePath)
	defer repository.Close()
	userRepository := repo.NewUserRepo()
	uc := usecase.NewURLConverter(repository, userRepository)
	handler := controller.NewURLHandler(uc, cfg.BaseURL)
	store := cookie.NewStore([]byte(cfg.CookiesStorageSecret))

	router := gin.New()
	router.Use(sessions.Sessions("session", store))

	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(gzip.Gzip(gzip.DefaultCompression, gzip.WithDecompressFn(gzip.DefaultDecompressHandle)))

	router.POST("/", handler.URLToID)
	router.POST("/api/shorten", handler.URLToIDInJSON)
	router.GET("/:id", handler.GetID)
	router.GET("/api/user/urls", handler.GetIDForUser)

	err := router.Run(cfg.ServerAddress)
	if err != nil {
		log.Fatal("Couldn't  start server ", err)
		return
	}
}
