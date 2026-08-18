package configs

import (
	"fmt"

	"github.com/spf13/viper"
)

func NewViper() (*viper.Viper, error) {
	config := viper.New()

	config.AddConfigPath(".")
	config.SetConfigFile(".env")

	err := config.ReadInConfig()

	if err != nil {
		return nil, fmt.Errorf("error setup viper configuration: %w", err)
	}

	return config, nil
}
