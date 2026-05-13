package config

import (
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type GRPCConfig struct {
	Address string `yaml:"address" env:"GRPC_ADDRESS" env-default:"localhost:50051"`
	Timeout int    `yaml:"timeout" env:"GRPC_TIMEOUT" env-default:"5"`
}

type Config struct {
	GRPC GRPCConfig `yaml:"grpc"`
}

func MustLoad(path string) Config {
	var cfg Config
	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		if os.IsNotExist(err) {
			cleanenv.ReadEnv(&cfg)
			return cfg
		}
		panic("failed to read config: " + err.Error())
	}
	return cfg
}
