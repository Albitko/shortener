// Package controller is used to process user requests.
// Each public method of `urlHandler` is associated with 1 API endpoint.
// It prepares the data for forwarding to the use case layer and returns the HTTP err codes
// depending on what the next layer returned.
package controller

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"

	"github.com/Albitko/shortener/internal/entity"
	"github.com/Albitko/shortener/internal/repo"
	"github.com/Albitko/shortener/internal/workers"
)

type urlConverter interface {
	URLToID(context.Context, entity.OriginalURL, string) (entity.URLID, error)
	IDToURL(context.Context, entity.URLID) (entity.OriginalURL, error)
	UserIDToURLs(c context.Context, userID string) (map[string]string, bool)
	PingDB() error
}

type urlHandler struct {
	uc      urlConverter
	baseURL string
	q       workers.Queue
}

// NewURLHandler create instance of `urlHandler` struct
func NewURLHandler(u urlConverter, envBaseURL string, queue *workers.Queue) *urlHandler {
	baseURL := "http://localhost:8080/"
	if envBaseURL != "" {
		baseURL = envBaseURL + "/"
	}
	return &urlHandler{
		uc:      u,
		baseURL: baseURL,
		q:       *queue,
	}
}

func processURL(c *gin.Context, h *urlHandler, originalURL, userID string) (entity.URLID, error) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()
	_, err := url.ParseRequestURI(originalURL)
	if err != nil {
		c.String(http.StatusBadRequest, "Should be URL in the body")
		log.Print("ERROR:", err, "\n")
	}
	return h.uc.URLToID(ctx, entity.OriginalURL(originalURL), userID)
}

func checkUserSession(c *gin.Context) (string, error) {
	session := sessions.Default(c)
	randID := make([]byte, 8)
	_, err := rand.Read(randID)
	if err != nil {
		log.Print("ERROR:", err, "\n")
		return "", err
	}
	userID := hex.EncodeToString(randID)

	if value := session.Get("user"); value == nil {
		session.Set("user", userID)
	} else {
		return fmt.Sprintf("%v", value), nil
	}
	err = session.Save()
	if err != nil {
		c.String(200, err.Error())
	}
	return userID, nil
}

// GetID handler gets full URL from storage by shorten.
func (h *urlHandler) GetID(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()
	id := c.Param("id")

	originalURL, err := h.uc.IDToURL(ctx, entity.URLID(id))
	switch {
	case err == nil:
		c.Header("Location", string(originalURL))
		c.Status(http.StatusTemporaryRedirect)
	case errors.Is(err, repo.ErrURLDeleted):
		c.String(http.StatusGone, "")
	default:
		c.Status(http.StatusBadRequest)
	}
}

// URLToID shorten URL for user.
func (h *urlHandler) URLToID(c *gin.Context) {
	userID, err := checkUserSession(c)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
	}

	originalURL, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
	}

	shortURL, urlError := processURL(c, h, string(originalURL), userID)

	if errors.Is(urlError, repo.ErrURLAlreadyExists) {
		c.String(http.StatusConflict, h.baseURL+string(shortURL))
		return
	}
	c.String(http.StatusCreated, h.baseURL+string(shortURL))
}

// DeleteURL deletes full URL for requested user.
func (h *urlHandler) DeleteURL(c *gin.Context) {
	userID, err := checkUserSession(c)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
	}

	var IDsForDelete []string
	if err := json.NewDecoder(c.Request.Body).Decode(&IDsForDelete); err != nil {
		log.Print("ERROR decoding IDs for deletion:", err, "\n")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	h.q.Push(&workers.Task{UserID: userID, IDsForDelete: IDsForDelete})
	c.String(http.StatusAccepted, "")
}

// BatchURLToIDInJSON receives a request from a user with plenty of URLs
// that need to be shortened in json format and shortens them.
func (h *urlHandler) BatchURLToIDInJSON(c *gin.Context) {
	var requestJSON []entity.ModelURLBatchRequest
	var shortenURL entity.ModelURLBatchResponse

	userID, err := checkUserSession(c)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
	}

	if err := json.NewDecoder(c.Request.Body).Decode(&requestJSON); err != nil {
		log.Print("ERROR:", err, "\n")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	response := make([]entity.ModelURLBatchResponse, 0, len(requestJSON))

	for i := range requestJSON {
		shortenURL.CorrelationID = requestJSON[i].CorrelationID
		shortID, _ := processURL(c, h, requestJSON[i].OriginalURL, userID)
		shortenURL.ShortURL = h.baseURL + string(shortID)
		response = append(response, shortenURL)
		log.Print("POST URL:", requestJSON[i].OriginalURL, " id: ", shortenURL.ShortURL, "\n")
	}

	c.JSON(http.StatusCreated, response)
}

// URLToIDInJSON receives a request from a user in json format with one URL
// that need to be shortened  and shortens them.
func (h *urlHandler) URLToIDInJSON(c *gin.Context) {
	userID, err := checkUserSession(c)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
	}
	requestJSON := make(map[string]string)
	if err := json.NewDecoder(c.Request.Body).Decode(&requestJSON); err != nil {
		log.Print("ERROR:", err, "\n")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.Header("Content-Type", "application/json")
	shortURL, urlError := processURL(c, h, requestJSON["url"], userID)

	log.Print("POST URL:", requestJSON["url"], " id: ", shortURL, "\n")

	if errors.Is(urlError, repo.ErrURLAlreadyExists) {
		c.String(http.StatusConflict, "{\"result\":\""+h.baseURL+string(shortURL)+"\"}")
		return
	}

	c.String(http.StatusCreated, "{\"result\":\""+h.baseURL+string(shortURL)+"\"}")
}

// GetIDForUser return all shorten URLs for requested user.
func (h *urlHandler) GetIDForUser(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()
	var urls []entity.UserURL
	session := sessions.Default(c)
	if userID := session.Get("user"); userID == nil {
		c.String(http.StatusNoContent, "There is no user in the session")
	} else {
		userURLs, ok := h.uc.UserIDToURLs(ctx, userID.(string))
		if ok {
			for shortURL, originalURL := range userURLs {
				var userURL entity.UserURL
				userURL.OriginalURL = originalURL
				userURL.ShortURL = h.baseURL + shortURL
				urls = append(urls, userURL)
			}
			c.JSON(http.StatusOK, urls)
		} else {
			c.String(http.StatusNoContent, "")
		}
	}
}

// CheckDBConnection ping DB.
func (h *urlHandler) CheckDBConnection(c *gin.Context) {
	err := h.uc.PingDB()
	if err != nil {
		c.String(http.StatusInternalServerError, "")
	} else {
		c.String(http.StatusOK, "")
	}
}
