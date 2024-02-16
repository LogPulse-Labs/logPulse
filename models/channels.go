package models

import "time"

type Channel struct {
	ID        string    `json:"id" bson:"_id,omitempty"`
	UserID    string    `json:"-" bson:"user_id"`
	ProjectID string    `json:"project_id" bson:"project_id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at,omitempty" bson:"created_at"`
	UpdatedAt time.Time `json:"-" bson:"updated_at"`
}
