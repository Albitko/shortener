package controller

import (
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"

	"github.com/Albitko/shortener/internal/entity"
)

type urlConverter interface {
	URLToID(url entity.OriginalURL) entity.URLID
	IDToURL(entity.URLID) (entity.OriginalURL, bool)
}

type urlHandler struct {
	uc urlConverter
}

func NewURLHandler(u urlConverter) *urlHandler {
	return &urlHandler{
		uc: u,
	}
}

func (h *urlHandler) GetID(c *gin.Context) {
	id := c.Param("id")

	if originalURL, ok := h.uc.IDToURL(entity.URLID(id)); ok {
		c.Header("Location", string(originalURL))
		log.Print("GET id:", id, " URL: ", originalURL, "\n")
		c.Status(http.StatusTemporaryRedirect)
	} else {
		c.Status(http.StatusBadRequest)
	}
}

func (h *urlHandler) URLToID(c *gin.Context) {
	originalURL, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
	}
	_, err = url.ParseRequestURI(string(originalURL))
	if err != nil {
		c.String(http.StatusBadRequest, "Should be URL in the body")
	}
	shortURL := h.uc.URLToID(entity.OriginalURL(originalURL[:]))

	log.Print("POST URL:", string(originalURL[:]), " id: ", shortURL, "\n")

	c.String(http.StatusCreated, "http://localhost:8080/"+string(shortURL))
}
