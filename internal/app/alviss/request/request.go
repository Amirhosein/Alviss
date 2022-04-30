package request

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type URLCreationRequest struct {
	LongURL string `json:"LongURL" binding:"required"`
	ExpDate string `json:"ExpTime" binding:"required"`
}

func (u URLCreationRequest) Validate() error {
	err := validation.ValidateStruct(&u,
		validation.Field(&u.LongURL, validation.Required, is.URL),
		validation.Field(&u.ExpDate, validation.Required),
	)
	if err != nil {
		return err
	}

	return nil
}
