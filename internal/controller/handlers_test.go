package controller

import (
	"bytes"
	"github.com/Albitko/shortener/internal/usecase"
	"github.com/Albitko/shortener/internal/usecase/repo"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func testRequest(t *testing.T, ts *httptest.Server, method, path string, body io.Reader) (*http.Response, string) {
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

	return resp, string(respBody)
}

func setupRouter() *gin.Engine {
	repository := repo.NewRepository()
	uc := usecase.NewURLConverter(repository)
	handler := NewURLHandler(uc)
	router := gin.New()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	router.POST("/", handler.URLToID)
	router.GET("/:id", handler.GetID)
	return router
}

func TestRouter(t *testing.T) {

	router := setupRouter()
	ts := httptest.NewServer(router)
	defer ts.Close()

	reqPost, body := testRequest(t, ts, "POST", "/", bytes.NewBuffer([]byte("https://google.com")))
	assert.Equal(t, http.StatusCreated, reqPost.StatusCode)
	assert.Equal(t, "http://localhost:8080/cv6VxVduxj", body)

	reqGet, _ := testRequest(t, ts, "GET", "/cv6VxVduxj", nil)
	assert.Equal(t, http.StatusTemporaryRedirect, reqGet.StatusCode)
	assert.Equal(t, "https://google.com", reqGet.Header.Get("Location"))

	reqBad, _ := testRequest(t, ts, "POST", "/", bytes.NewBuffer([]byte("SOME_STRING")))
	assert.Equal(t, http.StatusBadRequest, reqBad.StatusCode)
}
