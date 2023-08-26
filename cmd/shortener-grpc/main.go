package main

import (
	"fmt"

	"github.com/Albitko/shortener/internal/app"
	"github.com/Albitko/shortener/internal/config"
)

var (
	buildVersion string
	buildDate    string
	buildCommit  string
)

func checkBuild(setting, name string) {
	if setting == "" {
		fmt.Println("Build ", name, ": N/A")
		return
	}
	fmt.Println("Build ", name, ": ", setting)
}

func main() {
	checkBuild(buildVersion, "version")
	checkBuild(buildDate, "date")
	checkBuild(buildCommit, "commit")

	cfg := config.New()
	app.RunGRPC(cfg)
}
