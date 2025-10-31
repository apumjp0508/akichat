package friendsHandler

import (
	websocket "akichat/backend/internal/handler/webSocket"
)

func (h *FriendRequestHandler) NotifyFriendRequest(userID uint,FriendID uint) error {
	message := "You received a friend request!"
	return websocket.GlobalHub.NotifyUser(userID ,FriendID, message)
}
