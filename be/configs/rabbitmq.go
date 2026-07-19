package configs

import "github.com/spf13/viper"

func Producer(v *viper.Viper) RabbitMQConfig {
	return RabbitMQConfig{
		URL:          getOrDefaultString(v, "RABBITMQ_CONNECTION_URL", "amqp://guest:guest@localhost:5672/"),
		ExchangeName: getOrDefaultString(v, "MQ_EXCHANGE_NAME", "go-app-exchange"),
		ExchangeType: getOrDefaultString(v, "MQ_EXCHANGE_TYPE", "direct"),
	}
}
