package ports

import "context"

type Transaction interface {
	Begin(ctx context.Context) (Transaction, error)
	Commit() error
	Rollback() error
}
