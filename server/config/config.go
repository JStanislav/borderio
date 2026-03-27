package config

import "github.com/kelseyhightower/envconfig"

type Config struct {
	Port string `default:"8080"`
}

func LoadConfig() Config {
	var config Config
	envconfig.Process("", &config)
	return config
}
