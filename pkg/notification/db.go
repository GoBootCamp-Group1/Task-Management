package notification

import (
	"context"
	"errors"
	"fmt"
	"github.com/GoBootCamp-Group1/Task-Management/pkg/notification/notification"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"regexp"
)

type DatabaseNotifierConf struct {
	NotifierConf
	TableName string
	Db        *gorm.DB
}

type DatabaseNotifierConfInternal struct {
	NotifierConf
	TableName string
	db        *gorm.DB
}

func NewDatabaseNotifier(cfg *DatabaseNotifierConf) (*DatabaseNotifierConfInternal, error) {
	//init database
	db := cfg.Db

	// Validate table name
	validTableName := regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*$`).MatchString
	if !validTableName(cfg.TableName) {
		return nil, errors.New("invalid table name")
	}

	// Check if the table exists
	tableName := cfg.TableName
	if db.Migrator().HasTable(tableName) {
		//fmt.Printf("Table %s exists.\n", tableName)

		// Check for specific fields
		fieldsToCheck := []string{"id", "created_at", "deleted_at", "read_at", "user_id", "type", "message"}
		for _, field := range fieldsToCheck {
			if db.Migrator().HasColumn(tableName, field) {
				//fmt.Printf("Field %s exists in table %s.\n", field, tableName)
			} else {
				fmt.Printf("Field %s does not exist in table %s.\n", field, tableName)
			}
		}
	} else {
		fmt.Printf("Table %s does not exist.\n", tableName)
	}

	return &DatabaseNotifierConfInternal{
		db:        db,
		TableName: tableName,
	}, nil
}

func (s *DatabaseNotifierConfInternal) Send(ctx context.Context, userID uint, input *notification.InAppInput) error {
	id := uuid.NewString()
	sql := fmt.Sprintf("INSERT INTO %s (id, created_at, updated_at, user_id, type, message) VALUES (?, NOW(), NOW(), ?, ?, ?)", s.TableName)
	err := s.db.WithContext(ctx).Exec(sql, id, userID, input.Type, input.Message).Error

	if err != nil {
		return err
	}

	return nil
}

func (s *DatabaseNotifierConfInternal) Read(ctx context.Context, notificationID string) error {
	fmt.Printf("Reading Notification inside database for notification : %s ", notificationID)

	sql := "UPDATE ? SET updated_at = NOW(), read_at = NOW() WHERE id = ?"
	err := s.db.WithContext(ctx).Exec(sql, s.TableName, notificationID).Error

	if err != nil {
		return err
	}
	return nil
}

func (s *DatabaseNotifierConfInternal) Delete(ctx context.Context, notificationID string) error {
	fmt.Printf("Deleting Notification inside database for notification : %s", notificationID)
	sql := "UPDATE ? SET updated_at = NOW(), deleted_at = NOW() WHERE id = ?"
	err := s.db.WithContext(ctx).Exec(sql, s.TableName, notificationID).Error

	if err != nil {
		return err
	}
	return nil
}

//ctx := context.Background()
//noti := &notification2.InAppInput{
//	Type: "dsad",
//	Data: "some data",
//}
//notifier.InApp.Send(ctx, 1, noti)
