package dto

import "time"

type Note struct {
	ID        uint64    `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

type NotesPage struct {
	Notes []Note `json:"notes"`
}

type CreateNote struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}
