package users

import (
	"be-pokemon-club/internal/models"

	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) Repository {
	return &userRepository{db: db}
}

func (u *userRepository) GetLoginUser(request models.LoginRequest) (models.User, error) {
	var user models.User
	if err := u.db.Where(models.User{Email: request.Email}).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return user, err
		}
	}
	return user, nil
}
