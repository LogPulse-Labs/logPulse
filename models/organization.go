package models

import "time"

type Organization struct {
	ID        string    `json:"id,omitempty" bson:"_id,omitempty"`
	UserID    string    `json:"user_id,omitempty" bson:"user_id"`
	Name      string    `json:"name,omitempty" bson:"name"`
	CreatedAt time.Time `json:"created_at,omitempty" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at,omitempty" bson:"updated_at"`
}
