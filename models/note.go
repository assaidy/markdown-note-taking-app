package models

import (
	"time"
)

type Note struct {
	Id        int       `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"-"`
	CreatedAt time.Time `json:"createdAt"`
}
