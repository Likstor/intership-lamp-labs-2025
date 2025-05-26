package entity

import "time"

type Note struct {
	ID        uint64
	Title     string
	Content   string
	CreatedAt time.Time
}