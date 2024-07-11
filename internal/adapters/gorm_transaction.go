package adapters

import (
	"context"
	"github.com/GoBootCamp-Group1/Task-Management/internal/core/ports"
	"gorm.io/gorm"
)

type GormTransaction struct {
	db *gorm.DB
	tx *gorm.DB
}

type ValueKeyType string

const (
	CtxKey ValueKeyType = "CTX-KEY"
)

func (t *GormTransaction) GetTX() any {
	return t.tx
}

func NewGormTransaction(db *gorm.DB) *GormTransaction {
	return &GormTransaction{db: db}
}

func (t *GormTransaction) Begin(ctx context.Context) (ports.Transaction, error) {
	tx := t.db.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &GormTransaction{db: t.db, tx: tx}, nil
}

func (t *GormTransaction) Commit(ctx context.Context) error {
	if t.tx != nil {
		return t.tx.Commit().Error
	}
	return nil
}

func (t *GormTransaction) Rollback(ctx context.Context) error {
	if t.tx != nil {
		return t.tx.Rollback().Error
	}
	return nil
}

func (t *GormTransaction) Execute(ctx context.Context, fn func(ctx context.Context) error) error {
	tx := t.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	transactionCtx := context.WithValue(ctx, CtxKey, tx)

	if err := fn(transactionCtx); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}
