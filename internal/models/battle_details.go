package models

import "time"

type BattleDetails struct {
	ID          int       `gorm:"column:id;primaryKey;autoIncrement"`
	BattleID    int       `gorm:"column:battle_id;not null"`
	PokemonName string    `gorm:"column:pokemon_name;varchar(100);not null"`
	Score       int       `gorm:"column:score;not null"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt   time.Time `gorm:"column:updated_at;autoUpdateTime"`
}
