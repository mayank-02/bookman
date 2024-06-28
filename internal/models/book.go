package models

import (
	"time"
)

type Book struct {
	ID            int       `json:"id"`
	Title         string    `json:"title"`
	Author        string    `json:"author"`
	PublishedDate string    `json:"published_date"`
	Edition       string    `json:"edition"`
	Description   string    `json:"description"`
	Genre         string    `json:"genre"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
