package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"service/internal/delivery/http/v1/dto"
	"service/internal/delivery/http/v1/responses"
	"service/internal/domain/entity"
	"service/internal/pkg/apperror"
	"service/internal/pkg/logs"
	ucports "service/internal/usecase/ports"
	"strconv"
)

type NotesHandler struct {
	UseCase ucports.NotesUseCase
}

func (nh NotesHandler) GetMux() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /note", nh.createNote)
	mux.HandleFunc("GET /note/{id}", nh.getNote)
	mux.HandleFunc("GET /notes", nh.getNotes)

	return mux
}

func noteEntityToNoteDTO(note entity.Note) dto.Note {
	return dto.Note{
		ID:        note.ID,
		Title:     note.Title,
		Content:   note.Content,
		CreatedAt: note.CreatedAt,
	}
}

func noteDTOToNoteEntity(note dto.Note) entity.Note {
	return entity.Note{
		ID:        note.ID,
		Title:     note.Title,
		Content:   note.Content,
		CreatedAt: note.CreatedAt,
	}
}

func notesEntitiesToNotesPage(notes []entity.Note) dto.NotesPage {
	page := dto.NotesPage{
		Notes: make([]dto.Note, 0, len(notes)),
	}

	for _, note := range notes {
		page.Notes = append(page.Notes, noteEntityToNoteDTO(note))
	}

	return page
}

const opNotesGetNote = "delivery.http.v1.handlers.Notes.GetNote"

func (nh NotesHandler) getNote(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	if idStr == "" {
		responses.NotFound(r.Context(), w)

		return
	}

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		logs.Warn(
			r.Context(),
			err.Error(),
			opNotesGetNote,
		)

		responses.NotFound(r.Context(), w)

		return
	}

	note, err := nh.UseCase.GetNote(r.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, apperror.ErrNoteNotExists):
			responses.NotFound(r.Context(), w)
		default:
			responses.InternalServerError(r.Context(), w)
		}
		return
	}

	resp := noteEntityToNoteDTO(note)

	responses.JSON(
		r.Context(),
		w,
		http.StatusOK,
		resp,
	)
}

func (nh NotesHandler) getNotes(w http.ResponseWriter, r *http.Request) {
	notes, err := nh.UseCase.GetNotes(r.Context())
	if err != nil {
		responses.InternalServerError(r.Context(), w)

		return
	}

	resp := notesEntitiesToNotesPage(notes)

	responses.JSON(
		r.Context(),
		w,
		http.StatusOK,
		resp,
	)
}

const opNotesCreate = "delivery.http.v1.handlers.Notes.Create"

func (nh NotesHandler) createNote(w http.ResponseWriter, r *http.Request) {
	var noteDTO dto.Note

	if err := json.NewDecoder(r.Body).Decode(&noteDTO); err != nil {
		logs.Warn(
			r.Context(),
			err.Error(),
			opNotesCreate,
		)

		responses.Error(
			r.Context(),
			w,
			http.StatusBadRequest,
			http.StatusText(http.StatusBadRequest),
		)

		return
	}

	noteEntity := noteDTOToNoteEntity(noteDTO)

	if err := nh.UseCase.Create(r.Context(), noteEntity); err != nil {
		responses.InternalServerError(r.Context(), w)

		return
	}

	w.WriteHeader(http.StatusCreated)
}
