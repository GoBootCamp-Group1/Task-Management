package app

import (
	"github.com/GoBootCamp-Group1/Task-Management/config"
	storage2 "github.com/GoBootCamp-Group1/Task-Management/internal/adapters/storage"
	"github.com/GoBootCamp-Group1/Task-Management/internal/core/ops"
	"github.com/GoBootCamp-Group1/Task-Management/internal/core/service"
	"log"

	"gorm.io/gorm"
)

type AppContainer struct {
	cfg         config.Config
	dbConn      *gorm.DB
	userService *service.UserService
	authService *service.AuthService
}

func NewAppContainer(cfg config.Config) (*AppContainer, error) {
	app := &AppContainer{
		cfg: cfg,
	}

	app.mustInitDB()
	storage2.Migrate(app.dbConn)

	app.setUserService()
	app.setAuthService()

	return app, nil
}

func (a *AppContainer) RawRBConnection() *gorm.DB {
	return a.dbConn
}

func (a *AppContainer) UserService() *service.UserService {
	return a.userService
}

func (a *AppContainer) AuthService() *service.AuthService {
	return a.authService
}

func (a *AppContainer) setUserService() {
	if a.userService != nil {
		return
	}
	a.userService = service.NewUserService(ops.NewOps(storage2.NewUserRepo(a.dbConn)))
}

func (a *AppContainer) mustInitDB() {
	if a.dbConn != nil {
		return
	}

	db, err := storage2.NewPostgresGormConnection(a.cfg.DB)
	if err != nil {
		log.Fatal(err)
	}

	a.dbConn = db
}

func (a *AppContainer) setAuthService() {
	if a.authService != nil {
		return
	}

	a.authService = service.NewAuthService(ops.NewOps(storage2.NewUserRepo(a.dbConn)), []byte(a.cfg.Server.TokenSecret),
		a.cfg.Server.TokenExpMinutes,
		a.cfg.Server.RefreshTokenExpMinutes)
}
