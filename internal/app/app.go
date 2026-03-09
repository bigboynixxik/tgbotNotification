package app

import (
	"TGNotification/internal/transport/telegram"
	"TGNotification/pkg/config"
	"TGNotification/pkg/logger"
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

type App struct {
	TGBot *telegram.Bot
}

func NewApp() *App {
	cfg, err := config.LoadConfig(".env")
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}
	ctx := context.Background()
	logComponent := logger.FromContext(ctx).With("component", "app")
	ctx = logger.IntoContext(ctx, logComponent)

	tgBot, err := telegram.NewBot(cfg.TGToken)
	if err != nil {
		logger.FromContext(ctx).Error("Failed to create telegram bot", "error", err.Error())
		os.Exit(1)
	}

	return &App{TGBot: tgBot}
}

func (a *App) Run() {

	ctx := context.Background()

	logger.FromContext(ctx).Info("Starting application...")

	slog.Info("Starting application...")
	go a.TGBot.Start(ctx)

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit

	slog.Info("Application shutting down gracefully...")
}
