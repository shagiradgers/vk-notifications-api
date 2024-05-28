package server

import (
	"context"
	"fmt"

	"google.golang.org/grpc/codes"
	desc "vk-notifications-api/api/vk/notifcations"
	"vk-notifications-api/internal/dao"
	"vk-notifications-api/internal/errors"
	"vk-notifications-api/internal/utils/nulltypes"
)

func (s *server) GetUsersByFilter(
	ctx context.Context,
	req *desc.GetUsersByFilterRequest,
) (*desc.GetUsersByFilterResponse, error) {
	h, err := newGetUsersByFilterHandler(ctx, s.dao, req)
	if err != nil {
		return nil, err
	}
	err = h.handle()
	return h.response(), err
}

func (h *getUsersByFilterHandler) handle() error {
	if h == nil {
		return fmt.Errorf("nil receiver")
	}

	var err error
	h.users, err = h.dao.NewUserQuery().
		GetUserByFilter(h.ctx, h.userFilter, uint64(h.limit), uint64(h.offset), h.fields...)
	return errors.WrapToNetwork(err).ToGRPCError()
}

func (h *getUsersByFilterHandler) response() *desc.GetUsersByFilterResponse {
	users := make([]*desc.User, 0, len(h.users))
	var patronymic *string

	for idx := range h.users {
		if h.users[idx].Patronymic.Valid {
			patronymic = &h.users[idx].Patronymic.String
		}

		users = append(users, &desc.User{
			UserId:   h.users[idx].Id,
			VkId:     h.users[idx].VkID,
			UserRole: desc.UserRole(desc.UserRole_value[h.users[idx].Role]),
			Group:    h.users[idx].Group,
			Fio: &desc.FIO{
				Firstname:  h.users[idx].Firstname,
				Surname:    h.users[idx].Surname,
				Patronymic: patronymic,
			},
			MobilePhone: h.users[idx].MobilePhone,
			UserStatus:  desc.UserStatus(desc.UserRole_value[h.users[idx].Status]),
		})
	}

	return &desc.GetUsersByFilterResponse{
		Users:  users,
		Limit:  h.limit,
		Offset: h.offset,
		Count:  int64(len(h.users)),
	}
}

type getUsersByFilterHandler struct {
	ctx context.Context
	dao dao.DAO

	userFilter dao.UserTable
	fields     []string
	limit      int64
	offset     int64

	users []dao.UserTable
}

func newGetUsersByFilterHandler(
	ctx context.Context,
	dao dao.DAO,
	req *desc.GetUsersByFilterRequest,
) (*getUsersByFilterHandler, error) {
	h := &getUsersByFilterHandler{
		ctx:    ctx,
		dao:    dao,
		fields: make([]string, 0),
	}

	return h.adapt(req), h.validate()
}

func (h *getUsersByFilterHandler) adapt(
	req *desc.GetUsersByFilterRequest,
) *getUsersByFilterHandler {
	h.limit = req.GetLimit()
	h.offset = req.GetOffset()

	if req.UserStatus != nil {
		h.userFilter.Status = req.GetUserStatus().String()
		h.fields = append(h.fields, "status")
	}
	if req.UserRole != nil {
		h.userFilter.Role = req.GetUserRole().String()
		h.fields = append(h.fields, "role")
	}
	if req.VkId != nil {
		h.userFilter.VkID = req.GetVkId()
		h.fields = append(h.fields, "vk_id")
	}
	if req.Group != nil {
		h.userFilter.Group = req.GetGroup()
		h.fields = append(h.fields, "user_group")
	}
	if req.Firstname != nil {
		h.userFilter.Firstname = req.GetFirstname()
		h.fields = append(h.fields, "firstname")
	}
	if req.Surname != nil {
		h.userFilter.Surname = req.GetSurname()
		h.fields = append(h.fields, "surname")
	}
	if req.Patronymic != nil {
		h.userFilter.Patronymic = nulltypes.NewNullString(req.Patronymic)
		h.fields = append(h.fields, "patronymic")
	}
	if req.MobilePhone != nil {
		h.userFilter.MobilePhone = req.GetMobilePhone()
		h.fields = append(h.fields, "mobile_phone")
	}
	return h
}

func (h *getUsersByFilterHandler) validate() error {
	if h.userFilter.VkID != 0 && h.userFilter.VkID < 0 {
		return errors.
			NewNetworkError(codes.InvalidArgument, "vk_id must be specified").
			ToGRPCError()
	}
	if h.limit <= 0 {
		return errors.
			NewNetworkError(codes.InvalidArgument, "limit must be specified").
			ToGRPCError()
	}
	if h.offset < 0 {
		return errors.
			NewNetworkError(codes.InvalidArgument, "offset must be specified").
			ToGRPCError()
	}
	return nil

}
