package friends

import (
	"akichat/backend/internal/realtime"
	"akichat/backend/internal/repository"
)

type service struct {
	friendShipRepo    *repository.FriendShipRepository
	friendRequestRepo *repository.FriendRequestRepository
	rt                realtime.Gateway
}

func NewService(
	fsRepo *repository.FriendShipRepository,
	frRepo *repository.FriendRequestRepository,
	rt realtime.Gateway,
) Service {
	return &service{
		friendShipRepo:    fsRepo,
		friendRequestRepo: frRepo,
		rt:                rt,
	}
}

func (s *service) ListFriends(userID uint) ([]FriendDTO, error) {
	users, err := s.friendShipRepo.GetFriendsByUserID(userID)
	if err != nil {
		return nil, err
	}
	result := make([]FriendDTO, 0, len(users))
	for _, u := range users {
		result = append(result, FriendDTO{
			ID:       u.ID,
			Username: u.Username,
		})
	}
	return result, nil
}

func (s *service) RequestFriend(fromUserID, toUserID uint) error {
	if err := s.friendRequestRepo.CreateFriendRequest(fromUserID, toUserID); err != nil {
		return err
	}
	// 既存の通知メッセージに合わせる
	if s.rt != nil {
		_ = s.rt.NotifyUser(fromUserID, toUserID, "You received a friend request!")
	}
	return nil
}

// ApproveFriend: 既存ハンドラは RequestUserID と UserID を受け取り AddFriend している。
// フェーズ3では挙動維持のため、第1引数を「申請を出した側のユーザーID」として扱う。
func (s *service) ApproveFriend(requestUserID uint, approverID uint) error {
	return s.friendShipRepo.AddFriend(requestUserID, approverID)
}


