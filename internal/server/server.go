package server

import (
	"log/slog"

	desc "vk-notifications-api/api/vk/notifcations"
	"vk-notifications-api/internal/dao"
	"vk-notifications-api/internal/feature"
)

type server struct {
	clients feature.Clients
	dao     dao.DAO
	log     *slog.Logger

	desc.UnimplementedVkNotificationsApiServer
}

func NewServer(dao dao.DAO, clients feature.Clients, logger *slog.Logger) desc.VkNotificationsApiServer {
	s := &server{
		dao:     dao,
		clients: clients,
		log:     logger,
	}
	return s
}
