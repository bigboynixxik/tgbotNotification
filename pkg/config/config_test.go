package config

import (
	"testing"
)

func TestLoadConfig_Error(t *testing.T) {
	// Пытаемся загрузить конфиг из несуществующего файла
	_, err := LoadConfig("nil.env")

	if err == nil {
		t.Error("Expected error when loading non-existent config, but got nil")
	}
}
