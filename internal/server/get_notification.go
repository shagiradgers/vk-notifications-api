package server

import (
	"context"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/types/known/timestamppb"
	desc "vk-notifications-api/api/vk/notifcations"
	"vk-notifications-api/internal/dao"
	"vk-notifications-api/internal/errors"
)

func (s *server) GetNotification(ctx context.Context, req *desc.GetNotificationRequest) (*desc.GetNotificationResponse, error) {
	h, err := newGetNotificationHandler(ctx, s.dao, req)
	if err != nil {
		return nil, err
	}
	err = h.handle()
	return h.response(), err
}

func (h *getNotificationHandler) handle() error {
	if h == nil {
		return fmt.Errorf("nil receiver")
	}

	var err error
	h.notification, err = h.dao.NewNotificationQuery().GetNotification(h.ctx, h.notificationID)
	return errors.WrapToNetwork(err).ToGRPCError()
}

func (h *getNotificationHandler) response() *desc.GetNotificationResponse {
	var mediaContent *string
	if h.notification.MediaContent.Valid {
		mediaContent = &h.notification.MediaContent.String
	}

	return &desc.GetNotificationResponse{
		Notification: &desc.Notification{
			NotificationId:     h.notification.ID,
			SenderId:           h.notification.SenderID,
			ReceiverIds:        h.notification.ReceiverIDs,
			Message:            h.notification.Message,
			MediaContent:       mediaContent,
			NotificationStatus: desc.NotificationStatus(desc.NotificationStatus_value[h.notification.Status]),
			Date:               timestamppb.New(h.notification.Date),
		},
	}
}

type getNotificationHandler struct {
	ctx context.Context
	dao dao.DAO

	notificationID int64

	notification dao.NotificationTable
}

func newGetNotificationHandler(
	ctx context.Context,
	dao dao.DAO,
	req *desc.GetNotificationRequest,
) (*getNotificationHandler, error) {
	h := &getNotificationHandler{
		ctx:            ctx,
		dao:            dao,
		notificationID: req.GetNotificationId(),
	}

	return h, h.validate()
}

func (h *getNotificationHandler) validate() error {
	if h.notificationID <= 0 {
		return errors.
			NewNetworkError(codes.InvalidArgument, "notification_id must be specified").
			ToGRPCError()
	}
	return nil
}
