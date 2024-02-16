package models

import (
	"time"
)

type Timestamp struct {
	CreatedAt time.Time `json:"created_at,omitempty" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at,omitempty" bson:"updated_at"`
}

type User struct {
	ID            string    `json:"id,omitempty" bson:"_id,omitempty"`
	FullName      string    `json:"full_name" bson:"full_name"`
	Email         string    `json:"email" bson:"email"`
	EmailVerified bool      `json:"email_verified,omitempty" bson:"email_verified"`
	Password      string    `json:"-" bson:"password"`
	CreatedAt     time.Time `json:"created_at,omitempty" bson:"created_at"`
	UpdatedAt     time.Time `json:"updated_at,omitempty" bson:"updated_at"`
}
