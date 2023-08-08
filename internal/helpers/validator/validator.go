package validator

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type validationErrors struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func ValidationErrorResponse(c *fiber.Ctx, errs validator.ValidationErrors) error {
	var validationErrors []validationErrors
	for _, err := range errs {
		validationMessage := err.Tag()

		if err.Tag() == "email" {
			validationMessage = "Invalid email format"
		}

		validationErrors = append(validationErrors, struct {
			Field   string `json:"field"`
			Message string `json:"message"`
		}{
			Field:   err.Field(),
			Message: validationMessage,
		})
	}

	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
		"errors": validationErrors,
	})
}
