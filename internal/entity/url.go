package entity

import "time"

type Url struct {
	ID          string    `json:"id" db:"id"`
	Title       string    `json:"title" db:"title"`
	Path        string    `json:"path" db:"path"`
	CountClick  int       `json:"count_clicks" db:"count_clicks"`
	Destination string    `json:"destination" db:"destination"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}
