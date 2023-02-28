// Code generated by protoc-gen-go-http. DO NOT EDIT.
// versions:
// - protoc-gen-go-http v2.5.3
// - protoc             v3.21.6
// source: api/user_info.proto

package user_info

import (
	context "context"
	http "github.com/go-kratos/kratos/v2/transport/http"
	binding "github.com/go-kratos/kratos/v2/transport/http/binding"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the kratos package it is being compiled against.
var _ = new(context.Context)
var _ = binding.EncodeURL

const _ = http.SupportPackageIsVersion1

const OperationUserInfoGetUserInfo = "/api.user_info.UserInfo/GetUserInfo"

type UserInfoHTTPServer interface {
	// GetUserInfo rpc CreateUserInfo (CreateUserInfoRequest) returns (CreateUserInfoReply);
	// rpc UpdateUserInfo (UpdateUserInfoRequest) returns (UpdateUserInfoReply);
	// rpc DeleteUserInfo (DeleteUserInfoRequest) returns (DeleteUserInfoReply);
	GetUserInfo(context.Context, *GetUserInfoRequest) (*GetUserInfoReply, error)
}

func RegisterUserInfoHTTPServer(s *http.Server, srv UserInfoHTTPServer) {
	r := s.Route("/")
	r.GET("/user_info/{userID}", _UserInfo_GetUserInfo0_HTTP_Handler(srv))
}

func _UserInfo_GetUserInfo0_HTTP_Handler(srv UserInfoHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in GetUserInfoRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		if err := ctx.BindVars(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationUserInfoGetUserInfo)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.GetUserInfo(ctx, req.(*GetUserInfoRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*GetUserInfoReply)
		return ctx.Result(200, reply)
	}
}

type UserInfoHTTPClient interface {
	GetUserInfo(ctx context.Context, req *GetUserInfoRequest, opts ...http.CallOption) (rsp *GetUserInfoReply, err error)
}

type UserInfoHTTPClientImpl struct {
	cc *http.Client
}

func NewUserInfoHTTPClient(client *http.Client) UserInfoHTTPClient {
	return &UserInfoHTTPClientImpl{client}
}

func (c *UserInfoHTTPClientImpl) GetUserInfo(ctx context.Context, in *GetUserInfoRequest, opts ...http.CallOption) (*GetUserInfoReply, error) {
	var out GetUserInfoReply
	pattern := "/user_info/{userID}"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationUserInfoGetUserInfo))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}
