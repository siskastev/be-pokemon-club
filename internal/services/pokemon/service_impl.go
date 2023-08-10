package pokemon

import (
	pokemonIntegrations "be-pokemon-club/internal/integrations/pokemon"
	"be-pokemon-club/internal/models"
	"be-pokemon-club/internal/repositories/pokemon"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

type pokemonService struct {
	redisClient *redis.Client
	pokemonRepo pokemon.Repository
}

func NewPokemonService(redisClient *redis.Client, pokemonRepo pokemon.Repository) Service {
	return &pokemonService{redisClient: redisClient,
		pokemonRepo: pokemonRepo,
	}
}

func (s *pokemonService) GetAll(ctx context.Context, request models.PokedexListRequest) (result models.PokedexListResponse, err error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	result, err = pokemonIntegrations.GetAll(ctx, request.Page, request.PageSize)
	if err != nil {
		return result, err
	}

	return result, nil
}

func (s *pokemonService) GetByIdOrName(ctx context.Context, idOrName string) (result models.PokedexDetailsResponse, err error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	result, err = pokemonIntegrations.GetByIdOrName(ctx, idOrName)
	if err != nil {
		return result, err
	}

	return result, nil
}

func (s *pokemonService) fetchPokemonNamesFromAPI(ctx context.Context) (map[string]bool, error) {
	result, err := pokemonIntegrations.GetAll(ctx, 1, 1000)
	if err != nil {
		return nil, errors.New("failed to get pokemon list")
	}

	pokemonNameMap := make(map[string]bool)
	for _, pokemon := range result.Results {
		pokemonNameMap[pokemon.Name] = true
	}

	return pokemonNameMap, nil
}

func (s *pokemonService) getCachedPokemonNames(ctx context.Context) (map[string]bool, error) {
	cacheKey := "pokemonNames"
	if s.redisClient == nil {
		return nil, errors.New("redis client is nil")
	}

	cachedPokemonNames, err := s.redisClient.Get(ctx, cacheKey).Result()
	if err != nil && err != redis.Nil {
		return nil, errors.New("failed to get cached pokemon names")
	}

	if cachedPokemonNames == "" {
		pokemonNameMap, err := s.fetchPokemonNamesFromAPI(ctx)
		if err != nil {
			return nil, err
		}

		pokemonNamesJSON, err := json.Marshal(pokemonNameMap)
		if err != nil {
			return nil, err
		}

		err = s.redisClient.Set(ctx, cacheKey, pokemonNamesJSON, 1*time.Hour).Err()
		if err != nil {
			return nil, err
		}

		return pokemonNameMap, nil
	}

	var pokemonNameMap map[string]bool
	err = json.Unmarshal([]byte(cachedPokemonNames), &pokemonNameMap)
	if err != nil {
		return nil, err
	}

	return pokemonNameMap, nil
}

func (s *pokemonService) PokedexNameExist(ctx context.Context, request models.BattleRequest) ([]string, error) {
	cachedNames, err := s.getCachedPokemonNames(ctx)
	if err != nil {
		return nil, err
	}

	var invalidPokemonName []string
	for _, name := range request.Pokemon {
		if !cachedNames[name.Name] {
			invalidPokemonName = append(invalidPokemonName, fmt.Sprintf("pokemon with name %s not found", name.Name))
		}
	}
	return invalidPokemonName, nil
}

func (s *pokemonService) CreateBattles(ctx context.Context, request models.BattleRequest, userEmail string) (models.BattleResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	detailBattle := []models.BattleDetails{}

	for _, v := range request.Pokemon {
		detail := models.BattleDetails{
			PokemonName: v.Name,
			Score:       v.Score,
		}
		detailBattle = append(detailBattle, detail)
	}

	battle := models.Battle{
		TotalPokemon:  len(request.Pokemon),
		CreatedBy:     userEmail,
		BattleDetails: detailBattle,
	}

	result, err := s.pokemonRepo.CreateBattles(ctx, battle)
	if err != nil {
		return models.BattleResponse{}, err
	}

	response := models.BattleResponse{
		ID:                    result.ID,
		TotalPokemon:          result.TotalPokemon,
		CreatedAt:             result.CreatedAt,
		CreatedBy:             result.CreatedBy,
		BattleDetailsResponse: make([]models.BattleDetailsResponse, len(result.BattleDetails)),
	}

	for i, v := range result.BattleDetails {
		response.BattleDetailsResponse[i] = models.BattleDetailsResponse{
			ID:    v.ID,
			Name:  v.PokemonName,
			Score: v.Score,
		}
	}

	// Clear cache for all battles
	keys, err := s.redisClient.Keys(ctx, "battles:*").Result()
	if err != nil {
		return models.BattleResponse{}, err
	}

	if len(keys) > 0 {
		result := s.redisClient.Del(ctx, keys...)
		if result.Err() != nil {
			return models.BattleResponse{}, result.Err()
		}
	}

	return response, nil
}

func (s *pokemonService) GetBattlesByID(ctx context.Context, id int, createdBy string) (models.BattleResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	result, err := s.pokemonRepo.GetBattlesByID(ctx, id, createdBy)
	if err != nil {
		return models.BattleResponse{}, err
	}

	response := models.BattleResponse{
		ID:                    result.ID,
		TotalPokemon:          result.TotalPokemon,
		CreatedAt:             result.CreatedAt,
		CreatedBy:             result.CreatedBy,
		UpdatedBy:             result.UpdatedBy,
		BattleDetailsResponse: make([]models.BattleDetailsResponse, len(result.BattleDetails)),
	}

	for i, v := range result.BattleDetails {
		response.BattleDetailsResponse[i] = models.BattleDetailsResponse{
			ID:    v.ID,
			Name:  v.PokemonName,
			Score: v.Score,
		}
	}

	return response, nil
}

func (s *pokemonService) GetListBattles(ctx context.Context, request models.FilterPeriods, createdBy string) ([]models.BattleResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	cacheKey := fmt.Sprintf("battles:start_date=%s:end_date=%s:created_by=%s", request.StartDate, request.EndDate, createdBy)
	if cachedResult, found := s.redisClient.Get(ctx, cacheKey).Result(); found == nil {
		var response []models.BattleResponse
		err := json.Unmarshal([]byte(cachedResult), &response)
		if err != nil {
			return nil, err
		}
		return response, nil
	}

	result, err := s.pokemonRepo.GetListBattles(ctx, request, createdBy)
	if err != nil {
		return nil, err
	}

	response := make([]models.BattleResponse, len(result))

	for i, v := range result {
		response[i] = models.BattleResponse{
			ID:                    v.ID,
			TotalPokemon:          v.TotalPokemon,
			CreatedAt:             v.CreatedAt,
			CreatedBy:             v.CreatedBy,
			UpdatedBy:             v.UpdatedBy,
			BattleDetailsResponse: make([]models.BattleDetailsResponse, len(v.BattleDetails)),
		}

		for j, detail := range v.BattleDetails {
			response[i].BattleDetailsResponse[j] = models.BattleDetailsResponse{
				ID:    detail.ID,
				Name:  detail.PokemonName,
				Score: detail.Score,
			}
		}
	}

	cacheValue, err := json.Marshal(response)
	if err != nil {
		return nil, err
	}

	err = s.redisClient.Set(ctx, cacheKey, cacheValue, 1*time.Hour).Err()
	if err != nil {
		return nil, err
	}

	return response, nil
}
