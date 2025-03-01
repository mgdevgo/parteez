package shared

import "context"

type Repository[T any] interface {
	Save(ctx context.Context, entity T) error
	FindAll(ctx context.Context) ([]T, error)
	FindById(ctx context.Context, id int) (T, error)
	Delete(ctx context.Context, id int) error
}
