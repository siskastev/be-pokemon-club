package models

import "time"

type (
	PokedexListResponse struct {
		Count     int             `json:"count"`
		Page      int             `json:"page"`
		PageSize  int             `json:"page_size"`
		TotalPage int             `json:"total_page"`
		Results   []PokedexResult `json:"results"`
	}

	PokedexResult struct {
		Name string `json:"name"`
	}

	PokedexDetailsResponse struct {
		ID            int    `json:"id"`
		Name          string `json:"name"`
		IsMainSeries  bool   `json:"is_main_series"`
		VersionGroups []struct {
			Name string `json:"name"`
		} `json:"version_groups"`
		PokemonEntries []struct {
			EntryNumber    int `json:"entry_number"`
			PokemonSpecies struct {
				Name string `json:"name"`
			} `json:"pokemon_species"`
		} `json:"pokemon_entries"`
	}

	BattleResponse struct {
		ID                    int                     `json:"id"`
		TotalPokemon          int                     `json:"total_pokemon"`
		CreatedAt             time.Time               `json:"created_at,omitempty"`
		CreatedBy             string                  `json:"created_by"`
		UpdatedAt             *time.Time              `json:"updated_at,omitempty"`
		UpdatedBy             string                  `json:"updated_by,omitempty"`
		BattleDetailsResponse []BattleDetailsResponse `json:"battle_details"`
	}

	BattleDetailsResponse struct {
		ID    int    `json:"id"`
		Name  string `json:"name"`
		Score int    `json:"score"`
	}

	RankingResponse struct {
		PokemonName string `json:"pokemon_name"`
		TotalScore  int    `json:"total_score"`
		TotalBattle int    `json:"total_battle"`
	}
)
