package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"

	http_server "github.com/GoBootCamp-Group1/Task-Management/api/http"
	"github.com/GoBootCamp-Group1/Task-Management/cmd/api/app"
	"github.com/GoBootCamp-Group1/Task-Management/config"
)

var configPath = flag.String("config", "", "Configuration Path")

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
	flag.Parse()

	if cfgPathEnv := os.Getenv("APP_CONFIG_PATH"); len(cfgPathEnv) > 0 {
		*configPath = cfgPathEnv
	}

	if len(*configPath) == 0 {
		dir, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}

		*configPath = filepath.Join(dir, "config.yaml")
	}

	cfg := config.MustReadStandard(*configPath)

	return cfg
}
