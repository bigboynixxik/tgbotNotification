package telegram

import (
	"TGNotification/internal/clients"
	"TGNotification/pkg/logger"
	"context"
	"fmt"
	"log/slog"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	botAPI       *tgbotapi.BotAPI
	DjangoClient *clients.DjangoClient
}

func NewBot(token string, djangoClient *clients.DjangoClient) (*Bot, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, fmt.Errorf("tg.NewBot,error creating bot api: %w", err)
	}
	slog.Info("Created telegram BotAPI")
	return &Bot{
		botAPI:       bot,
		DjangoClient: djangoClient,
	}, nil
}

// Start - запускает бесконечный цикл прослушивания обновлений от Telegram
func (b *Bot) Start(ctx context.Context) {
	log := logger.FromContext(ctx).With("component", "telegram_bot")
	log.Info("Starting Telegram Bot", slog.String("account", b.botAPI.Self.UserName))

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60 // Ждем 60 секунд. Если сообщений нет, переподключаемся (экономит трафик)

	updates := b.botAPI.GetUpdatesChan(u)

	for update := range updates {

		if update.Message == nil {
			continue
		}

		msgCtx := logger.IntoContext(ctx, log.With(
			"chat_id", update.Message.Chat.ID,
			"username", update.Message.From.UserName,
		))

		logMsg := logger.FromContext(msgCtx)

		logMsg.Debug("Received message",
			slog.String("username", update.Message.From.UserName),
			slog.String("text", update.Message.Text),
			slog.Int64("chat_id", update.Message.Chat.ID), // ID чата - самое важное для нас!
		)

		if update.Message.IsCommand() {
			b.handleCommand(msgCtx, update.Message)
			continue
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Я бот-уведомитель. Я понимаю только команды. Нажми /start")

		_, err := b.botAPI.Send(msg)
		if err != nil {
			logMsg.Error("tg.start, Failed to send message", slog.String("error", err.Error()))
		}
	}
}

// handleCommand - маршрутизатор команд
func (b *Bot) handleCommand(ctx context.Context, message *tgbotapi.Message) {
	log := logger.FromContext(ctx)
	var responseText string

	switch message.Command() {
	case "start":
		log.Info("User started bot")
		args := message.CommandArguments()
		if len(args) == 0 {
			responseText = "🐾 Привет! Я бот зоомагазина.\nВаш Chat ID: " + fmt.Sprint(message.Chat.ID) + "\nСообщите этот ID администратору для привязки уведомлений."
		} else {
			resp, msg, err := b.DjangoClient.LinkUser(ctx, args, message.From.ID, message.From.UserName)
			if err != nil || !resp {
				log.Error("telegram.handleCommand, Failed to link user", slog.String("error", fmt.Sprint(err)))
				responseText = "Не получилось привязать ваш токен. Попробуйте ещё раз."
			} else {
				// Если всё ок, выводим сообщение, которое нам прислал Django
				responseText = msg
			}
		}
	case "help":
		responseText = "Я могу отправлять вам уведомления о статусе заказов. Ожидайте новых сообщений!"

	default:
		responseText = "Неизвестная команда. Попробуйте /start или /help"
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, responseText)
	_, err := b.botAPI.Send(msg)
	if err != nil {
		log.Error("tg.handleCommand, Failed to send command response", slog.String("error", err.Error()))
	}
}

func (b *Bot) SendMessage(ctx context.Context, chatID int64, message string) error {
	log := logger.FromContext(ctx)
	msg := tgbotapi.NewMessage(chatID, message)
	msg.ParseMode = tgbotapi.ModeHTML

	_, err := b.botAPI.Send(msg)
	if err != nil {
		log.Error("tg.sendMessage, Failed to send message", slog.String("error", err.Error()))
		return fmt.Errorf("tg.sendMessage, Failed to send message: %w", err)
	}
	log.Info("Successfully sent message", slog.Int64("chat_id", chatID))
	return nil
}
