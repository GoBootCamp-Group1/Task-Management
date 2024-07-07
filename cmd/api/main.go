package main

import (
	http_server "github.com/GoBootCamp-Group1/Task-Management/api/http"
	"github.com/GoBootCamp-Group1/Task-Management/config"
	"github.com/GoBootCamp-Group1/Task-Management/internal/core/service"
	"log"
	"os"
	"path/filepath"
)

func main() {
	cfg := readConfig()

	app, err := service.NewAppContainer(cfg)
	if err != nil {
		log.Fatal(err)
	}

	http_server.Run(cfg.Server, app)
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

	cfg, err := config.ReadStandard(configPath)

	if err != nil {
		log.Fatal(err)
	}

	return cfg
}
