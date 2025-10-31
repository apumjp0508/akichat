package model

import(
	"time"
)

type FriendShip struct {
	ID       uint   `gorm:"primaryKey"`
	User1ID   uint   `json:"user_1_id"`
	User2ID uint   `json:"user_2_id"`

	User1 User `gorm:"foreignKey:User1ID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	User2 User `gorm:"foreignKey:User2ID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	
	Status    string `json:"status" gorm:"default:'active'"` // active / blocked / etc.
	CreatedAt time.Time
}