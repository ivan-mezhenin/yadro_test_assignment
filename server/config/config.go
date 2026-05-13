package config

import "github.com/ilyakaznacheev/cleanenv"

type GRPCConfig struct {
	Port int `yaml:"port" env:"GRPC_PORT" env-default:"50051"`
}

type ServerConfig struct {
	ResolvConf string `yaml:"resolv_conf" env:"RESOLV_CONF" env-default:"/etc/resolv.conf"`
	BackupPath string `yaml:"backup_path" env:"BACKUP_PATH" env-default:"/etc/resolv.conf.bak"`
}

type LoggerConfig struct {
	Level string `yaml:"level" env:"LOG_LEVEL" env-default:"info"`
	Debug bool   `yaml:"debug" env:"DEBUG" env-default:"false"`
}

type Config struct {
	GRPC   GRPCConfig   `yaml:"grpc"`
	Server ServerConfig `yaml:"server"`
	Logger LoggerConfig `yaml:"logger"`
}

func MustLoad(path string) Config {
	var cfg Config
	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		panic("failed to read config: " + err.Error())
	}
	return cfg
}
