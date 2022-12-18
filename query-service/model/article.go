package model

import "time"

type Article struct {
	ID        int       `bson:"_id,omitempty"`
	Uuid      string    `bson:"uuid"`
	Author    string    `bson:"author"`
	Title     string    `bson:"title"`
	Body      string    `bson:"body"`
	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
}
