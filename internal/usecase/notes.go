package usecase

import (
	"context"
	cr "service/internal/domain/criteria"
	"service/internal/domain/entity"
	"service/internal/domain/ports/repository"
	"service/internal/pkg/apperror"
	"service/internal/pkg/logs"
)

var (
	dummyNote = entity.Note{}
)

type NotesService struct {
	Repository repository.Notes
}

const opNotesGetNote = "usecase.Notes.GetNote"

func (ns NotesService) GetNote(ctx context.Context, id uint64) (entity.Note, error) {
	criteria := &cr.Criteria{
		Condition: &cr.SimpleCondition{
			Field: "id",
			Operator: cr.EQ,
			Value: id,
		},
	}

	notes, err := ns.Repository.Get(ctx, criteria)
	if err != nil {
		return dummyNote, err
	}

	if len(notes) == 0 {
		logs.Warn(
			ctx,
			apperror.ErrNoteNotExists.Error(),
			opNotesGetNote,
		)
		
		return dummyNote, apperror.ErrNoteNotExists
	}

	return notes[0], nil
}

func (ns NotesService) GetNotes(ctx context.Context) ([]entity.Note, error) {
	notes, err := ns.Repository.Get(ctx, nil)
	if err != nil {
		return nil, err
	}
	
	return notes, nil
}

func (ns NotesService) Create(ctx context.Context, note entity.Note) error {
	if err := ns.Repository.Create(ctx, note); err != nil {
		return err
	}

	return nil
}

