package vk

import (
	"context"
	"errors"
	"fmt"
	"net/http"
)

func (c *client) SendMessage(ctx context.Context, userID int64, message string) error {
	const op = "client.SendMessage"
	var (
		err error
		r   *http.Response
	)
	if err = c.validateSendMessage(userID, message); err != nil {
		return fmt.Errorf("ERROR: %s: %s", op, err)
	}
	url := c.getUrl("messages.send", c.getDefaultParams(params{"user_id": userID, "message": message}))

	r, err = c.makeRequest(ctx, http.MethodPost, url, nil)
	if err != nil {
		return fmt.Errorf("ERROR: %s: %s", op, err)
	}

	if r == nil {
		return fmt.Errorf("ERROR: %s: response is nil", op)
	}
	if r.StatusCode != http.StatusOK {
		return fmt.Errorf("ERROR: %s: code not eq 200", op)
	}
	return nil
}

func (c *client) validateSendMessage(userID int64, message string) error {
	if userID <= 0 {
		return errors.New("userID not be less or eq zero")
	}
	if len(message) > 9000 {
		return errors.New("message len not be greater that 9000")
	}
	return nil
}
