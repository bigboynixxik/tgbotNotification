package models

type Notification struct {
	ChatID  int64  `json:"chat_id"`
	Message string `json:"message"`
}
