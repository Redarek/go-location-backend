package app

import (
	"github.com/rs/zerolog/log"
	"location-backend/internal/config"
	"location-backend/internal/db"
	"location-backend/internal/logger"
	"location-backend/internal/server"
	"location-backend/internal/socket"
	"os"
	"os/signal"
	"syscall"
)

type App struct {
	DB     db.Service
	Fiber  *server.Fiber
	Socket *socket.Server
}

func New() App {
	logger.Setup()
	config.Init()
	postgres := db.New()
	app := App{DB: postgres, Fiber: server.New(postgres), Socket: socket.New(postgres)}
	return app
}

func (app App) Run() {
	var err error
	app.Fiber.RegisterFiberRoutes()
	err = app.Fiber.App.Listen(":" + config.App.Port)
	if err != nil {
		log.Error().Err(err).Msg("Cannot start fiber server")
	}
}

// TODO
func GracefulShutdown(app App) {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	//s := <-sig
	//log.WithField("signal", s).
	//	Debug("received signal")
	//utils.WaitFor(app.Stop, config.App.ShutdownTimeout)
	log.Printf("Shutdown app")
}
