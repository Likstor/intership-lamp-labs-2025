package ucports

import (
	"context"
	"service/internal/domain/entity"
)

type NotesUseCase interface {
	GetNote(ctx context.Context, id uint64) (entity.Note, error)
	GetNotes(ctx context.Context) ([]entity.Note, error)
	Create(ctx context.Context, note entity.Note) error
}