package model

import "time"

// Event is a model for events in calendar with fields id, user_id, title and date
type Event struct {
	ID     uint64    `json:"uuid"`
	UserID uint64    `json:"user_id"`
	Title  string    `json:"title"`
	Date   time.Time `json:"date"`
}
