package routes

import (
	"be-pokemon-club/internal/routes/auth"
	"be-pokemon-club/internal/routes/pokemon"

	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {
	api := app.Group("/api")
	//add route here
	auth.RegisterRoutes(api)
	pokemon.RegisterRoutes(api)

}
