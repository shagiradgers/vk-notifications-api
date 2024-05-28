package server

import (
	"context"
	"fmt"

	"google.golang.org/grpc/codes"
	desc "vk-notifications-api/api/vk/notifcations"
	"vk-notifications-api/internal/dao"
	"vk-notifications-api/internal/errors"
)

func (s *server) GetUser(ctx context.Context, req *desc.GetUserRequest) (*desc.GetUserResponse, error) {
	h, err := newGetUserHandler(ctx, s.dao, req)
	if err != nil {
		return nil, err
	}
	err = h.handle()
	return h.response(), err
}

func (h *getUserHandler) response() *desc.GetUserResponse {
	var patronymic *string
	if h.user.Patronymic.Valid {
		patronymic = &h.user.Patronymic.String
	}

	return &desc.GetUserResponse{
		User: &desc.User{
			UserId:   h.user.Id,
			VkId:     h.user.VkID,
			UserRole: desc.UserRole(desc.UserRole_value[h.user.Role]),
			Group:    h.user.Group,
			Fio: &desc.FIO{
				Firstname:  h.user.Firstname,
				Surname:    h.user.Surname,
				Patronymic: patronymic,
			},
			MobilePhone: h.user.MobilePhone,
			UserStatus:  desc.UserStatus(desc.UserStatus_value[h.user.Status]),
		},
	}
}

func (h *getUserHandler) handle() error {
	if h == nil {
		return fmt.Errorf("nil receiver")
	}

	var err error
	h.user, err = h.dao.NewUserQuery().GetUser(h.ctx, h.userID)
	return errors.WrapToNetwork(err).ToGRPCError()
}

type getUserHandler struct {
	ctx context.Context
	dao dao.DAO

	userID int64

	user dao.UserTable
}

func newGetUserHandler(
	ctx context.Context,
	dao dao.DAO,
	req *desc.GetUserRequest,
) (*getUserHandler, error) {
	h := &getUserHandler{
		ctx:    ctx,
		dao:    dao,
		userID: req.GetUserId(),
	}

	return h, h.validate()
}

func (h *getUserHandler) validate() error {
	if h.userID <= 0 {
		return errors.
			NewNetworkError(codes.InvalidArgument, "user_id must be specified").
			ToGRPCError()
	}
	return nil
}
