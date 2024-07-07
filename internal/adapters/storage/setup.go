package storage

import (
	"fmt"
	"github.com/GoBootCamp-Group1/Task-Management/config"
	entities2 "github.com/GoBootCamp-Group1/Task-Management/internal/adapters/storage/entities"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewPostgresGormConnection(dbConfig config.DB) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=UTC", dbConfig.Host, dbConfig.User, dbConfig.Pass, dbConfig.DBName, dbConfig.Port)
	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}

func Migrate(db *gorm.DB) {
	migrator := db.Migrator()

	migrator.AutoMigrate(&entities2.User{})
}
