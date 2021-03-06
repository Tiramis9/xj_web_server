// Code generated by protoc-gen-micro. DO NOT EDIT.
// source: public.proto

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

// Client API for Public service

type PublicService interface {
	// 获取服务器列表
	GetHost(ctx context.Context, in *ReqHost, opts ...client.CallOption) (*RespHost, error)
}

type publicService struct {
	c    client.Client
	name string
}

func NewPublicService(name string, c client.Client) PublicService {
	if c == nil {
		c = client.NewClient()
	}
	if len(name) == 0 {
		name = "proto"
	}
	return &publicService{
		c:    c,
		name: name,
	}
}

func (c *publicService) GetHost(ctx context.Context, in *ReqHost, opts ...client.CallOption) (*RespHost, error) {
	req := c.c.NewRequest(c.name, "Public.GetHost", in)
	out := new(RespHost)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Public service

type PublicHandler interface {
	// 获取服务器列表
	GetHost(context.Context, *ReqHost, *RespHost) error
}

func RegisterPublicHandler(s server.Server, hdlr PublicHandler, opts ...server.HandlerOption) error {
	type public interface {
		GetHost(ctx context.Context, in *ReqHost, out *RespHost) error
	}
	type Public struct {
		public
	}
	h := &publicHandler{hdlr}
	return s.Handle(s.NewHandler(&Public{h}, opts...))
}

type publicHandler struct {
	PublicHandler
}

func (h *publicHandler) GetHost(ctx context.Context, in *ReqHost, out *RespHost) error {
	return h.PublicHandler.GetHost(ctx, in, out)
}
