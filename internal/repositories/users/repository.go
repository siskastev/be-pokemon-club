package users

import "be-pokemon-club/internal/models"

type Repository interface {
	GetLoginUser(request models.LoginRequest) (models.User, error)
}
