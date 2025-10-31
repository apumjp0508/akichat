package repository

import (
	"fmt"
	"gorm.io/gorm"
	"akichat/backend/internal/model"
)

type FriendShipRepository struct {
	DB *gorm.DB
}

func NewFriendShipRepository(db *gorm.DB) *FriendShipRepository {
	return &FriendShipRepository{DB: db}
}

// 友達リストの取得
func (r *FriendShipRepository) GetFriendsByUserID(userID uint) ([]model.User, error) {
	var friendships []model.FriendShip
	if r.DB == nil {
		fmt.Println("DB is nil")
	}
	err := r.DB.
		Preload("User1").
		Preload("User2").
		Where("user1_id = ? OR user2_id = ?", userID, userID).
		Find(&friendships).Error
	if err != nil {
		return nil, err
	}

	var friendUsers []model.User
	for _, f := range friendships {
		if f.User1ID == userID {
			friendUsers = append(friendUsers, f.User2)
		} else {
			friendUsers = append(friendUsers, f.User1)
		}
	}

	return friendUsers, nil
}

// 友達追加
func (r *FriendShipRepository) AddFriend(user1ID uint, user2ID uint) error {
	friendship := model.FriendShip{
		User1ID: user1ID,
		User2ID: user2ID,
	}
	if r.DB == nil {
		fmt.Println("DB is nil")
	}
	return r.DB.Create(&friendship).Error
}
