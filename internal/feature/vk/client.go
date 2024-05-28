package vk

import "context"

type Client interface {
	SendMessage(ctx context.Context, userID int64, message string) error
}

type client struct {
	token string
}

func NewClient(token string) Client {
	return &client{token: token}
}
