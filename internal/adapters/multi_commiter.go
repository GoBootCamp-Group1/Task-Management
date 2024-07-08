package adapters

import (
	"github.com/GoBootCamp-Group1/Task-Management/pkg/valuecontext"
)

type multiCommitter struct {
	committers []valuecontext.Committer
	txs        []valuecontext.Committer
}

func NewMultiCommitter(committers []valuecontext.Committer) *multiCommitter {
	mc := &multiCommitter{
		committers: committers,
	}

	txs := make([]valuecontext.Committer, len(committers))
	mc.txs = txs

	return mc
}

func (c *multiCommitter) Begin() valuecontext.Committer {
	for i, cm := range c.committers {
		c.txs[i] = cm.Begin()
	}

	return c
}

func (c *multiCommitter) Commit() error {
	for _, tx := range c.txs {
		err := tx.Commit()
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *multiCommitter) Rollback() error {
	for _, tx := range c.txs {
		err := tx.Rollback()
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *multiCommitter) Tx() any {
	return c.txs
}
