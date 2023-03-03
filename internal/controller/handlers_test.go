package controller

import (
	"bytes"
	gz "compress/gzip"
	"github.com/Albitko/shortener/internal/config"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Albitko/shortener/internal/usecase"
	"github.com/Albitko/shortener/internal/usecase/repo"
)

func testRequest(t *testing.T, ts *httptest.Server, method, path string, body []byte, needCompress bool) (int, http.Header, string) {
	var req *http.Request
	var err error
	var reader io.Reader

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	if needCompress {
		var buf bytes.Buffer
		zw := gz.NewWriter(&buf)
		_, _ = zw.Write(body)
		_ = zw.Close()
		req, err = http.NewRequest(method, ts.URL+path, bytes.NewBuffer(buf.Bytes()))
		require.NoError(t, err)
		req.Header.Add("Accept-Encoding", "gzip")
		req.Header.Add("Content-Encoding", "gzip")
	} else {
		req, err = http.NewRequest(method, ts.URL+path, bytes.NewBuffer(body))
		req.Header.Set("Accept-Encoding", "identity")
		require.NoError(t, err)
	}

	resp, err := client.Do(req)
	require.NoError(t, err)

	defer resp.Body.Close()
	if resp.Header.Get("Content-Encoding") == "gzip" {
		gzReader, _ := gz.NewReader(resp.Body)
		reader = gzReader
		defer gzReader.Close()
	} else {
		reader = resp.Body
	}

	respBody, err := io.ReadAll(reader)
	require.NoError(t, err)

	return resp.StatusCode, resp.Header, string(respBody)
}

func setupRouter() *gin.Engine {
	cfg := config.NewConfig()

	repository := repo.NewRepository("")
	uc := usecase.NewURLConverter(repository)
	handler := NewURLHandler(uc, cfg.BaseURL)
	router := gin.Default()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(gzip.Gzip(gzip.DefaultCompression, gzip.WithDecompressFn(gzip.DefaultDecompressHandle)))

	router.POST("/", handler.URLToID)
	router.POST("/api/shorten", handler.URLToIDInJSON)
	router.GET("/:id", handler.GetID)
	return router
}

func TestRouter(t *testing.T) {

	router := setupRouter()
	ts := httptest.NewServer(router)
	defer ts.Close()

	pStatus, _, body := testRequest(t, ts, "POST", "/", []byte("https://google.com"), false)
	assert.Equal(t, http.StatusCreated, pStatus)
	assert.Equal(t, "http://localhost:8080/cv6VxVduxj", body)

	gStatus, gHeaders, _ := testRequest(t, ts, "GET", "/cv6VxVduxj", nil, false)
	assert.Equal(t, http.StatusTemporaryRedirect, gStatus)
	assert.Equal(t, "https://google.com", gHeaders.Get("Location"))

	badStatus, _, _ := testRequest(t, ts, "POST", "/", []byte("SOME_STRING"), false)
	assert.Equal(t, http.StatusBadRequest, badStatus)

	jStatus, _, body := testRequest(t, ts, "POST", "/api/shorten", []byte(`{"url":"https://yandex.ru"}`), false)
	assert.Equal(t, http.StatusCreated, jStatus)
	assert.Equal(t, `{"result":"http://localhost:8080/4eVSAfM3-P"}`, body)

	cStatus, _, cBody := testRequest(t, ts, "POST", "/", []byte(`https://bing.com`), true)
	assert.Equal(t, http.StatusCreated, cStatus)
	assert.Equal(t, `http://localhost:8080/asnI5ScKGD`, cBody)

}
