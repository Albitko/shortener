package app

import (
	"github.com/Albitko/shortener/internal/controller"
	"github.com/Albitko/shortener/internal/usecase"
	"github.com/Albitko/shortener/internal/usecase/repo"
	"log"
	"net/http"
)

func Run() {
	repository := repo.NewRepository()
	uc := usecase.NewUrlConverter(repository)
	handler := controller.NewUrlHandler(uc)

	http.Handle("/", handler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
