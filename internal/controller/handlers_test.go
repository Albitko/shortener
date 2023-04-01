package controller

import (
	"bytes"
	gz "compress/gzip"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"runtime"
	"sync"
	"testing"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"golang.org/x/sync/errgroup"

	"github.com/Albitko/shortener/internal/config"
	"github.com/Albitko/shortener/internal/repo"
	"github.com/Albitko/shortener/internal/workers"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Albitko/shortener/internal/usecase"
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
	//queue := workers.NewQueue()
	//wrkrs := make([]*workers.Worker, 0, runtime.NumCPU())
	var db *repo.DB
	repository := repo.NewRepository(cfg.FileStoragePath)
	defer repository.Close()
	userRepository := repo.NewUserRepo()
	uc := usecase.NewURLConverter(repository, userRepository, db)
	// Init Workers
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	g, _ := errgroup.WithContext(ctx)
	recordCh := make(chan workers.Task, 50)
	doneCh := make(chan struct{})
	mu := &sync.Mutex{}
	inWorker := workers.NewInputWorker(recordCh, doneCh, ctx, mu)
	for i := 1; i <= runtime.NumCPU(); i++ {
		outWorker := workers.NewOutputWorker(i, recordCh, doneCh, ctx, db, mu)
		g.Go(outWorker.Do)
	}
	g.Go(inWorker.Loop)
	if cfg.DatabaseDSN != "" {
		db = repo.NewPostgres(context.Background(), cfg.DatabaseDSN)
		defer db.Close()
		uc = usecase.NewURLConverter(db, db, db)
		//for i := 0; i < runtime.NumCPU(); i++ {
		//	wrkrs = append(wrkrs, workers.NewWorker(i, queue, workers.NewResizer(db)))
		//}
		//
		//for _, w := range wrkrs {
		//	go w.Loop()
		//}
	}

	handler := NewURLHandler(uc, cfg.BaseURL, inWorker)
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
}
