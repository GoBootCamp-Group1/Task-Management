package domains

import (
	"github.com/google/uuid"
	"time"
)

type Notification struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
	ReadAt    *time.Time
	UserID    uint
	Type      string
	Message   string

	User *User
}
