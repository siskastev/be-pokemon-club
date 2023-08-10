package models

import "time"

type Battle struct {
	ID            int             `gorm:"column:id;primaryKey;autoIncrement"`
	TotalPokemon  int             `gorm:"column:total_pokemon;not null"`
	CreatedAt     time.Time       `gorm:"column:created_at;autoCreateTime"`
	CreatedBy     string          `gorm:"column:created_by;varchar(100);not null"`
	UpdatedAt     time.Time       `gorm:"column:updated_at;autoUpdateTime"`
	UpdatedBy     string          `gorm:"column:updated_by;varchar(100);not null"`
	BattleDetails []BattleDetails `gorm:"foreignKey:BattleID"`
}
