package app

import (
	"context"
	"main/internal/service"
	"os"
	"os/signal"
	"syscall"

	"github.com/yanun0323/pkg/logs"
)

func Run() {
	ctx := context.Background()
	l := logs.New("bot", 0)

	bot, err := service.NewBot("General", "SOL/USDT")
	if err != nil {
		l.WithError(err).Fatal("create bot")
	}

	if err := bot.Run(ctx); err != nil {
		l.WithError(err).Fatal("run bot")
	}

	/* Graceful shutdown */
	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)
	<-sigterm
	l.Info("shutdown server")
	if err := bot.Shutdown(ctx); err != nil {
		l.WithError(err).Error("shutdown server")
	}
}
