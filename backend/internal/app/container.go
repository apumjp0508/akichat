package app

import (
	"akichat/backend/internal/db"
	ws "akichat/backend/internal/handler/webSocket"
	"akichat/backend/internal/realtime"
	"akichat/backend/internal/repository"
	authsvc "akichat/backend/internal/service/auth"
	friendssvc "akichat/backend/internal/service/friends"
	signaling "akichat/backend/internal/service/signaling"
)

type Container struct {
	// Infra
	DB  interface{}
	Hub *ws.Hub
	RT  realtime.Gateway

	// Repository
	UserRepo          *repository.UserRepository
	FriendShipRepo    *repository.FriendShipRepository
	FriendRequestRepo *repository.FriendRequestRepository

	// Services
	AuthService      authsvc.Service
	FriendsService   friendssvc.Service
	SignalingService *signaling.Service
}

func NewContainer() (*Container, error) {
	// 1. DB
	database, err := db.InitDB()
	if err != nil {
		return nil, err
	}

	// 2. Hub（既存の GlobalHub を利用）
	hub := ws.GlobalHub

	// 3. Repositories
	userRepo := repository.NewUserRepository(database)
	friendShipRepo := repository.NewFriendShipRepository(database)
	friendRequestRepo := repository.NewFriendRequestRepository(database)

	// 4. Services
	authService := authsvc.NewService(userRepo)
	friendsService := friendssvc.NewService(friendShipRepo, friendRequestRepo, hub)
	sigService := &signaling.Service{RT: hub}

	return &Container{
		DB:                database,
		Hub:               hub,
		RT:                hub,
		UserRepo:          userRepo,
		FriendShipRepo:    friendShipRepo,
		FriendRequestRepo: friendRequestRepo,
		AuthService:       authService,
		FriendsService:    friendsService,
		SignalingService:  sigService,
	}, nil
}



