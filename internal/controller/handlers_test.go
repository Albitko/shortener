package controller

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Albitko/shortener/internal/usecase"
	"github.com/Albitko/shortener/internal/usecase/repo"
)

func testRequest(t *testing.T, ts *httptest.Server, method, path string, body io.Reader) (int, http.Header, string) {
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	req, err := http.NewRequest(method, ts.URL+path, body)
	require.NoError(t, err)

	resp, err := client.Do(req)
	require.NoError(t, err)

	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	return resp.StatusCode, resp.Header, string(respBody)
}

func setupRouter() *gin.Engine {
	repository := repo.NewRepository()
	uc := usecase.NewURLConverter(repository)
	handler := NewURLHandler(uc)
	router := gin.New()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	router.POST("/", handler.URLToID)
	router.POST("/api/shorten", handler.URLToIDInJSON)
	router.GET("/:id", handler.GetID)
	return router
}

func TestRouter(t *testing.T) {

	router := setupRouter()
	ts := httptest.NewServer(router)
	defer ts.Close()

	pStatus, _, body := testRequest(t, ts, "POST", "/", bytes.NewBuffer([]byte("https://google.com")))
	assert.Equal(t, http.StatusCreated, pStatus)
	assert.Equal(t, "http://localhost:8080/cv6VxVduxj", body)

	gStatus, gHeaders, _ := testRequest(t, ts, "GET", "/cv6VxVduxj", nil)
	assert.Equal(t, http.StatusTemporaryRedirect, gStatus)
	assert.Equal(t, "https://google.com", gHeaders.Get("Location"))

	badStatus, _, _ := testRequest(t, ts, "POST", "/", bytes.NewBuffer([]byte("SOME_STRING")))
	assert.Equal(t, http.StatusBadRequest, badStatus)

	jStatus, _, body := testRequest(t, ts, "POST", "/api/shorten", bytes.NewBuffer([]byte(`{"url":"https://yandex.ru"}`)))
	assert.Equal(t, http.StatusCreated, jStatus)
	assert.Equal(t, `{"result":"http://localhost:8080/4eVSAfM3-P"}`, body)
}
