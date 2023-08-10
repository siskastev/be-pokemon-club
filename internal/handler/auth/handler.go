package auth

import (
	"be-pokemon-club/internal/helpers/jwt"
	response "be-pokemon-club/internal/helpers/response"
	helpers "be-pokemon-club/internal/helpers/validator"
	"be-pokemon-club/internal/models"
	"be-pokemon-club/internal/services/auth"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type HandlerAuth struct {
	authService auth.Service
}

func NewHandlerAuth(authService auth.Service) *HandlerAuth {
	return &HandlerAuth{authService: authService}
}

func (h *HandlerAuth) Login(c *fiber.Ctx) error {
	logger := c.Locals("logger").(*logrus.Logger)

	var request models.LoginRequest

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request payload",
		})
	}

	if err := validator.New().Struct(request); err != nil {
		return helpers.ValidationErrorResponse(c, err.(validator.ValidationErrors))
	}

	userResponse, err := h.authService.LoginUser(request)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"method": c.Method(),
			"route":  c.Path(),
			"error":  err.Error(),
		}).Error("Failed to parse login user")
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}

	// Generate JWT token
	token, err := jwt.GenerateJWT(userResponse)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"method": c.Method(),
			"route":  c.Path(),
			"error":  err.Error(),
		}).Error("Failed to generate token")
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}

	result := models.UserResponseWithToken{
		UserResponse: userResponse,
		Token: models.UserSession{
			JWTToken: token,
		},
	}

	logger.WithFields(logrus.Fields{
		"method": c.Method(),
		"route":  c.Path(),
		"error":  nil,
	}).Info("Success login user")
  
	return response.Success(c, result)
}
