package project_models

import "time"

type ChannelResponse struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type ProjectResponse struct {
	ID        string             `json:"id"`
	Name      string             `json:"name"`
	CreatedAt time.Time          `json:"created_at"`
	Channels  []*ChannelResponse `json:"channels"`
}
