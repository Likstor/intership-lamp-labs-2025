package apperror

import "errors"

var (
	ErrUseCase = errors.New("usecase internal error")
)

// Notes
var (
	ErrNoteNotExists = errors.New("note not exists")
)