package service

import "context"

type BotProvider interface {
	SendMessage(ctx context.Context, chatID int64, message string) error
}

type QueueProvider interface {
	ListenQueue(ctx context.Context, queueName string) (string, error)
	PushQueue(ctx context.Context, queueName string, message string) error
}
