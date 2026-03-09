package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	TGToken        string `mapstructure:"TG_TOKEN"`
	AppEnv         string `mapstructure:"APP_ENV"`
	RedisAddr      string `mapstructure:"REDIS_ADDR"`
	DjangoGRPCAddr string `mapstructure:"DJANGO_GRPC_ADDR"`
}

func LoadConfig(path string) (*Config, error) {
	v := viper.New()
	v.SetConfigFile(path)
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("config.LoadConfig: ошибка чтения конфига %w", err)
	}
	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("config.LoadConfig: Ошибка анмаршала конфига %w", err)
	}
	return &cfg, nil
}
