package app

import (
	"log"

	"github.com/gin-gonic/gin"

	"github.com/Albitko/shortener/internal/controller"
	"github.com/Albitko/shortener/internal/usecase"
	"github.com/Albitko/shortener/internal/usecase/repo"
)

func Run() {
	repository := repo.NewRepository()
	uc := usecase.NewURLConverter(repository)
	handler := controller.NewURLHandler(uc)

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	router.POST("/", handler.URLToID)
	router.GET("/:id", handler.GetID)

	err := router.Run(":8080")
	if err != nil {
		log.Fatal("Couldn't  start server ", err)
		return
	}
}
