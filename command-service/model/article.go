package model

import "time"

type Article struct {
	ID        int       `db:"id"`
	Uuid      string    `db:"uuid"`
	Author    string    `db:"author"`
	Title     string    `db:"title"`
	Body      string    `db:"body"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
