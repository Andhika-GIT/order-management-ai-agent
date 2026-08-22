package configs

import "github.com/spf13/viper"

func ProductImportConsumer(v *viper.Viper) RabbitMQConfig {
	return RabbitMQConfig{
		URL:          getOrDefaultString(v, "RABBITMQ_CONNECTION_URL", "amqp://guest:guest@localhost:5672/"),
		Queue:        getOrDefaultString(v, "MQ_Q_PRODUCT_IMPORT", "product.import.q"),
		RoutingKey:   getOrDefaultString(v, "MQ_RK_PRODUCT_IMPORT", "product.import"),
		ExchangeName: getOrDefaultString(v, "MQ_EXCHANGE_NAME", "go-app-exchange"),
		ExchangeType: getOrDefaultString(v, "MQ_EXCHANGE_TYPE", "direct"),
	}
}
