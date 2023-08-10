package pokemon

import (
	"be-pokemon-club/internal/database"
	handler "be-pokemon-club/internal/handler/pokemon"
	"be-pokemon-club/internal/middleware"
	"be-pokemon-club/internal/redis"
	repository "be-pokemon-club/internal/repositories/pokemon"
	service "be-pokemon-club/internal/services/pokemon"
	"os"

	"github.com/gofiber/fiber/v2"
)

const (
	RoleBos  = "bos"
	RoleOps  = "ops"
	RoleTeam = "team"
)

func RegisterRoutes(route fiber.Router) {
	repo := repository.NewUserRepository(database.DB)
	pokemonService := service.NewPokemonService(redis.RedisClient, repo)
	handler := handler.NewHandlerPokemon(pokemonService)

	pokemonGroup := route.Group("/pokemon")

	jwtMiddleware := middleware.NewAuthMiddleware(os.Getenv("JWT_PRIVATE_KEY"))

	pokemonGroup.Use(jwtMiddleware.AuthRequired())

	pokemonGroup.Get("", handler.GetAll)
	pokemonGroup.Get("/ranking", jwtMiddleware.HasRoles(RoleBos), handler.GetRanking)
	pokemonGroup.Post("/battles", jwtMiddleware.HasRoles(RoleOps, RoleTeam), handler.CreateBattles)
	pokemonGroup.Get("/battles/:id", handler.GetBattlesByID)
	pokemonGroup.Get("/battles", handler.GetBattles)
	pokemonGroup.Get("/:id_or_name", handler.GetByIdOrName)
}
