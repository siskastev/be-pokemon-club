package auth

import (
	"be-pokemon-club/internal/models"
	repository "be-pokemon-club/internal/repositories/users"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type userService struct {
	userRepo repository.Repository
}

func NewUserService(userRepo repository.Repository) Service {
	return &userService{userRepo: userRepo}
}

func verifyPassword(hashPasswordDb, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashPasswordDb), []byte(password))
	if err != nil {
		return err
	}
	return nil
}

func (u *userService) LoginUser(request models.LoginRequest) (models.UserResponse, error) {
	result, err := u.userRepo.GetLoginUser(request)
	if err != nil {
		return models.UserResponse{}, err
	}

	if err := verifyPassword(result.Password, request.Password); err != nil {
		return models.UserResponse{}, errors.New("password is incorrect")
	}

	response := models.UserResponse{
		ID:        result.ID,
		Name:      result.Name,
		Email:     result.Email,
		Role:      result.Role,
		CreatedAt: result.CreatedAt,
		UpdatedAt: result.UpdatedAt,
	}

	return response, nil
}
