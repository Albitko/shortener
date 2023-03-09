package controller

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/Albitko/shortener/internal/entity"
)

type urlConverter interface {
	URLToID(url entity.OriginalURL) entity.URLID
	IDToURL(entity.URLID) (entity.OriginalURL, bool)
	UserIDToURLs(userID string) (map[string]string, bool)
	AddUserURL(userID string, shortURL string, originalURL string)
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

func checkUserSession(c *gin.Context) (string, error) {
	session := sessions.Default(c)
	randID := make([]byte, 8)
	_, err := rand.Read(randID)
	if err != nil {
		log.Print("ERROR:", err, "\n")
	}
	userID := hex.EncodeToString(randID)

	if value := session.Get("user"); value == nil {
		session.Set("user", userID)
	} else {
		return userID, nil
	}
	err = session.Save()
	if err != nil {
		c.String(200, err.Error())
	}
	return userID, nil
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
	userID, err := checkUserSession(c)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
	}

	originalURL, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
	}

	shortURL := processURL(c, h, string(originalURL))

	h.uc.AddUserURL(userID, h.baseURL+string(shortURL), string(originalURL[:]))

	log.Print("POST URL:", string(originalURL[:]), " id: ", shortURL, "\n")

	c.String(http.StatusCreated, h.baseURL+string(shortURL))
}

func (h *urlHandler) URLToIDInJSON(c *gin.Context) {
	userId, err := checkUserSession(c)
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
	shortURL := processURL(c, h, requestJSON["url"])

	h.uc.AddUserURL(userId, h.baseURL+string(shortURL), requestJSON["url"])

	log.Print("POST URL:", requestJSON["url"], " id: ", shortURL, "\n")

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
				userURL.ShortURL = shortURL
				urls = append(urls, userURL)
			}
			response, _ := json.Marshal(urls)
			c.JSON(http.StatusOK, string(response))
		} else {
			c.String(http.StatusNoContent, "")
		}
	}
}
