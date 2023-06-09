// Code generated by protoc-gen-micro. DO NOT EDIT.
// source: proto/user.proto

package user

import (
	fmt "fmt"
	proto "google.golang.org/protobuf/proto"
	math "math"
)

import (
	context "context"
	api "go-micro.dev/v4/api"
	client "go-micro.dev/v4/client"
	server "go-micro.dev/v4/server"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// Reference imports to suppress errors if they are not otherwise used.
var _ api.Endpoint
var _ context.Context
var _ client.Option
var _ server.Option

// Api Endpoints for User service

func NewUserEndpoints() []*api.Endpoint {
	return []*api.Endpoint{}
}

// Client API for User service

type UserService interface {
	SendSms(ctx context.Context, in *Request, opts ...client.CallOption) (*Response, error)
	Register(ctx context.Context, in *RegReq, opts ...client.CallOption) (*Response, error)
	UploadAvatar(ctx context.Context, in *UploadReq, opts ...client.CallOption) (*Response, error)
	AuthUpdate(ctx context.Context, in *AuthReq, opts ...client.CallOption) (*Response, error)
}

type userService struct {
	c    client.Client
	name string
}

func NewUserService(name string, c client.Client) UserService {
	return &userService{
		c:    c,
		name: name,
	}
}

func (c *userService) SendSms(ctx context.Context, in *Request, opts ...client.CallOption) (*Response, error) {
	req := c.c.NewRequest(c.name, "User.SendSms", in)
	out := new(Response)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userService) Register(ctx context.Context, in *RegReq, opts ...client.CallOption) (*Response, error) {
	req := c.c.NewRequest(c.name, "User.Register", in)
	out := new(Response)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userService) UploadAvatar(ctx context.Context, in *UploadReq, opts ...client.CallOption) (*Response, error) {
	req := c.c.NewRequest(c.name, "User.UploadAvatar", in)
	out := new(Response)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userService) AuthUpdate(ctx context.Context, in *AuthReq, opts ...client.CallOption) (*Response, error) {
	req := c.c.NewRequest(c.name, "User.AuthUpdate", in)
	out := new(Response)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for User service

type UserHandler interface {
	SendSms(context.Context, *Request, *Response) error
	Register(context.Context, *RegReq, *Response) error
	UploadAvatar(context.Context, *UploadReq, *Response) error
	AuthUpdate(context.Context, *AuthReq, *Response) error
}

func RegisterUserHandler(s server.Server, hdlr UserHandler, opts ...server.HandlerOption) error {
	type user interface {
		SendSms(ctx context.Context, in *Request, out *Response) error
		Register(ctx context.Context, in *RegReq, out *Response) error
		UploadAvatar(ctx context.Context, in *UploadReq, out *Response) error
		AuthUpdate(ctx context.Context, in *AuthReq, out *Response) error
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

func (h *userHandler) SendSms(ctx context.Context, in *Request, out *Response) error {
	return h.UserHandler.SendSms(ctx, in, out)
}

func (h *userHandler) Register(ctx context.Context, in *RegReq, out *Response) error {
	return h.UserHandler.Register(ctx, in, out)
}

func (h *userHandler) UploadAvatar(ctx context.Context, in *UploadReq, out *Response) error {
	return h.UserHandler.UploadAvatar(ctx, in, out)
}

func (h *userHandler) AuthUpdate(ctx context.Context, in *AuthReq, out *Response) error {
	return h.UserHandler.AuthUpdate(ctx, in, out)
}
