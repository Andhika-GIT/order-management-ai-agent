package configs

import "github.com/spf13/viper"

func OrderImportConsumer(v *viper.Viper) RabbitMQConfig {
	return RabbitMQConfig{
		URL:          getOrDefaultString(v, "RABBITMQ_CONNECTION_URL", "amqp://guest:guest@localhost:5672/"),
		Queue:        getOrDefaultString(v, "MQ_Q_ORDER_IMPORT", "order.import.q"),
		RoutingKey:   getOrDefaultString(v, "MQ_RK_ORDER_IMPORT", "order.import"),
		ExchangeName: getOrDefaultString(v, "MQ_EXCHANGE_NAME", "go-app-exchange"),
		ExchangeType: getOrDefaultString(v, "MQ_EXCHANGE_TYPE", "direct"),
	}
}

func UserImportConsumer(v *viper.Viper) RabbitMQConfig {
	return RabbitMQConfig{
		URL:          getOrDefaultString(v, "RABBITMQ_CONNECTION_URL", "amqp://guest:guest@localhost:5672/"),
		Queue:        getOrDefaultString(v, "MQ_Q_USER_IMPORT", "user.import.q"),
		RoutingKey:   getOrDefaultString(v, "MQ_RK_USER_IMPORT", "user.import"),
		ExchangeName: getOrDefaultString(v, "MQ_EXCHANGE_NAME", "go-app-exchange"),
		ExchangeType: getOrDefaultString(v, "MQ_EXCHANGE_TYPE", "direct"),
	}
}
