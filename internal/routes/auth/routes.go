package auth

import (
	"be-pokemon-club/internal/database"
	handler "be-pokemon-club/internal/handler/auth"
	repository "be-pokemon-club/internal/repositories/users"
	service "be-pokemon-club/internal/services/auth"

	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(route fiber.Router) {
	repo := repository.NewUserRepository(database.DB)
	service := service.NewUserService(repo)
	handler := handler.NewHandlerAuth(service)
	route.Post("/login", handler.Login)
}
