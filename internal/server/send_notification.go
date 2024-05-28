package server

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	"google.golang.org/grpc/codes"
	desc "vk-notifications-api/api/vk/notifcations"
	"vk-notifications-api/internal/dao"
	"vk-notifications-api/internal/errors"
	"vk-notifications-api/internal/feature"
	"vk-notifications-api/internal/utils/nulltypes"
)

const enableDrainSendMessage = true

func (s *server) SendNotification(
	ctx context.Context,
	req *desc.SendNotificationRequest,
) (*desc.SendNotificationResponse, error) {
	h, err := newSendNotificationHandler(ctx, s.dao, s.clients, s.log, time.Now(), req)
	if err != nil {
		return nil, err
	}
	err = h.handle()
	return h.response(), err
}

func (h *sendNotificationHandler) handle() error {
	if h == nil {
		return fmt.Errorf("nil receiver")
	}
	var err error

	h.notification, err = h.dao.
		NewNotificationQuery().
		AddNotification(h.ctx, h.senderId, h.receiverIds, h.message, h.mediaContent, h.now)
	if err != nil {
		return errors.WrapToNetwork(err).ToGRPCError()
	}

	var receiver dao.UserTable
	for idx := range h.notification.ReceiverIDs {
		receiver, err = h.dao.NewUserQuery().GetUser(h.ctx, h.notification.ReceiverIDs[idx])
		if err != nil {
			return errors.WrapToNetwork(err).ToGRPCError()
		}

		err = h.clients.VkClient().
			SendMessage(h.ctx, receiver.VkID, h.notification.Message)

		// if u want send messages with ignoring errors set enableDrainSendMessage = true
		if err != nil && !enableDrainSendMessage {
			h.notification.Status = desc.NotificationStatus_PROBLEM.String()
			err = h.dao.NewNotificationQuery().
				UpdateNotificationStatus(h.ctx, h.notification.ID, desc.NotificationStatus_PROBLEM.String())
			return errors.WrapToNetwork(err).ToGRPCError()
		}

		if err != nil && enableDrainSendMessage {
			h.log.Error("SendNotification: message not send, but enableDrainSendMessage = true: " + err.Error())
		}
	}

	h.notification.Status = desc.NotificationStatus_SEND.String()
	err = h.dao.NewNotificationQuery().
		UpdateNotificationStatus(h.ctx, h.notification.ID, desc.NotificationStatus_SEND.String())
	return errors.WrapToNetwork(err).ToGRPCError()
}

func (h *sendNotificationHandler) response() *desc.SendNotificationResponse {
	return &desc.SendNotificationResponse{
		NotificationId: h.notification.ID,
		MessageStatus:  desc.NotificationStatus(desc.NotificationStatus_value[h.notification.Status]),
	}
}

type sendNotificationHandler struct {
	ctx     context.Context
	dao     dao.DAO
	clients feature.Clients
	log     *slog.Logger

	senderId     int64
	receiverIds  []int64
	message      string
	mediaContent sql.NullString
	now          time.Time

	notification dao.NotificationTable
}

func newSendNotificationHandler(
	ctx context.Context,
	dao dao.DAO,
	clients feature.Clients,
	log *slog.Logger,
	now time.Time,
	req *desc.SendNotificationRequest,
) (*sendNotificationHandler, error) {
	h := &sendNotificationHandler{
		ctx:     ctx,
		dao:     dao,
		clients: clients,
		log:     log,
		now:     now,
	}
	return h.adapt(req), h.validate()
}

func (h *sendNotificationHandler) adapt(
	req *desc.SendNotificationRequest,
) *sendNotificationHandler {
	h.mediaContent = nulltypes.NewNullString(req.MediaContent)
	h.message = req.GetMessage()
	h.senderId = req.GetSenderId()
	h.receiverIds = req.GetReceiverIds()
	return h
}

func (h *sendNotificationHandler) validate() error {
	if len(h.message) == 0 {
		return errors.
			NewNetworkError(codes.InvalidArgument, "message must be specified").
			ToGRPCError()
	}
	if len(h.receiverIds) == 0 {
		return errors.
			NewNetworkError(codes.InvalidArgument, "receiverIds must be specified").
			ToGRPCError()
	}
	if h.senderId <= 0 {
		return errors.
			NewNetworkError(codes.InvalidArgument, "senderId must be specified").
			ToGRPCError()
	}
	return nil
}
