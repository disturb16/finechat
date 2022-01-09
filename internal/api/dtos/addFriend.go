package dtos

type AddFriend struct {
	Email       string `json:"email"`
	FriendEmail string `json:"friend_email"`
}
