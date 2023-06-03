package main

import (
	"github.com/Albitko/shortener/internal/app"
	"github.com/Albitko/shortener/internal/config"
)

func main() {
	cfg := config.New()
	app.Run(cfg)
}
