package main

import (
	http_server "github.com/GoBootCamp-Group1/Task-Management/api/http"
	"github.com/GoBootCamp-Group1/Task-Management/cmd/api/app"
	"github.com/GoBootCamp-Group1/Task-Management/config"
	"log"
	"os"
	"path/filepath"
)

//	@title			Task Manager
//	@version		1.0
//	@description	This is a task management application.

//	@contact.name	GoBootcamp-Group1
//	@contact.url	https://github.com/GoBootCamp-Group1

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

// @Security ApiKeyAuth

// @host		localhost:8082
// @BasePath	/api/v1
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
