package service

import (
	"context"
	"user_info/internal/biz"

	pb "user_info/api/user_info"
)

type UserInfoService struct {
	pb.UnimplementedUserInfoServer

	bizHandler *biz.UserBizHandler
}

func NewUserInfoService(h *biz.UserBizHandler) *UserInfoService {
	return &UserInfoService{
		bizHandler: h,
	}
}

func (s *UserInfoService) GetUserInfo(ctx context.Context, req *pb.GetUserInfoRequest) (*pb.GetUserInfoReply, error) {
	user, err := s.bizHandler.GetUserInfo(ctx, req.UserID)
	if err != nil {
		return nil, err
	}
	pbUser := pb.User{
		UserID: req.UserID,
		Name:   user.Name,
		Email:  user.Email,
		Visits: user.Visits,
	}
	resp := pb.GetUserInfoReply{User: &pbUser}
	return &resp, nil
}
