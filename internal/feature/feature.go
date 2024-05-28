package feature

import (
	"vk-notifications-api/internal/config"
	"vk-notifications-api/internal/feature/vk"
)

type Clients interface {
	VkClient() vk.Client
}

type clients struct {
	vkClient vk.Client
}

func (c *clients) VkClient() vk.Client {
	return c.vkClient
}

func NewClients(cfg config.Config) Clients {
	return &clients{
		vkClient: vk.NewClient(cfg.MustGetVkToken()),
	}
}
