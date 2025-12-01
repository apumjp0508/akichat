package auth

import (
	"errors"

	"akichat/backend/internal/handler/auth/token/JWTToken"
	"akichat/backend/internal/model"
	"akichat/backend/internal/repository"
)

type service struct {
	users *repository.UserRepository
}

// NewService は認証用サービスを生成します。
func NewService(users *repository.UserRepository) Service {
	return &service{
		users: users,
	}
}

func (s *service) Login(in LoginInput) (LoginOutput, error) {
	user, err := s.users.GetUserByEmail(in.Email)
	if err != nil {
		return LoginOutput{}, err
	}
	// 現状は平文パスワードのため単純比較（将来はハッシュへ移行）
	if user.Password != in.Password {
		return LoginOutput{}, errors.New("invalid credentials")
	}

	access, refresh, err := JWTHandler.GenerateTokens(user.ID, user.Email)
	if err != nil {
		return LoginOutput{}, err
	}
	return LoginOutput{
		AccessToken:  access,
		RefreshToken: refresh,
		UserID:       user.ID,
		UserName:     user.Username,
		UserEmail:    user.Email,
	}, nil
}

func (s *service) Register(in RegisterInput) (RegisterOutput, error) {
	u := &model.User{
		Username: in.Name,
		Email:    in.Email,
		Password: in.Password, // 将来はハッシュ化へ移行
	}
	if err := s.users.CreateUser(u); err != nil {
		return RegisterOutput{}, err
	}
	access, refresh, err := JWTHandler.GenerateTokens(u.ID, u.Email)
	if err != nil {
		return RegisterOutput{}, err
	}
	return RegisterOutput{
		AccessToken:  access,
		RefreshToken: refresh,
		UserID:       u.ID,
		UserName:     u.Username,
		UserEmail:    u.Email,
	}, nil
}


