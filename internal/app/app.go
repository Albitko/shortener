package app

import (
	"context"
	"log"
	"runtime"
	"sync"

	"github.com/gin-contrib/gzip"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"

	"github.com/Albitko/shortener/internal/controller"
	"github.com/Albitko/shortener/internal/entity"
	"github.com/Albitko/shortener/internal/repo"
	"github.com/Albitko/shortener/internal/usecase"
	"github.com/Albitko/shortener/internal/workers"
)

func Run(cfg entity.Config) {
	var db *repo.DB

	//queue := workers.NewQueue()
	//wrkrs := make([]*workers.Worker, 0, runtime.NumCPU())

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

	if cfg.DatabaseDSN != "" {
		db = repo.NewPostgres(context.Background(), cfg.DatabaseDSN)
		defer db.Close()
		uc = usecase.NewURLConverter(db, db, db)

		for i := 1; i <= runtime.NumCPU(); i++ {
			outWorker := workers.NewOutputWorker(i, recordCh, doneCh, ctx, db, mu)
			g.Go(outWorker.Do)
		}
		g.Go(inWorker.Loop)

		//for i := 0; i < runtime.NumCPU(); i++ {
		//	wrkrs = append(wrkrs, workers.NewWorker(i, queue, workers.NewResizer(db)))
		//}
		//
		//for _, w := range wrkrs {
		//	go w.Loop()
		//}
	}

	handler := controller.NewURLHandler(uc, cfg.BaseURL, inWorker)
	//handler := controller.NewURLHandler(uc, cfg.BaseURL, queue)

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

	err := router.Run(cfg.ServerAddress)
	if err != nil {
		log.Print("Couldn't  start server ", err)
		return
	}
}
