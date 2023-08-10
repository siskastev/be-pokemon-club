package pokemon

import (
	response "be-pokemon-club/internal/helpers/response"
	helpers "be-pokemon-club/internal/helpers/validator"
	"be-pokemon-club/internal/models"
	"be-pokemon-club/internal/services/pokemon"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type HandlerPokemon struct {
	pokemonService pokemon.Service
}

func NewHandlerPokemon(pokemonService pokemon.Service) *HandlerPokemon {
	return &HandlerPokemon{pokemonService: pokemonService}
}

func (h *HandlerPokemon) GetAll(c *fiber.Ctx) error {
	logger := c.Locals("logger").(*logrus.Logger)

	request := new(models.PokedexListRequest)
	request.Page = c.QueryInt("page", 1)
	request.PageSize = c.QueryInt("page_size", 20)

	if err := c.QueryParser(request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"errors": err.Error()})
	}

	pokemonList, err := h.pokemonService.GetAll(c.Context(), *request)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"method": c.Method(),
			"route":  c.Path(),
			"error":  err.Error(),
		}).Error("Failed to get pokemon list")
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": pokemonList.Results,
		"metadata": fiber.Map{
			"page":       pokemonList.Page,
			"page_size":  pokemonList.PageSize,
			"total":      pokemonList.Count,
			"total_page": pokemonList.TotalPage,
		},
	})
}

func (h *HandlerPokemon) GetByIdOrName(c *fiber.Ctx) error {
	logger := c.Locals("logger").(*logrus.Logger)

	pokemonList, err := h.pokemonService.GetByIdOrName(c.Context(), c.Params("id_or_name"))
	if err != nil {
		logger.WithFields(logrus.Fields{
			"method": c.Method(),
			"route":  c.Path(),
			"error":  err.Error(),
		}).Error("Failed to get pokemon list")
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}

	return response.Success(c, pokemonList)
}

func (h *HandlerPokemon) CreateBattles(c *fiber.Ctx) error {
	logger := c.Locals("logger").(*logrus.Logger)

	request := new(models.BattleRequest)
	if err := c.BodyParser(request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"errors": err.Error()})
	}

	if err := validator.New().Struct(request); err != nil {
		return helpers.ValidationErrorResponse(c, err.(validator.ValidationErrors))
	}

	//Check if each pokemon name exists in pokedex
	var errors []string
	invalidPokemonName, err := h.pokemonService.PokedexNameExist(c.Context(), *request)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"method": c.Method(),
			"route":  c.Path(),
			"error":  err.Error(),
		}).Error("Failed to get pokemon list")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"errors": err.Error()})
	}

	errors = append(errors, invalidPokemonName...)
	if len(errors) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"errors": fiber.Map{"field": "pokemon", "message": errors}})
	}

	numPlayers := len(request.Pokemon)

	// check for unique scores and within range
	scoreMap := make(map[int]bool)
	for _, pokemon := range request.Pokemon {
		if pokemon.Score < 1 || pokemon.Score > numPlayers {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"errors": fiber.Map{"field": "pokemon", "message": "Scores must be between 1 and number of players"}})
		}
		if _, exists := scoreMap[pokemon.Score]; exists {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"errors": fiber.Map{"field": "pokemon", "message": "Scores must be unique"}})
		}
		scoreMap[pokemon.Score] = true
	}

	createdBy := c.Locals("user").(*models.UserResponse).Email
	result, err := h.pokemonService.CreateBattles(c.Context(), *request, createdBy)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"method": c.Method(),
			"route":  c.Path(),
			"error":  err.Error(),
		}).Error("Failed to get pokemon list")
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"data": result})
}

func (h *HandlerPokemon) GetBattlesByID(c *fiber.Ctx) error {
	logger := c.Locals("logger").(*logrus.Logger)

	id, _ := strconv.Atoi(c.Params("id"))

	role := c.Locals("user").(*models.UserResponse).Role
	createdBy := c.Locals("user").(*models.UserResponse).Email
	if role == "bos" {
		createdBy = ""
	}
	result, err := h.pokemonService.GetBattlesByID(c.Context(), id, createdBy)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"method": c.Method(),
			"route":  c.Path(),
			"error":  err.Error(),
		}).Error("Failed to get pokemon list")
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}

	return response.Success(c, result)
}

func (h *HandlerPokemon) GetBattles(c *fiber.Ctx) error {

	logger := c.Locals("logger").(*logrus.Logger)

	request := new(models.FilterPeriods)

	if err := c.QueryParser(request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"errors": err.Error()})
	}

	if request.StartDate != "" && !isValidDate(request.StartDate) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"errors": "Invalid start date format, use YYYY-MM-DD"})
	}

	if request.EndDate != "" && !isValidDate(request.EndDate) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"errors": "Invalid end date format, use YYYY-MM-DD"})
	}

	// Set default end date to the end of the current month
	if request.EndDate == "" {
		now := time.Now()
		startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
		endOfMonth := startOfMonth.AddDate(0, 1, -1)
		request.EndDate = endOfMonth.Format("2006-01-02")
	}

	role := c.Locals("user").(*models.UserResponse).Role
	createdBy := c.Locals("user").(*models.UserResponse).Email
	if role == "bos" {
		createdBy = ""
	}

	result, err := h.pokemonService.GetListBattles(c.Context(), *request, createdBy)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"method": c.Method(),
			"route":  c.Path(),
			"error":  err,
		}).Error("Failed to get pokemon list")
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}

	return response.Success(c, result)
}

func isValidDate(date string) bool {
	_, err := time.Parse("2006-01-02", date)
	return err == nil
}

func (h *HandlerPokemon) GetRanking(c *fiber.Ctx) error {

	logger := c.Locals("logger").(*logrus.Logger)

	result, err := h.pokemonService.GetRanking(c.Context())
	if err != nil {
		logger.WithFields(logrus.Fields{
			"method": c.Method(),
			"route":  c.Path(),
			"error":  err,
		}).Error("Failed to get pokemon list")
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}

	return response.Success(c, result)
}
