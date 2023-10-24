package config

import (
	"log/slog"

	"github.com/rs/zerolog/log"

	"github.com/spf13/viper"
)

type Config struct {
	APP_NAME                   string `mapstructure:"app_name"`
	APP_ENV                    string `mapstructure:"app_env"`
	APP_KEY                    string `mapstructure:"app_key"`
	DB_URL                     string `mapstructure:"db_url"`
	PORT                       int    `mapstructure:"port"`
	DISABLE_GRAPHQL_PLAYGROUND bool   `mapstructure:"disable_graphql_playground"`
}

func New() *Config {
	config := Config{}

	viper.AddConfigPath(".")
	viper.SetConfigFile(".env.yml")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			slog.Warn("failed to read config file")
		}
	}

	if err := viper.BindEnv("APP_NAME"); err != nil {
		log.Err(err).Msg("failed to bind environment variable")
	}

	if err := viper.BindEnv("APP_ENV"); err != nil {
		log.Err(err).Msg("failed to bind environment variable")
	}

	if err := viper.BindEnv("APP_KEY"); err != nil {
		log.Err(err).Msg("failed to bind environment variable")
	}

	if err := viper.BindEnv("DB_URL"); err != nil {
		log.Err(err).Msg("failed to bind environment variable")
	}

	if err := viper.BindEnv("PORT"); err != nil {
		log.Err(err).Msg("failed to bind environment variable")
	}

	if err := viper.BindEnv("JWT_SECRET"); err != nil {
		log.Err(err).Msg("failed to bind environment variable")
	}

	if err := viper.Unmarshal(&config); err != nil {
		log.Err(err).Msg("unable to parse config file")
	}

	return &config
}
