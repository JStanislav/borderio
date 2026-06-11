package config

import "github.com/kelseyhightower/envconfig"

type Config struct {
	Port                 string `default:"8080"`
	TimeoutAfterGameOver int    `default:"60" split_words:"true"` // in seconds
}

func LoadConfig() Config {
	var config Config
	envconfig.Process("", &config)
	return config
}
