package data

import "time"

type Movie struct {
	ID       int64     `json:"id"`
	Title    string    `json:"title"`
	CreateAt time.Time `json:"created_at"`
	Year     int32     `json:"year,omitempty"`
	Runtime  int32     `json:"runtime,omitempty"`
	Genres   []string  `json:"genres,omitempty"`
	Version  int32     `json:"version"`
}
