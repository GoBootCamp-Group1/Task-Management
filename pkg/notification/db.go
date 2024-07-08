package notification

import (
	"context"
	"fmt"
	"github.com/GoBootCamp-Group1/Task-Management/pkg/notification/notification"
	"gorm.io/gorm"
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

	// Check if the table exists
	tableName := cfg.TableName
	if db.Migrator().HasTable(tableName) {
		fmt.Printf("Table %s exists.\n", tableName)

		// Check for specific fields
		fieldsToCheck := []string{"id", "created_at", "deleted_at", "read_at", "user_id", "type", "data"}
		for _, field := range fieldsToCheck {
			if db.Migrator().HasColumn(tableName, field) {
				fmt.Printf("Field %s exists in table %s.\n", field, tableName)
			} else {
				fmt.Printf("Field %s does not exist in table %s.\n", field, tableName)
			}
		}
	} else {
		fmt.Printf("Table %s does not exist.\n", tableName)
	}

	return &DatabaseNotifierConfInternal{
		db: db,
	}, nil
}

func (s *DatabaseNotifierConfInternal) Send(ctx context.Context, userID uint, input *notification.InAppInput) error {
	fmt.Printf("Storing Notification inside database for user : %d :: %v", userID, input)
	return nil
}

func (s *DatabaseNotifierConfInternal) Read(ctx context.Context, notificationID uint) error {
	fmt.Printf("Reading Notification inside database for notification : %d ", notificationID)
	return nil
}

func (s *DatabaseNotifierConfInternal) Delete(ctx context.Context, notificationID uint) error {
	fmt.Printf("Deleting Notification inside database for notification : %d ", notificationID)
	return nil
}

//ctx := context.Background()
//noti := &notification2.InAppInput{
//	Type: "dsad",
//	Data: "some data",
//}
//notifier.InApp.Send(ctx, 1, noti)
