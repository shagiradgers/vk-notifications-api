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

func (s *server) EditUser(ctx context.Context, req *desc.EditUserRequest) (*desc.EditUserResponse, error) {
	h, err := newEditUserHandler(ctx, s.dao, req)
	if err != nil {
		return nil, err
	}
	err = h.handle()
	return h.response(), err
}

func (h *editUserHandler) handle() error {
	if h == nil {
		return fmt.Errorf("nil receiver")
	}

	var err error
	h.user, err = h.dao.NewUserQuery().UpdateUser(h.ctx, h.userToEdit, h.fields...)
	return errors.WrapToNetwork(err).ToGRPCError()
}

func (h *editUserHandler) response() *desc.EditUserResponse {
	var patronymic *string
	if h.user.Patronymic.Valid {
		patronymic = &h.user.Patronymic.String
	}

	return &desc.EditUserResponse{
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
			UserStatus:  desc.UserStatus(desc.UserRole_value[h.user.Status]),
		},
	}
}

type editUserHandler struct {
	ctx context.Context
	dao dao.DAO

	userToEdit dao.UserTable
	fields     []string

	user dao.UserTable
}

func newEditUserHandler(
	ctx context.Context,
	dao dao.DAO,
	req *desc.EditUserRequest,
) (*editUserHandler, error) {
	h := editUserHandler{
		ctx:    ctx,
		dao:    dao,
		fields: make([]string, 0),
	}
	return h.adapt(req), h.validate()
}

func (h *editUserHandler) adapt(req *desc.EditUserRequest) *editUserHandler {
	h.userToEdit.Id = req.GetUserId()

	if req.UserStatus != nil {
		h.userToEdit.Status = req.GetUserStatus().String()
		h.fields = append(h.fields, "status")
	}
	if req.UserRole != nil {
		h.userToEdit.Role = req.GetUserRole().String()
		h.fields = append(h.fields, "role")
	}
	if req.VkId != nil {
		h.userToEdit.VkID = req.GetVkId()
		h.fields = append(h.fields, "vk_id")
	}
	if req.Group != nil {
		h.userToEdit.Group = req.GetGroup()
		h.fields = append(h.fields, "user_group")
	}
	if req.Firstname != nil {
		h.userToEdit.Firstname = req.GetFirstname()
		h.fields = append(h.fields, "firstname")
	}
	if req.Surname != nil {
		h.userToEdit.Surname = req.GetSurname()
		h.fields = append(h.fields, "surname")
	}
	if req.Patronymic != nil {
		h.userToEdit.Patronymic = nulltypes.NewNullString(req.Patronymic)
		h.fields = append(h.fields, "patronymic")
	}
	if req.MobilePhone != nil {
		h.userToEdit.MobilePhone = req.GetMobilePhone()
		h.fields = append(h.fields, "mobile_phone")
	}
	return h
}

func (h *editUserHandler) validate() error {
	if h.userToEdit.Id <= 0 {
		return errors.
			NewNetworkError(codes.InvalidArgument, "id must be specified").
			ToGRPCError()
	}
	if h.userToEdit.VkID != 0 && h.userToEdit.VkID < 0 {
		return errors.
			NewNetworkError(codes.InvalidArgument, "vk_id must be specified").
			ToGRPCError()
	}
	return nil
}
