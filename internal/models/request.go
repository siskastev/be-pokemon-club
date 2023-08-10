package models

type (
	PokedexListRequest struct {
		Page     int `query:"page"`
		PageSize int `query:"page_size"`
	}

	BattleRequest struct {
		Pokemon []PokemonRequest `json:"pokemon" validate:"required,dive"`
	}

	PokemonRequest struct {
		Name  string `json:"name" validate:"required"`
		Score int    `json:"score" validate:"required"`
	}

	FilterPeriods struct {
		StartDate string `query:"start_date"`
		EndDate   string `query:"end_date"`
	}
)
