package auth

import "be-pokemon-club/internal/models"

type Service interface {
	LoginUser(request models.LoginRequest) (models.UserResponse, error)
}
