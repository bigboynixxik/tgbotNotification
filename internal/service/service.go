package service

import (
	"TGNotification/internal/models"
	"TGNotification/pkg/logger"
	"context"
	"encoding/json"
	"log/slog"
	"time"
)

const QueueName = "django_notification_queue"

type NotifierService struct {
	Bot   BotProvider
	Queue QueueProvider
}

func NewNotifierService(bot BotProvider, queueProvider QueueProvider) *NotifierService {
	return &NotifierService{
		Bot:   bot,
		Queue: queueProvider,
	}
}

func (ns *NotifierService) StartWorker(ctx context.Context) error {
	log := logger.FromContext(ctx)
	log.Info("service.StartWorker starting notifier worker")

	for {
		select {
		case <-ctx.Done():
			log.Info("service.StartWorker stopping notifier worker")
			return nil
		default:
		}
		jsonStr, err := ns.Queue.ListenQueue(ctx, QueueName)
		if err != nil {
			log.Error("service.StartWorker error listening queue",
				slog.String("queue", QueueName),
				slog.String("error", err.Error()))
			time.Sleep(1 * time.Second)
			continue
		}
		var notification models.Notification
		if err := json.Unmarshal([]byte(jsonStr), &notification); err != nil {
			log.Error("service.StartWorker error unmarshalling json notification",
				slog.String("json", jsonStr),
				slog.String("error", err.Error()))
			continue
		}
		log.Debug("service.StartWorker got notification", slog.Int64("notification", notification.ChatID))
		err = ns.Bot.SendMessage(ctx, notification.ChatID, notification.Message)
		if err != nil {
			log.Error("service.StartWorker error sending message to notifier",
				slog.String("error", err.Error()))
			log.Info("service.StartWorker pushing task back to queue (Retry)",
				slog.Int64("notification", notification.ChatID))
			pushErr := ns.Queue.PushQueue(ctx, QueueName, jsonStr)
			if pushErr != nil {
				log.Error("service.StartWorker CRITICAL error: Failed to push task back to queue. Message lost!",
					slog.String("error", pushErr.Error()))
				time.Sleep(2 * time.Second)
			} else {
				log.Info("service.StartWorker Successfully processed and retry ",
					slog.Int64("notification", notification.ChatID))
			}

		}

	}
}
