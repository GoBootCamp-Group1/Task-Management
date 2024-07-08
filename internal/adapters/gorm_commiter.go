package adapters

import (
	"github.com/GoBootCamp-Group1/Task-Management/pkg/valuecontext"

	"gorm.io/gorm"
)

type GormCommitter struct {
	db *gorm.DB
	tx *gorm.DB
}

func NewGormCommitter(db *gorm.DB) valuecontext.Committer {
	return &GormCommitter{db: db}
}

func (c *GormCommitter) Tx() any {
	return c.tx
}

func (c *GormCommitter) Begin() valuecontext.Committer {
	tx := c.db.Begin()
	c.tx = tx
	return c
}

func (c *GormCommitter) Commit() error {
	if c.tx == nil {
		return nil
	}
	return c.tx.Commit().Error
}

func (c *GormCommitter) Rollback() error {
	if c.tx == nil {
		return nil
	}
	return c.tx.Rollback().Error
}
