package notification

import (
	"fmt"
	"github.com/GoBootCamp-Group1/Task-Management/pkg/notification/notification"
	// Assume we have an SMS package for sending SMS
)

type DatabaseNotifierConf struct {
	NotifierConf
	// Add necessary fields, like SMS gateway info
}

func NewDatabaseNotifier(cfg *DatabaseNotifierConf) (*DatabaseNotifierConf, error) {
	return &DatabaseNotifierConf{}, nil
}

func (s *DatabaseNotifierConf) Send(userId uint, notification *notification.InAppNotificationInput) error {
	fmt.Printf("Storing Notification inside database")
	return nil
}
