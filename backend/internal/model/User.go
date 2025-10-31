package model

type User struct {
    ID       uint   `gorm:"primaryKey"`
    Username string `json:"username"`
    Password string `json:"password"`
    Email    string `json:"email"`
}