package biz

import "context"

type UserBizHandler struct {
}

type User struct {
	UserID int64
}

func NewUserBizHandler() *UserBizHandler {
	return &UserBizHandler{}
}

func (h *UserBizHandler) GetUserInfo(ctx context.Context, userID int64) (user *User, err error) {
	user = &User{UserID: 100024}
	return user, nil
}
