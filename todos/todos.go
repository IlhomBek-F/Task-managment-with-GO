package todos

import (
	"net/http"
	"time"

	"github.com/go-playground/validator"
	"github.com/labstack/echo"
)

type (
	Todo struct {
		ID        uint      `json:"id"`
		Title     string    `json:"title" validate:"required"`
		Completed *bool     `json:"completed" validate:"required"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}

	UpdateTodo struct {
		ID        uint   `json:"id" validate:"required,gt=0"`
		Title     string `json:"title" validate:"required"`
		Completed *bool  `json:"completed"`
	}

	CustomValidator struct {
		Validator *validator.Validate
	}
)

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.Validator.Struct(i); err != nil {
		// Optionally, you could return the error to give each route more control over the status code
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}
