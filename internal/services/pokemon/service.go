package pokemon

import (
	"be-pokemon-club/internal/models"
	"context"
)

type Service interface {
	GetAll(ctx context.Context, request models.PokedexListRequest) (models.PokedexListResponse, error)
	GetByIdOrName(ctx context.Context, idOrName string) (models.PokedexDetailsResponse, error)
	PokedexNameExist(ctx context.Context, request models.BattleRequest) ([]string, error)
	CreateBattles(ctx context.Context, request models.BattleRequest, userEmail string) (models.BattleResponse, error)
	GetBattlesByID(ctx context.Context, id int, createdBy string) (models.BattleResponse, error)
	GetListBattles(ctx context.Context, request models.FilterPeriods, createdBy string) ([]models.BattleResponse, error)
}
