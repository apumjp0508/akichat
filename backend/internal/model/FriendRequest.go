package model

import(
	"time"
)

type FriendRequest struct {
	ID       uint   `gorm:"primaryKey"`
	RequestUserID   uint   `json:"request_user_id"`
	ReceiverUserID uint   `json:"Receiver_friend_id"`

	Status string `json:"status" gorm:"default:'pending'"` // "pending" / "accepted" / "rejected"
	CreatedAt time.Time
	UpdatedAt time.Time

	// リレーション
	RequestUser  User `gorm:"foreignKey:RequestUserID;constraint:OnDelete:CASCADE"`
	ReceiverUser User `gorm:"foreignKey:ReceiverUserID;constraint:OnDelete:CASCADE"`
}