package models

import "time"

type LoginAttempts struct {
	ID              string     `json:"id,omitempty" bson:"_id,omitempty"`
	UserId          string     `json:"user_id" bson:"user_id"`
	Attempts        *int       `json:"attempts,omitempty" bson:"attempts,omitempty"`
	LastAttemptTime *time.Time `json:"last_attempt_time,omitempty" bson:"last_attempt_time,omitempty"`
	CreatedAt       *time.Time `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt       *time.Time `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}
