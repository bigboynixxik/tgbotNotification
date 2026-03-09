package service

import (
	"context"
	"errors"
	"testing"
	"time"
)

type mockBot struct {
	shouldFail bool
}

func (m *mockBot) SendMessage(ctx context.Context, chatID int64, message string) error {
	if m.shouldFail {
		return errors.New("fake network error") // Имитируем падение Телеграма
	}
	return nil
}

type mockQueue struct {
	tasksToGive     []string
	tasksPushed     []string
	listenCallCount int
}

func (m *mockQueue) ListenQueue(ctx context.Context, queueName string) (string, error) {
	if m.listenCallCount >= len(m.tasksToGive) {
		return "", errors.New("no more tasks in mock queue")
	}
	task := m.tasksToGive[m.listenCallCount]
	m.listenCallCount++
	return task, nil
}

func (m *mockQueue) PushQueue(ctx context.Context, queueName string, message string) error {
	m.tasksPushed = append(m.tasksPushed, message)
	return nil
}

func TestNotifierService_WorkerRetry(t *testing.T) {
	jsonTask := `{"chat_id": 111, "message": "test"}`

	fakeBot := &mockBot{shouldFail: true}
	fakeQueue := &mockQueue{tasksToGive: []string{jsonTask}}

	svc := NewNotifierService(fakeBot, fakeQueue)

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	_ = svc.StartWorker(ctx)

	if len(fakeQueue.tasksPushed) == 0 {
		t.Fatal("Expected task to be pushed back to queue (Retry), but it wasn't")
	}

	if fakeQueue.tasksPushed[0] != jsonTask {
		t.Errorf("Expected pushed task to be '%s', got '%s'", jsonTask, fakeQueue.tasksPushed[0])
	}
}
