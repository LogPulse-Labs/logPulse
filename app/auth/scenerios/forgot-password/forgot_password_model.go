package forgot_password

type ForgotPasswordRequest struct {
	Email string `json:"email" validate:"required,email"`
}
