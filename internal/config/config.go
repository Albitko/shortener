package config

import (
	"flag"
	"log"

	"github.com/caarlos0/env/v6"

	"github.com/Albitko/shortener/internal/entity"
)

func NewConfig() entity.Config {
	var cfg entity.Config

	flag.StringVar(&cfg.ServerAddress, "a", ":8080", "port to listen on")
	flag.StringVar(&cfg.BaseURL, "b", "http://localhost:8080", "http://HOST:PORT")
	flag.StringVar(&cfg.FileStoragePath, "f", "", "File that stores URL -> ID")
	flag.StringVar(&cfg.CookiesStorageSecret, "c", "secret", "secret for cookies storage")
	flag.StringVar(&cfg.DatabaseDSN, "d", "postgresql://localhost:5432", "connection URL for Postgres")
	flag.Parse()

	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}

	return cfg
}
