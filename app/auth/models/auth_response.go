package auth_models

import (
	"time"
)

type AuthResponsePayload struct {
	ID           string    `json:"id,omitempty"`
	FullName     string    `json:"full_name"`
	Email        string    `json:"email"`
	CreatedAt    time.Time `json:"created_at,omitempty"`
	Organization struct {
		ID   string `json:"id,omitempty"`
		Name string `json:"name,omitempty"`
	} `json:"organization,omitempty"`
}
