package requests

import (
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
)

type (
	CalcReq struct {
		Type        int8    `json:"Type" query:"Type" validate:"required,number,min=1,max=20"`
		GateId      string  `json:"GateId" query:"GateId" validate:"required"`
		CountryIso  uint16  `json:"CountryIso" query:"CountryIso" validate:"required,number"`
		Weight      float32 `json:"Weight" query:"Weight" validate:"required,number"`
		OrderAmount float32 `json:"OrderAmount" query:"OrderAmount" validate:"required,number"`
	}

	calcValidator struct {
		Validator *validator.Validate
	}
)

func (cv *calcValidator) Validate(i interface{}) error {
	if err := cv.Validator.Struct(i); err != nil {
		return err
	}
	return nil
}

func CalcRegister(e *echo.Echo) {
	e.Validator = &calcValidator{Validator: validator.New()}
}
