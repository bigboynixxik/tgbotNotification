package app

import (
	"TGNotification/internal/clients"
	"TGNotification/internal/repository"
	"TGNotification/internal/service"
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
	TGBot   *telegram.Bot
	Service *service.NotifierService
}

func NewApp() *App {
	cfg, err := config.LoadConfig(".env")
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}
	ctx := context.Background()
	logComponent := logger.FromContext(ctx).With("component", "app")
	ctx = logger.IntoContext(ctx, logComponent)

	djangoClient, err := clients.NewDjangoClient(cfg.DjangoGRPCAddr)
	if err != nil {
		logger.FromContext(ctx).Error("app.NewApp, failed to connect to django grpc",
			slog.String("error", err.Error()))
		os.Exit(1)
	}

	repo, err := repository.NewRedisRepository(ctx, cfg.RedisAddr)
	if err != nil {
		logger.FromContext(ctx).Error("Failed to connect to Redis", slog.String("error", err.Error()))
		os.Exit(1)
	}

	tgBot, err := telegram.NewBot(cfg.TGToken, djangoClient)
	if err != nil {
		logger.FromContext(ctx).Error("Failed to create telegram bot", "error", err.Error())
		os.Exit(1)
	}

	notifierSvc := service.NewNotifierService(tgBot, repo)

	return &App{
		TGBot:   tgBot,
		Service: notifierSvc,
	}
}

func (a *App) Run() {

	ctx := context.Background()

	logs := logger.FromContext(ctx)
	logs.Info("Starting application...")

	go a.TGBot.Start(ctx)

	go func() {
		if err := a.Service.StartWorker(ctx); err != nil {
			logs.Error("app.Run worker stop with error",
				slog.String("error", err.Error()))
		}
	}()

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit

	logs.Info("Application shutting down gracefully...")
}
