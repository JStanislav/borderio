package config

import "github.com/kelseyhightower/envconfig"

type Config struct {
	Port                 string `default:"8080"`
	TimeoutAfterGameOver int    `default:"60" split_words:"true"` // in seconds

	Cors Cors
}

type Cors struct {
	AllowedOrigins []string `required:"true" split_words:"true"`
	AllowedMethods []string `default:"GET,POST,PUT,PATCH,DELETE,OPTIONS" split_words:"true"`
}

func LoadConfig() Config {
	var config Config
	if err := envconfig.Process("", &config); err != nil {
		panic(err)
	}

	return config
}
