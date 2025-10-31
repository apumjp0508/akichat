package repository

import (
	"fmt"
	"gorm.io/gorm"
	"akichat/backend/internal/model"
)	

type FriendRequestRepository struct {
	DB *gorm.DB
}

func NewFriendRequestRepository(db *gorm.DB) *FriendRequestRepository {
	return &FriendRequestRepository{DB: db}
}

//フレンド申請作成
func (r *FriendRequestRepository) CreateFriendRequest(requestUserID uint, ReceiverUserID uint) error {
	// 重複チェック
	results := r.DuplicateRequestCheckByID(requestUserID, ReceiverUserID)
	if results {
		return fmt.Errorf("Friend request already exists")
	}
	// 逆パターンの重複チェック
	results_part2 := r.DuplicateRequestCheckByID(ReceiverUserID, requestUserID)
	if results_part2 {
		return fmt.Errorf("Friend request already exists in reverse direction")
	}

	friendRequest := model.FriendRequest{
		RequestUserID: requestUserID,
		ReceiverUserID: ReceiverUserID,
		Status: "pending",
	}
	if r.DB == nil {
		fmt.Println("DB is nil")
	}
	return r.DB.Create(&friendRequest).Error
}

func (r *FriendRequestRepository) GetFriendRequest(requestUserID uint, ReceiverUserID uint) (*model.FriendRequest, error) {
	var friendRequest model.FriendRequest
	if r.DB == nil {
		fmt.Println("DB is nil")
	}
	err := r.DB.Where("request_user_id = ? AND Receiver_user_id = ?", requestUserID, ReceiverUserID).First(&friendRequest).Error
	if err != nil {
		return nil, err
	}
	return &friendRequest, nil
}

func (r *FriendRequestRepository) DuplicateRequestCheckByID(requestUserID uint, ReceiverUserID uint) bool {
	var friendRequest model.FriendRequest
	if r.DB == nil {
		fmt.Println("DB is nil")
	}
	err := r.DB.Where("request_user_id = ? AND Receiver_user_id = ?", requestUserID, ReceiverUserID).First(&friendRequest).Error
	if err != nil {
		return false
	}
	return true
}

//フレンド申請削除
func (r *FriendRequestRepository) DeleteFriendRequest(requestUserID uint, ReceiverUserID uint) error {
	if r.DB == nil {
		fmt.Println("DB is nil")
	}
	return r.DB.Where("request_user_id = ? AND Receiver_user_id = ?", requestUserID, ReceiverUserID).Delete(&model.FriendRequest{}).Error
}

//ユーザーIDからフレンド申請取得
func (r *FriendRequestRepository) GetFriendRequestsByUserID(userID uint) ([]model.FriendRequest, error) {
	var friendRequests []model.FriendRequest
	if r.DB == nil {
		fmt.Println("DB is nil")
	}
	err := r.DB.Preload("RequestUser").Preload("ReceiverUser").Where("Receiver_user_id = ? AND status = ?", userID, "pending").Find(&friendRequests).Error
	if err != nil {
		return nil, err
	}
	return friendRequests, nil
}

//フレンド申請ステータス更新
func (r *FriendRequestRepository) UpdateFriendRequestStatus(ID uint, status string) error {
	if r.DB == nil {
		fmt.Println("DB is nil")
	}
	return r.DB.Model(&model.FriendRequest{}).Where("ID = ?",ID).Update("status", status).Error
}
