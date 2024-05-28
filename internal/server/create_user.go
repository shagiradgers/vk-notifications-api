package server

import (
	"context"
	"database/sql"
	"fmt"

	"google.golang.org/grpc/codes"
	desc "vk-notifications-api/api/vk/notifcations"
	"vk-notifications-api/internal/dao"
	"vk-notifications-api/internal/errors"
	"vk-notifications-api/internal/utils/nulltypes"
)

func (s *server) CreateUser(ctx context.Context, req *desc.CreateUserRequest) (*desc.CreateUserResponse, error) {
	h, err := newCreateUserHandler(ctx, s.dao, req)
	if err != nil {
		return nil, err
	}
	err = h.handle()
	return h.response(), err
}

func (h *createUserHandler) handle() error {
	if h == nil {
		return fmt.Errorf("nil receiver")
	}

	user, err := h.dao.NewUserQuery().
		CreateUser(
			h.ctx,
			h.vkId,
			h.userRole,
			h.group,
			h.firstname,
			h.surname,
			h.patronymic,
			h.mobilePhone,
			desc.UserStatus_ACTIVE.String(),
		)
	h.user = user
	return errors.WrapToNetwork(err).ToGRPCError()
}

func (h *createUserHandler) response() *desc.CreateUserResponse {
	var patronymic *string
	if h.patronymic.Valid {
		patronymic = &h.patronymic.String
	}

	return &desc.CreateUserResponse{
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

type createUserHandler struct {
	ctx context.Context
	dao dao.DAO

	vkId        int64
	userRole    string
	group       string
	mobilePhone string
	firstname   string
	surname     string
	patronymic  sql.NullString

	user dao.UserTable
}

func newCreateUserHandler(ctx context.Context, dao dao.DAO, req *desc.CreateUserRequest) (*createUserHandler, error) {
	h := createUserHandler{
		ctx: ctx,
		dao: dao,
	}
	return h.adapt(req), h.validate()
}

func (h *createUserHandler) adapt(req *desc.CreateUserRequest) *createUserHandler {
	h.vkId = req.GetVkId()
	h.userRole = req.GetUserRole().String()
	h.group = req.GetGroup()
	h.mobilePhone = req.GetMobilePhone()
	h.firstname = req.GetFio().GetFirstname()
	h.surname = req.GetFio().GetSurname()
	h.patronymic = nulltypes.NewNullString(req.GetFio().Patronymic)
	return h
}

func (h *createUserHandler) validate() error {
	if h.vkId <= 0 {
		return errors.
			NewNetworkError(codes.InvalidArgument, "vk_id must be specified").
			ToGRPCError()
	}
	if h.group == "" {
		return errors.
			NewNetworkError(codes.InvalidArgument, "group must be specified").
			ToGRPCError()
	}
	if h.mobilePhone == "" {
		return errors.
			NewNetworkError(codes.InvalidArgument, "mobile_phone must be specified").
			ToGRPCError()
	}
	if h.firstname == "" {
		return errors.
			NewNetworkError(codes.InvalidArgument, "firstname must be specified").
			ToGRPCError()
	}
	if h.surname == "" {
		return errors.
			NewNetworkError(codes.InvalidArgument, "surname must be specified").
			ToGRPCError()
	}
	return nil
}
