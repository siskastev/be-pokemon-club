package pokemon

import (
	"be-pokemon-club/internal/models"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
)

const POKE_API_URL = "https://pokeapi.co/api/v2/pokedex"

func GetAll(ctx context.Context, page, pageSize int) (result models.PokedexListResponse, err error) {
	if page < 1 {
		page = 1
	}

	offset := (page - 1) * pageSize
	limit := pageSize

	response, err := http.Get(fmt.Sprintf("%s?offset=%d&limit=%d", POKE_API_URL, offset, limit))
	if err != nil {
		return result, err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return result, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return result, err
	}

	result.Page = page
	result.PageSize = pageSize
	result.TotalPage = int(math.Ceil(float64(result.Count) / float64(pageSize)))

	return result, nil
}

func GetByIdOrName(ctx context.Context, idOrName string) (result models.PokedexDetailsResponse, err error) {
	response, err := http.Get(fmt.Sprintf("%s/%s", POKE_API_URL, idOrName))
	if err != nil {
		return result, err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return result, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return result, err
	}

	return result, nil
}
