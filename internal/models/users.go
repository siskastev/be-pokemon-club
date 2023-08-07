package models

import (
	"time"
)

type Role string

const (
	BOS_ROLE Role = "bos"
	OPS_ROLE Role = "ops"
	TEA_ROLE Role = "team"
)

type User struct {
	ID        int       `gorm:"column:id;primaryKey;autoIncrement"`
	Name      string    `gorm:"column:name;type:varchar(100);not null"`
	Email     string    `gorm:"column:email;type:varchar(100);unique;not null"`
	Role      Role      `gorm:"column:role;not null;default:'team'"`
	Password  string    `gorm:"column:password;type:varchar(100);not null"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type UserResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Role      Role      `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserSession struct {
	JWTToken string `json:"jwt_token"`
}

type UserResponseWithToken struct {
	UserResponse
	Token UserSession `json:"token"`
}
