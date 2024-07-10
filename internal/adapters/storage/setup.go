package storage

import (
	"fmt"

	"github.com/GoBootCamp-Group1/Task-Management/config"
	"github.com/GoBootCamp-Group1/Task-Management/internal/adapters/storage/entities"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewPostgresGormConnection(dbConfig config.DB) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=UTC", dbConfig.Host, dbConfig.User, dbConfig.Pass, dbConfig.DBName, dbConfig.Port)
	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}

func Migrate(db *gorm.DB) {
	migrator := db.Migrator()

	err := migrator.AutoMigrate(&entities.User{})
	if err != nil {
		panic("migration failed")
	}
}
