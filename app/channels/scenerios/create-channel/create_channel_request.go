package create_channel

type CreateChannelRequest struct {
	Name string `json:"name" validate:"required"`
}
