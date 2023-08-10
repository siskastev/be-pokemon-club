package pokemon

import (
	"be-pokemon-club/internal/models"
	"context"
)

type Repository interface {
	CreateBattles(ctx context.Context, battle models.Battle) (models.Battle, error)
	GetBattlesByID(ctx context.Context, id int, createdBy string) (models.Battle, error)
	GetListBattles(ctx context.Context, request models.FilterPeriods, createdBy string) ([]models.Battle, error)
}
