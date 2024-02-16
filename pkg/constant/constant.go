package constant

import "errors"

var (
	ErrSomethingWentWrong = errors.New("Something went wrong. Please try again later.")
	UnableToGenerateToken = errors.New("Unable to generate access token. Please try login later.")
)
