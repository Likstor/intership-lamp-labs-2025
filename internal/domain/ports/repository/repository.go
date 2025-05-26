package repository

import (
	"context"
	cr "service/internal/domain/criteria"
	"service/internal/domain/entity"
)

type Notes interface {
	Get(ctx context.Context, criteria *cr.Criteria) ([]entity.Note, error)
	Create(ctx context.Context, note entity.Note) error
}