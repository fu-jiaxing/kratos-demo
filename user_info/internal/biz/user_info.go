package biz

import (
	"context"
	"fmt"
	"strconv"
	"user_info/internal/data"
)

type UserBizHandler struct {
	repo      data.UserRepo
	redisRepo *data.RedisRepo
}


func NewUserBizHandler(repo data.UserRepo, redisRepo *data.RedisRepo) *UserBizHandler {
	return &UserBizHandler{repo: repo, redisRepo: redisRepo}
}

func (h *UserBizHandler) GetUserInfo(ctx context.Context, userID int64) (user *data.User, err error) {
	user, err = h.repo.FindByID(ctx, userID)
	key := fmt.Sprintf("U:%d:VISITS", userID)
	v, err := h.redisRepo.Get(key)
	if err != nil {
		return user, nil
	}
	vint, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		return user, nil
	}
	user.Visits = vint
	return user, err
}
