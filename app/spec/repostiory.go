package spec

import (
	"context"
)

type Repository[T any] interface {
	Create(ctx context.Context, entity *T) (*T, error)
	FindByID(ctx context.Context, id string) (*T, error)
	FindAll(ctx context.Context, queryOptions *QueryOptions, filter *Filter) ([]T, error)
	Count(ctx context.Context, filter *Filter) (int64, error)
	Update(ctx context.Context, entity *T) error
	Delete(ctx context.Context, id string) error
	Exists(ctx context.Context, id string) (bool, error)
}

type QueryOptions struct {
	Limit  int
	Offset int
}
