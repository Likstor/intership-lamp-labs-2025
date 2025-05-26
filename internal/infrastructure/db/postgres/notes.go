package postgres

import (
	"context"
	cr "service/internal/domain/criteria"
	"service/internal/domain/entity"
	"service/internal/pkg/apperror"
	"service/internal/pkg/logs"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type NotesRepository struct {
	Pool *pgxpool.Pool
}

const opNoteGet = "infrastructure.db.postgres.Note.Get"

const queryNoteGet = `
	SELECT 
		id, title, content, created_at
	FROM
		public.notes
`

func (ns NotesRepository) Get(ctx context.Context, criteria *cr.Criteria) ([]entity.Note, error) {
	sql := queryNoteGet

	if criteria != nil {
		sql = cr.Build(queryNoteGet, criteria, 1, "$%d")
	}
	
	rows, err := ns.Pool.Query(ctx, sql)
	if err != nil {
		logs.Error(
			ctx,
			err.Error(),
			opNoteGet,
		)

		return nil, apperror.ErrRepository
	}

	notes, err := pgx.CollectRows(rows, pgx.RowToStructByNameLax[entity.Note])
	if err != nil {
		logs.Error(
			ctx,
			err.Error(),
			opNoteGet,
		)

		return nil, apperror.ErrRepository
	}

	return notes, nil
}

const opNoteCreate = "infrastructure.db.postgres.Note.Create"

const queryNoteCreate = `
	INSERT INTO public.notes
		(title, content)
	VALUES
		($1, $2)
`

func (ns NotesRepository) Create(ctx context.Context, note entity.Note) error {
	_, err := ns.Pool.Exec(ctx, queryNoteCreate, note.Title, note.Content)
	if err != nil {
		logs.Error(
			ctx,
			err.Error(),
			opNoteCreate,
		)

		return apperror.ErrRepository
	}

	return nil
}
