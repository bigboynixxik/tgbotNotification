package models

import (
	"encoding/json"
	"testing"
)

// TestNotification_Unmarshal проверяет правильность парсинга JSON
func TestNotification_Unmarshal(t *testing.T) {
	jsonStr := `{"chat_id": 12345678, "message": "Ваш заказ готов!"}`

	var n Notification
	err := json.Unmarshal([]byte(jsonStr), &n)

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err) // Fatalf сразу остановит тест
	}

	if n.ChatID != 12345678 {
		t.Errorf("Expected ChatID 12345678, got: %d", n.ChatID)
	}

	if n.Message != "Ваш заказ готов!" {
		t.Errorf("Expected Message 'Ваш заказ готов!', got: '%s'", n.Message)
	}
}
