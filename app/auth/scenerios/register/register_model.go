package register_scenerio

type RegisterRequest struct {
	FullName     string `json:"full_name" validate:"required,min=2"`
	Organization string `json:"organization" validate:"required"`
	Email        string `json:"email" validate:"required,email"`
	Password     string `json:"password" validate:"required,alphanum,min=8"`
}
