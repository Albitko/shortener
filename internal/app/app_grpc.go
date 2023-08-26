package app

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"

	gh "github.com/Albitko/shortener/internal/controller/grpc"
	"github.com/Albitko/shortener/internal/entity"
	"github.com/Albitko/shortener/internal/repo/memstorage"
	"github.com/Albitko/shortener/internal/repo/postgres"
	"github.com/Albitko/shortener/internal/repo/usermemstorage"
	"github.com/Albitko/shortener/internal/usecase"
	"github.com/Albitko/shortener/internal/workers"
	pb "github.com/Albitko/shortener/proto"
)

// RunGRPC grpc application func that runs App
func RunGRPC(cfg entity.Config) {
	var db *postgres.DB
	var r rep
	var err error

	ctx, cancel := context.WithCancel(context.Background())
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	repository := memstorage.New(cfg.FileStoragePath)

	userRepository := usermemstorage.New()
	uc := usecase.New(repository, userRepository, db)
	r = repository
	if cfg.DatabaseDSN != "" {
		db = postgres.New(ctx, cfg.DatabaseDSN)
		uc = usecase.New(db, db, db)
		r = db
	}

	queue := workers.Init(ctx, r)
	handler := gh.New(uc, cfg, queue)

	listen, err := net.Listen("tcp", cfg.ServerAddress)
	if err != nil {
		log.Fatal(err)
	}
	s := grpc.NewServer()
	pb.RegisterShortenerServer(s, handler)

	log.Println("gRPC server start")
	if err = s.Serve(listen); err != nil {
		log.Fatal(err)
	}

	go func() {
		<-signalChan

		log.Println("Shutting down...")
		cancel()
		s.Stop()

		err = repository.Close()
		if err != nil {
			log.Println("error closing database: ", err)
			return
		}
		if cfg.DatabaseDSN != "" {
			db.Close()
		}
	}()
}
