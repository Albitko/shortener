package controller

import (
	"encoding/json"
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
	uc      urlConverter
	baseURL string
}

func NewURLHandler(u urlConverter, envBaseURL string) *urlHandler {
	baseURL := "http://localhost:8080/"
	if envBaseURL != "" {
		baseURL = envBaseURL + "/"
	}
	return &urlHandler{
		uc:      u,
		baseURL: baseURL,
	}
}

func processURL(c *gin.Context, h *urlHandler, originalURL string) entity.URLID {
	_, err := url.ParseRequestURI(originalURL)
	if err != nil {
		c.String(http.StatusBadRequest, "Should be URL in the body")
		log.Print("ERROR:", err, "\n")
	}
	return h.uc.URLToID(entity.OriginalURL(originalURL))
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

	shortURL := processURL(c, h, string(originalURL))

	log.Print("POST URL:", string(originalURL[:]), " id: ", shortURL, "\n")

	c.String(http.StatusCreated, h.baseURL+string(shortURL))
}

func (h *urlHandler) URLToIDInJSON(c *gin.Context) {
	requestJSON := make(map[string]string)
	if err := json.NewDecoder(c.Request.Body).Decode(&requestJSON); err != nil {
		log.Print("ERROR:", err, "\n")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.Header("Content-Type", "application/json")
	shortURL := processURL(c, h, requestJSON["url"])

	log.Print("POST URL:", requestJSON["url"], " id: ", shortURL, "\n")

	c.String(http.StatusCreated, "{\"result\":\""+h.baseURL+string(shortURL)+"\"}")
}
