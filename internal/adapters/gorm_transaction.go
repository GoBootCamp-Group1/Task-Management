package adapters

import (
	"context"
	"github.com/GoBootCamp-Group1/Task-Management/internal/core/ports"
	"gorm.io/gorm"
)

// GormTransaction is an implementation of the Transaction interface using GORM.
type GormTransaction struct {
	tx *gorm.DB
}

// Begin starts a new transaction.
func (g *GormTransaction) Begin(ctx context.Context) (ports.Transaction, error) {
	tx := g.tx.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &GormTransaction{tx: tx}, nil
}

// Commit commits the transaction.
func (g *GormTransaction) Commit() error {
	return g.tx.Commit().Error
}

// Rollback rolls back the transaction.
func (g *GormTransaction) Rollback() error {
	return g.tx.Rollback().Error
}
