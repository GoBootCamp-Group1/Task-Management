package main

import (
	http_server "github.com/GoBootCamp-Group1/Task-Management/api/http"
	"github.com/GoBootCamp-Group1/Task-Management/cmd/api/app"
	"github.com/GoBootCamp-Group1/Task-Management/config"
	"log"
	"os"
	"path/filepath"
)

func main() {
	cfg := readConfig()

	appContainer, err := app.NewAppContainer(cfg)
	if err != nil {
		log.Fatal(err)
	}

	http_server.Run(cfg.Server, appContainer)
}

func readConfig() config.Config {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	configPath := filepath.Join(dir, "config.yaml")

	if len(configPath) == 0 {
		log.Fatal("configuration file not found")
	}

	cfg := config.MustReadStandard(configPath)

	return cfg
}
