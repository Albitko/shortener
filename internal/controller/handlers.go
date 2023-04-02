package controller

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"

	"github.com/Albitko/shortener/internal/entity"
	"github.com/Albitko/shortener/internal/repo"
	"github.com/Albitko/shortener/internal/workers"
)

type urlConverter interface {
	URLToID(entity.OriginalURL, string) (entity.URLID, error)
	IDToURL(entity.URLID) (entity.OriginalURL, error)
	UserIDToURLs(userID string) (map[string]string, bool)
	AddUserURL(userID string, shortURL string, originalURL string)
	BatchDeleteURL(userID string, shortURLs []string)
	PingDB() error
}

type urlHandler struct {
	uc      urlConverter
	baseURL string
	q       workers.Queue
}

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
	_, err := url.ParseRequestURI(originalURL)
	if err != nil {
		c.String(http.StatusBadRequest, "Should be URL in the body")
		log.Print("ERROR:", err, "\n")
	}
	return h.uc.URLToID(entity.OriginalURL(originalURL), userID)
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
		return value.(string), nil
	}
	err = session.Save()
	if err != nil {
		c.String(200, err.Error())
	}
	return userID, nil
}

func (h *urlHandler) GetID(c *gin.Context) {
	id := c.Param("id")

	originalURL, err := h.uc.IDToURL(entity.URLID(id))
	if err == nil {
		c.Header("Location", string(originalURL))
		//log.Print("GET id:", id, " URL: ", originalURL, "\n")
		c.Status(http.StatusTemporaryRedirect)
	} else if errors.Is(err, repo.ErrURLDeleted) {
		c.String(http.StatusGone, "")
	} else {
		c.Status(http.StatusBadRequest)
	}
}

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

	//log.Print("POST URL:", string(originalURL), " id: ", shortURL, "\n")

	if errors.Is(urlError, repo.ErrURLAlreadyExists) {
		c.String(http.StatusConflict, h.baseURL+string(shortURL))
		return
	}
	c.String(http.StatusCreated, h.baseURL+string(shortURL))
}

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

	//h.uc.BatchDeleteURL(userID, IDsForDelete)
	h.q.Push(&workers.Task{UserID: userID, IDsForDelete: IDsForDelete})
	c.String(http.StatusAccepted, "")
}

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

	for _, val := range requestJSON {
		shortenURL.CorrelationID = val.CorrelationID
		shortID, _ := processURL(c, h, val.OriginalURL, userID)
		shortenURL.ShortURL = h.baseURL + string(shortID)
		response = append(response, shortenURL)
		log.Print("POST URL:", val.OriginalURL, " id: ", shortenURL.ShortURL, "\n")
	}

	c.JSON(http.StatusCreated, response)
}

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

func (h *urlHandler) GetIDForUser(c *gin.Context) {
	var urls []entity.UserURL
	session := sessions.Default(c)
	if userID := session.Get("user"); userID == nil {
		c.String(http.StatusNoContent, "There is no user in the session")
	} else {
		userURLs, ok := h.uc.UserIDToURLs(userID.(string))
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

func (h *urlHandler) CheckDBConnection(c *gin.Context) {
	err := h.uc.PingDB()
	if err != nil {
		c.String(http.StatusInternalServerError, "")
	} else {
		c.String(http.StatusOK, "")
	}
}
