package models

import "time"

type Project struct {
	ID        string    `json:"id,omitempty" bson:"_id,omitempty"`
	UserID    string    `json:"user_id,omitempty" bson:"user_id"`
	Name      string    `json:"name,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty" bson:"created_at"`
	UpdatedAt time.Time `json:"-" bson:"updated_at"`
}
