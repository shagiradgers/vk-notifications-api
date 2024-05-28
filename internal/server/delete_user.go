package server

import (
	"context"
	"fmt"

	"google.golang.org/grpc/codes"
	desc "vk-notifications-api/api/vk/notifcations"
	"vk-notifications-api/internal/dao"
	"vk-notifications-api/internal/errors"
)

func (s *server) DeleteUser(ctx context.Context, req *desc.DeleteUserRequest) (*desc.DeleteUserResponse, error) {
	h, err := newDeleteUserHandler(ctx, s.dao, req)
	if err != nil {
		return nil, err
	}
	err = h.handle()
	return h.response(), err
}

func (h *deleteUserHandler) handle() error {
	if h == nil {
		return fmt.Errorf("nil receiver")
	}

	err := h.dao.NewUserQuery().DeleteUser(h.ctx, h.userID)
	return errors.WrapToNetwork(err).ToGRPCError()
}

func (h *deleteUserHandler) response() *desc.DeleteUserResponse {
	return &desc.DeleteUserResponse{}
}

type deleteUserHandler struct {
	ctx context.Context
	dao dao.DAO

	userID int64
}

func newDeleteUserHandler(
	ctx context.Context,
	dao dao.DAO,
	req *desc.DeleteUserRequest,
) (*deleteUserHandler, error) {
	h := &deleteUserHandler{
		ctx:    ctx,
		dao:    dao,
		userID: req.GetUserId(),
	}
	return h, h.validate()
}

func (h *deleteUserHandler) validate() error {
	if h.userID <= 0 {
		return errors.
			NewNetworkError(codes.InvalidArgument, "user_id must be specified").
			ToGRPCError()
	}
	return nil
}
