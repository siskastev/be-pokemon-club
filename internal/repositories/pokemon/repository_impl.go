package pokemon

import (
	"be-pokemon-club/internal/models"
	"context"
	"fmt"

	"gorm.io/gorm"
)

type pokemonRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) Repository {
	return &pokemonRepository{db: db}
}

func (o *pokemonRepository) CreateBattles(ctx context.Context, battle models.Battle) (models.Battle, error) {
	var created models.Battle
	err := o.db.Transaction(func(tx *gorm.DB) error {

		if err := tx.Create(&battle).Error; err != nil {
			return err
		}
		created = battle
		return nil
	})

	if err != nil {
		return created, err
	}

	return battle, nil
}

func (o *pokemonRepository) GetBattlesByID(ctx context.Context, id int, createdBy string) (models.Battle, error) {
	var battle models.Battle

	query := o.db.Where("id = ?", id).Preload("BattleDetails", func(db *gorm.DB) *gorm.DB {
		return o.db.Order("battle_details.score DESC")
	})

	if createdBy != "" {
		query = query.Where("id = ? AND created_by = ?", id, createdBy)
	}

	if err := query.First(&battle).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return battle, err
		}
	}

	return battle, nil
}

func (o *pokemonRepository) GetListBattles(ctx context.Context, request models.FilterPeriods, createdBy string) ([]models.Battle, error) {
	var battles []models.Battle

	query := o.db.Preload("BattleDetails", func(db *gorm.DB) *gorm.DB {
		return o.db.Order("battle_details.score DESC")
	})

	if request.StartDate != "" {
		query = query.Where("created_at >= ?", request.StartDate)
	}

	if request.EndDate != "" {
		query = query.Where("created_at <= ?", request.EndDate)
	}

	if createdBy != "" {
		query = query.Where("created_by = ?", createdBy)
	}

	fmt.Println("query", query)

	if err := query.Find(&battles).Error; err != nil {
		return battles, err
	}

	return battles, nil
}

func (o *pokemonRepository) GetListBattleDetails(ctx context.Context) ([]models.BattleDetails, error) {
	var battleDetails []models.BattleDetails

	if err := o.db.Find(&battleDetails).Error; err != nil {
		return battleDetails, err
	}

	return battleDetails, nil
}
