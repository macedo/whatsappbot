package main

import (
	"github.com/spf13/viper"
)

type AppConfig struct {
	Bucket      string `mapstructure:"BUCKET"`
	DatabaseURL string `mapstructure:"DATABASE_URL"`
	JIDs        string `mapstructure:"JIDS"`
}

func LoadConfig(path string) (AppConfig, error) {
	var cfg AppConfig

	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		return cfg, err
	}

	err = viper.Unmarshal(&cfg)

	return cfg, err
}

func LoadDefaultConfig() (AppConfig, error) {
	return LoadConfig(".")
}
