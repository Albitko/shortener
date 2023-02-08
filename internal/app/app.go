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
	uc := usecase.NewURLConverter(repository)
	handler := controller.NewURLHandler(uc)

	http.Handle("/", handler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
