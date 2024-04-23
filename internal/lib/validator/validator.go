package validator

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		echo.New().AcquireContext().Logger().Errorf("Failed to validate request body: %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, "Bad Request")
	}
	return nil
}

func SetCustomValidator() *CustomValidator {
	return &CustomValidator{
		validator: validator.New(validator.WithRequiredStructEnabled()),
	}
}
