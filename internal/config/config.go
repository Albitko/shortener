package config

import (
	"flag"
	"log"

	"github.com/caarlos0/env/v6"

	"github.com/Albitko/shortener/internal/entity"
)

// New create app config. If the application is running with keys, then it sets their values.
// If environment variables are set, their values will overwrite the values set by the keys.
func New() entity.Config {
	var cfg entity.Config

	flag.StringVar(&cfg.ServerAddress, "a", ":8080", "port to listen on")
	flag.StringVar(&cfg.BaseURL, "b", "http://localhost:8080", "http://HOST:PORT")
	flag.StringVar(&cfg.FileStoragePath, "f", "", "File that stores URL -> ID")
	flag.StringVar(&cfg.CookiesStorageSecret, "x", "secret", "secret for cookies storage")
	flag.StringVar(&cfg.DatabaseDSN, "d", "", "connection URL for Postgres")
	flag.BoolVar(&cfg.EnableHTTPS, "s", false, "enable HTTPS serve")
	flag.StringVar(&cfg.Config, "c", "", "path to JSON config")
	flag.StringVar(&cfg.TrustedSubnet, "t", "", "trusted subnet IPs mask")
	flag.Parse()

	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}

	if cfg.Config != "" {
		err = parseJSON(cfg.Config, &cfg)
		if err != nil {
			log.Fatal(err)
		}
	}

	return cfg
}
