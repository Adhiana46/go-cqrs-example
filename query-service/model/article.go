package model

import "time"

type Article struct {
	ID        string    `bson:"_id,omitempty" json:"id"`
	Uuid      string    `bson:"uuid" json:"uuid"`
	Author    string    `bson:"author" json:"author"`
	Title     string    `bson:"title" json:"title"`
	Body      string    `bson:"body" json:"body"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}
