package repository

import (
	"github.com/FaisalMashuri/my-wallet/internal/domain/user"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) user.UserRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) FindUserByEmail(email string) (user *user.User, err error) {
	return user, nil
}

func (r *userRepository) CreateUser(user *user.User) (*user.User, error) {
	return user, nil
}

func (r *userRepository) UpdateUser(updatedUser *user.User) (*user.User, error) {
	return updatedUser, nil
}

func (r *userRepository) GetUserByID(id string) (user *user.User, err error) {
	return user, nil
}

func (r *userRepository) GetAllUser() (user []*user.User, err error) {
	return user, nil
}

func (r *userRepository) DeleteUser(id string) error {
	return nil
}
