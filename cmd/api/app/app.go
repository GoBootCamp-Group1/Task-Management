package app

import (
	"github.com/GoBootCamp-Group1/Task-Management/config"
	"github.com/GoBootCamp-Group1/Task-Management/internal/adapters/cache"
	"github.com/GoBootCamp-Group1/Task-Management/internal/adapters/notifier"
	"github.com/GoBootCamp-Group1/Task-Management/internal/adapters/storage"
	"github.com/GoBootCamp-Group1/Task-Management/internal/core/services"
	"github.com/GoBootCamp-Group1/Task-Management/pkg/notification"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"log"
)

type Container struct {
	cfg           config.Config
	dbConn        *gorm.DB
	cacheClient   *redis.Client
	notifier      *notification.Notifier
	userService   *services.UserService
	authService   *services.AuthService
	boardService  *services.BoardService
	taskService   *services.TaskService
	columnService *services.ColumnService
}

func NewAppContainer(cfg config.Config) (*Container, error) {
	app := &Container{
		cfg: cfg,
	}

	app.mustInitCache()
	app.mustInitDB()
	storage.Migrate(app.dbConn)

	app.initNotifier()

	app.setUserService()
	app.setAuthService()
	app.setBoardService()
	app.setTaskService()
	app.setColumnService()

	return app, nil
}

func (a *Container) RawRBConnection() *gorm.DB {
	return a.dbConn
}

func (a *Container) UserService() *services.UserService {
	return a.userService
}

func (a *Container) AuthService() *services.AuthService {
	return a.authService
}

func (a *Container) BoardService() *services.BoardService {
	return a.boardService
}

func (a *Container) TaskService() *services.TaskService {
	return a.taskService
}

func (a *Container) ColumnService() *services.ColumnService {
	return a.columnService
}

func (a *Container) setUserService() {
	if a.userService != nil {
		return
	}
	a.userService = services.NewUserService(storage.NewUserRepo(a.dbConn))
}

func (a *Container) mustInitDB() {
	if a.dbConn != nil {
		return
	}

	db, err := storage.NewPostgresGormConnection(a.cfg.DB)
	if err != nil {
		log.Fatal(err)
	}

	a.dbConn = db
}

func (a *Container) mustInitCache() {
	if a.cacheClient != nil {
		return
	}

	redisClient, err := cache.NewRedisConnection(a.cfg.Redis)
	if err != nil {
		log.Fatal(err)
	}

	a.cacheClient = redisClient
}

func (a *Container) initNotifier() {
	if a.notifier != nil {
		return
	}

	//needed services
	notifiersConf := make(map[string]notification.NotifierConf)

	//init database conf
	notifiersConf[notification.DB_NOTIFIER] = &notification.DatabaseNotifierConf{
		TableName: "notifications",
		Db:        a.dbConn,
	}

	//init email conf
	notifiersConf[notification.EMAIL_NOTIFIER] = &notification.EmailNotifierConf{
		SmtpHost:        a.cfg.Email.SmtpHost,
		SmtpPort:        a.cfg.Email.SmtpPort,
		SmtpUsername:    a.cfg.Email.SmtpUsername,
		SmtpPassword:    a.cfg.Email.SmtpPassword,
		SmtpFromAddress: a.cfg.Email.SmtpFromAddress,
		SmtpEncryption:  a.cfg.Email.SmtpEncryption,
		SmtpFromName:    a.cfg.Email.SmtpFromName,
	}

	n, err := notification.NewNotifier(notifiersConf)
	if err != nil {
		log.Fatal(err)
	}

	a.notifier = n
}

func (a *Container) setAuthService() {
	if a.authService != nil {
		return
	}

	a.authService = services.NewAuthService(storage.NewUserRepo(a.dbConn), []byte(a.cfg.Server.TokenSecret),
		a.cfg.Server.TokenExpMinutes,
		a.cfg.Server.RefreshTokenExpMinutes)
}

func (a *Container) setBoardService() {
	if a.boardService != nil {
		return
	}
	a.boardService = services.NewBoardService(storage.NewBoardRepo(a.dbConn), storage.NewBoardMemberRepo(a.dbConn), storage.NewUserRepo(a.dbConn), storage.NewRoleRepo(a.dbConn))
}

func (a *Container) setTaskService() {
	if a.taskService != nil {
		return
	}
	taskRepository := storage.NewTaskRepo(a.dbConn)
	notifierAdapter := notifier.NewNotifierAdapter(a.notifier)
	a.taskService = services.NewTaskService(taskRepository, notifierAdapter, a.boardService)
}

func (a *Container) setColumnService() {
	if a.columnService != nil {
		return
	}
	a.columnService = services.NewColumnService(storage.NewColumnRepo(a.dbConn))
}
