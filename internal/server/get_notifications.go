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

func (s *server) GetNotifications(
	ctx context.Context,
	req *desc.GetNotificationsRequest,
) (*desc.GetNotificationsResponse, error) {
	h, err := newGetNotificationsHandler(ctx, s.dao, req)
	if err != nil {
		return nil, err
	}
	err = h.handle()
	return h.response(), err
}

func (h *getNotificationsHandler) handle() error {
	if h == nil {
		return fmt.Errorf("nil receiver")
	}

	var err error
	h.notifications, err = h.dao.NewNotificationQuery().
		GetNotificationsByIDs(h.ctx, h.notificationIDs, h.limit, h.offset)
	return err
}

func (h *getNotificationsHandler) response() *desc.GetNotificationsResponse {
	var mediaContent *string
	notifications := make([]*desc.Notification, 0, len(h.notifications))
	for idx := range h.notifications {
		if h.notifications[idx].MediaContent.Valid {
			mediaContent = &h.notifications[idx].MediaContent.String
		}

		notifications = append(notifications, &desc.Notification{
			NotificationId:     h.notifications[idx].ID,
			SenderId:           h.notifications[idx].SenderID,
			ReceiverIds:        h.notifications[idx].ReceiverIDs,
			Message:            h.notifications[idx].Message,
			MediaContent:       mediaContent,
			NotificationStatus: desc.NotificationStatus(desc.NotificationStatus_value[h.notifications[idx].Status]),
			Date:               timestamppb.New(h.notifications[idx].Date),
		})
	}

	return &desc.GetNotificationsResponse{
		Notification: notifications,
		Limit:        h.limit,
		Offset:       h.offset,
		Count:        int64(len(h.notifications)),
	}
}

type getNotificationsHandler struct {
	ctx context.Context
	dao dao.DAO

	notificationIDs []int64
	limit           int64
	offset          int64

	notifications []dao.NotificationTable
}

func newGetNotificationsHandler(
	ctx context.Context,
	dao dao.DAO,
	req *desc.GetNotificationsRequest,
) (*getNotificationsHandler, error) {
	h := &getNotificationsHandler{
		ctx: ctx,
		dao: dao,
	}
	return h.adapt(req), h.validate()
}

func (h *getNotificationsHandler) adapt(
	req *desc.GetNotificationsRequest,
) *getNotificationsHandler {
	h.limit = req.GetLimit()
	h.offset = req.GetOffset()
	h.notificationIDs = req.GetNotificationIds()
	return h
}

func (h *getNotificationsHandler) validate() error {
	if h.limit <= 0 {
		return errors.
			NewNetworkError(codes.InvalidArgument, "limit must be specified").
			ToGRPCError()
	}
	if h.limit > 500 {
		return errors.
			NewNetworkError(codes.InvalidArgument, "limit greater that 500").
			ToGRPCError()
	}
	if h.offset < 0 {
		return errors.
			NewNetworkError(codes.InvalidArgument, "offset must be specified").
			ToGRPCError()
	}
	if len(h.notificationIDs) == 0 {
		return errors.
			NewNetworkError(codes.InvalidArgument, "notification_ids must be specified").
			ToGRPCError()
	}
	return nil
}
