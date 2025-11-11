package repository

import (
    "fmt"
	"gorm.io/gorm"
	"akichat/backend/internal/model"
)

type UserRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{DB: db}
}

// ユーザー登録
func (r *UserRepository) CreateUser(user *model.User) error {
    if r.DB == nil {
        fmt.Println("DB is nil")
    }
    return r.DB.Create(user).Error
}

func (repo *UserRepository) GetUserByUserID(userID uint) (*model.User, error) {
    var user model.User
    err := repo.DB.Where("id = ?", userID).First(&user).Error
    if err != nil {
        return nil, err
    }
    return &user, nil
}

func (repo *UserRepository) GetUserByEmail(email string) (*model.User, error) {
    var user model.User
    err := repo.DB.Where("email = ?", email).First(&user).Error
    if err != nil {
        return nil, err
    }
    return &user, nil
}

func (repo *UserRepository) GetUserByUserName(username string) (*model.User, error) {
    var user model.User
    err := repo.DB.Where("username = ?", username).First(&user).Error
    if err != nil {
        return nil, err
    }
    return &user, nil
}
