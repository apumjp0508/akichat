package db

import (
    "fmt"
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
	"akichat/backend/internal/model"
)

func InitDB() (*gorm.DB, error) {
    dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
        "meeting",           // ユーザー名
        "meetingpass",       // パスワード
        "mysql",           // ホスト名（Dockerのサービス名）
        3306,             // ポート
        "meetingdb",        // データベース名
    )

    db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
    if err != nil {
        return nil, fmt.Errorf("DB接続失敗: %w", err)
    }else {
		fmt.Println("DB接続成功")
	}

	err = db.AutoMigrate(
        &model.User{},
        &model.FriendRequest{},
        &model.FriendShip{},
    )
    if err != nil {
        return nil, fmt.Errorf("マイグレーション失敗: %w", err)
    }else {
		fmt.Println("マイグレーション成功")
	}

    return db, nil
}
