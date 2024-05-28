package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"vk-notifications-api/internal/app"
	cfg "vk-notifications-api/internal/config"
	"vk-notifications-api/internal/dao"
	"vk-notifications-api/internal/feature"
	"vk-notifications-api/internal/storage"
)

var logger *slog.Logger

func init() {
	logger = slog.New(slog.NewTextHandler(os.Stdout, nil))
}

func main() {
	config, err := cfg.NewConfig(cfg.LocalEnv)
	if err != nil {
		logger.Error("can't create config: ", err)
		return
	}

	s, err := storage.NewStorage(config.MustGetDatabaseConnectionString())
	if err != nil {
		logger.Error("can't create storage: ", err)
		return
	}

	c := feature.NewClients(config)
	d := dao.NewDAO(s)
	a := app.New(logger, d, c, config.MustGetServerHost(), config.MustGetServerPort())
	go func() {
		if err = a.Run(); err != nil {
			logger.Error("can't run app: ", err)
			return
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	<-stop
	a.Stop()
	logger.Info("Gracefully stopped")
}
