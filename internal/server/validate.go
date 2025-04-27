package server

import (
	"github.com/go-playground/validator/v10"

	"github.com/trysourcetool/onprem-portal/internal/errdefs"
)

func validateRequest(p any) error {
	v := validator.New()

	if err := v.Struct(p); err != nil {
		return errdefs.ErrInvalidArgument(err)
	}

	return nil
}
