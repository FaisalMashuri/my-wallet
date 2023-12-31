package repository

import (
	"github.com/FaisalMashuri/my-wallet/internal/domain/user"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"time"
)

type userRepository struct {
	db  *gorm.DB
	log *logrus.Logger
}

func NewRepository(db *gorm.DB, log *logrus.Logger) user.UserRepository {
	return &userRepository{
		db:  db,
		log: log,
	}
}

func (r *userRepository) FindUserByEmail(email string) (user *user.User, err error) {
	r.log.Debug("Start find user")
	err = r.db.Debug().Take(&user, "email = ?", email).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			r.log.Info("User not found")
			return nil, nil
		}
		r.log.WithField("error", err.Error()).Info("failed find user by email")
		return nil, err
	}
	r.log.Debug("User found")

	return user, nil
}

func (r *userRepository) CreateUser(user *user.User) (*user.User, error) {
	r.log.Debug("Start creating user")
	err := r.db.Debug().Create(&user).Error
	if err != nil {
		r.log.WithField("error", err.Error()).Info("failed create user")
		return nil, err
	}
	r.log.Debug("User created")
	return user, nil
}

func (r *userRepository) UpdateUser(updatedUser *user.User) (*user.User, error) {
	return updatedUser, nil
}

func (r *userRepository) GetUserByID(id string) (user *user.User, err error) {
	err = r.db.Debug().First(&user, "id = ?", id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return user, nil
		}
		return nil, err
	}
	return user, nil
}

func (r *userRepository) GetAllUser() (user []*user.User, err error) {
	return user, nil
}

func (r *userRepository) DeleteUser(id string) error {
	return nil
}

func (r *userRepository) VerifyUser(id string) error {
	err := r.db.Debug().Model(&user.User{}).Where("id = ?", id).Update("verified_at", time.Now()).Error
	if err != nil {
		return err
	}
	return nil
}
