package controller

import (
	"github.com/Albitko/shortener/internal/entity"
	"github.com/Albitko/shortener/internal/usecase"
	"io"
	"log"
	"net/http"
)

type URLHandler interface {
	ServeHTTP(http.ResponseWriter, *http.Request)
}

type urlHandler struct {
	uc usecase.URLConverter
}

func NewURLHandler(u usecase.URLConverter) URLHandler {
	return &urlHandler{
		uc: u,
	}
}

func (h *urlHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				log.Print("Could`t read request body")
			}
		}(r.Body)

		originalURL, err := io.ReadAll(r.Body)

		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		w.WriteHeader(http.StatusCreated)
		shortURL := h.uc.URLToID(entity.OriginalURL(originalURL[:]))

		log.Print("POST URL:", string(originalURL[:]), " id: ", shortURL, "\n")

		response := "http://localhost:8080/" + shortURL
		_, err = w.Write([]byte(response))
		if err != nil {
			log.Print("Could`t write response")
			return
		}

	case http.MethodGet:
		shortPath := r.URL.EscapedPath()

		//val, ok := h.uc.IDToURL(entity.URLID(shortPath[1:]))
		//fmt.Print(val, ok, shortPath[1:])
		if val, ok := h.uc.IDToURL(entity.URLID(shortPath[1:])); ok {
			w.Header().Set("Location", string(val))
			w.WriteHeader(http.StatusTemporaryRedirect)
			log.Print("GET id:", shortPath[1:], " URL: ", val, "\n")
		}

	default:
		w.WriteHeader(400)
	}
}
