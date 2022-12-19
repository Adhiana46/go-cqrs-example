package model

import "time"

type Article struct {
	ID        int       `db:"id" json:"id"`
	Uuid      string    `db:"uuid" json:"uuid"`
	Author    string    `db:"author" json:"author"`
	Title     string    `db:"title" json:"title"`
	Body      string    `db:"body" json:"body"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}
