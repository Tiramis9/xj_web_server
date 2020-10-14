// Code generated by protoc-gen-micro. DO NOT EDIT.
// source: user.proto

package proto

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

import (
	context "context"
	client "github.com/micro/go-micro/client"
	server "github.com/micro/go-micro/server"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ client.Option
var _ server.Option

// Client API for User service

type UserService interface {
	// 用户注册
	SignUp(ctx context.Context, in *ReqSignUp, opts ...client.CallOption) (*RespSignUp, error)
}

type userService struct {
	c    client.Client
	name string
}

func NewUserService(name string, c client.Client) UserService {
	if c == nil {
		c = client.NewClient()
	}
	if len(name) == 0 {
		name = "proto"
	}
	return &userService{
		c:    c,
		name: name,
	}
}

func (c *userService) SignUp(ctx context.Context, in *ReqSignUp, opts ...client.CallOption) (*RespSignUp, error) {
	req := c.c.NewRequest(c.name, "User.SignUp", in)
	out := new(RespSignUp)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for User service

type UserHandler interface {
	// 用户注册
	SignUp(context.Context, *ReqSignUp, *RespSignUp) error
}

func RegisterUserHandler(s server.Server, hdlr UserHandler, opts ...server.HandlerOption) error {
	type user interface {
		SignUp(ctx context.Context, in *ReqSignUp, out *RespSignUp) error
	}
	type User struct {
		user
	}
	h := &userHandler{hdlr}
	return s.Handle(s.NewHandler(&User{h}, opts...))
}

type userHandler struct {
	UserHandler
}

func (h *userHandler) SignUp(ctx context.Context, in *ReqSignUp, out *RespSignUp) error {
	return h.UserHandler.SignUp(ctx, in, out)
}
