package http

import (
	"bytes"
	gz "compress/gzip"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-contrib/gzip"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Albitko/shortener/internal/config"
	"github.com/Albitko/shortener/internal/entity"
	"github.com/Albitko/shortener/internal/repo/memstorage"
	"github.com/Albitko/shortener/internal/repo/postgres"
	"github.com/Albitko/shortener/internal/repo/usermemstorage"
	"github.com/Albitko/shortener/internal/usecase"
	"github.com/Albitko/shortener/internal/workers"
)

type rep interface {
	BatchDeleteShortURLs(context.Context, []entity.ModelURLForDelete) error
}

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
		req.Header.Add("X-Real-IP", "127.0.0.1")
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
	cfg := config.New()
	cfg.TrustedSubnet = "127.0.0.0/24"
	var db *postgres.DB
	var r rep
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	repository := memstorage.New(cfg.FileStoragePath)
	defer repository.Close()
	r = repository
	userRepository := usermemstorage.New()
	uc := usecase.New(repository, userRepository, db)
	if cfg.DatabaseDSN != "" {
		db = postgres.New(ctx, cfg.DatabaseDSN)
		defer db.Close()
		uc = usecase.New(db, db, db)
		r = db
	}

	queue := workers.Init(ctx, r)
	handler := New(uc, cfg, queue)
	store := cookie.NewStore([]byte(cfg.CookiesStorageSecret))

	router := gin.New()
	router.Use(sessions.Sessions("session", store))
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(gzip.Gzip(gzip.DefaultCompression, gzip.WithDecompressFn(gzip.DefaultDecompressHandle)))

	router.POST("/", handler.URLToID)
	router.POST("/api/shorten", handler.URLToIDInJSON)
	router.POST("/api/shorten/batch", handler.BatchURLToIDInJSON)
	router.GET("/:id", handler.GetID)
	router.GET("/api/user/urls", handler.GetIDForUser)
	router.GET("/api/internal/stats", handler.Stats)
	router.GET("/ping", handler.CheckDBConnection)
	router.DELETE("/api/user/urls", handler.DeleteURL)
	return router
}

func TestRouter(t *testing.T) {
	router := setupRouter()
	ts := httptest.NewServer(router)
	defer ts.Close()

	pStatus, _, body := testRequest(t, ts, "POST", "/", []byte("https://google.com"), false)
	assert.Equal(t, http.StatusCreated, pStatus)
	assert.Equal(t, "http://localhost:8080/BQRvJsg-jI", body)

	gStatus, gHeaders, _ := testRequest(t, ts, "GET", "/BQRvJsg-jI", nil, false)
	assert.Equal(t, http.StatusTemporaryRedirect, gStatus)
	assert.Equal(t, "https://google.com", gHeaders.Get("Location"))

	badStatus, _, _ := testRequest(t, ts, "POST", "/", []byte("SOME_STRING"), false)
	assert.Equal(t, http.StatusBadRequest, badStatus)

	jStatus, _, body := testRequest(t, ts, "POST", "/api/shorten", []byte(`{"url":"https://yandex.ru"}`), false)
	assert.Equal(t, http.StatusCreated, jStatus)
	assert.Equal(t, `{"result":"http://localhost:8080/FgAJzmBKgR"}`, body)

	cStatus, _, cBody := testRequest(t, ts, "POST", "/", []byte(`https://bing.com`), true)
	assert.Equal(t, http.StatusCreated, cStatus)
	assert.Equal(t, `http://localhost:8080/DVgElL_ZX_`, cBody)

	bStatus, _, body := testRequest(t, ts, "POST", "/api/shorten/batch", []byte(`[{"correlation_id": "qwerty123", "original_url": "https://news.com"}, {"correlation_id": "qwerty123", "original_url": "https://mail.com"}]`), false)
	assert.Equal(t, http.StatusCreated, bStatus)
	assert.Equal(t, `[{"correlation_id":"qwerty123","short_url":"http://localhost:8080/3aAJHI89Bk"},{"correlation_id":"qwerty123","short_url":"http://localhost:8080/ojVFbv-Meo"}]`, body)

	pingStatus, _, _ := testRequest(t, ts, "GET", "/ping", nil, false)
	assert.Equal(t, http.StatusInternalServerError, pingStatus)

	// BatchDeleteShortURLs DELETE
	originalURLs := []string{"https://1.com", "https://2.com", "https://3.com"}
	shortenURLs := []string{"L4hn_ks5jl", "Ibh-weQcwK", "2J_SGlAgcs"}

	for i, url := range originalURLs {
		s, _, b := testRequest(t, ts, "POST", "/", []byte(url), false)
		assert.Equal(t, http.StatusCreated, s)
		assert.Equal(t, "http://localhost:8080/"+shortenURLs[i], b)
	}

	for i, url := range shortenURLs {
		s, h, _ := testRequest(t, ts, "GET", "/"+url, nil, false)
		assert.Equal(t, http.StatusTemporaryRedirect, s)
		assert.Equal(t, originalURLs[i], h.Get("Location"))
	}

	s, _, _ := testRequest(t, ts, "DELETE", "/api/user/urls", []byte(`["L4hn_ks5jl", "Ibh-weQcwK", "2J_SGlAgcs"]`), false)
	assert.Equal(t, http.StatusAccepted, s)

	for i := range shortenURLs {
		s, _, _ := testRequest(t, ts, "GET", "/"+shortenURLs[i], nil, false)
		assert.Equal(t, http.StatusGone, s)
	}
	st, _, _ := testRequest(t, ts, "GET", "/api/internal/stats", nil, false)
	assert.Equal(t, http.StatusOK, st)
}
