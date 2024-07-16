package check_status

import (
	"context"
)

type Logger interface {
	Save(message string, ID string, options ...any)
	SaveMessage(message string)
}

type Database interface {
	Update(context.Context, *Transaction) error
	GetID(ctx context.Context, ID string) (*Transaction, error)
}
