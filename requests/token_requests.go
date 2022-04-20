package requests

import validation "github.com/go-ozzo/ozzo-validation"

type BasicToken struct {
	Token  string `json:"token" validate:"required"`
	Active int16  `json:"active"`
}

func (bp BasicToken) Validate() error {
	return validation.ValidateStruct(&bp,
		validation.Field(&bp.Token, validation.Required),
	)
}

type CreateTokenRequest struct {
}

type UpdateTokenRequest struct {
	BasicToken
}
