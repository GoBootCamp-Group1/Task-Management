package ports

import "context"

type Transaction interface {
	Begin(ctx context.Context) (Transaction, error)
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
	Execute(ctx context.Context, fn func(ctx context.Context) error) error
	GetTX() any
}
