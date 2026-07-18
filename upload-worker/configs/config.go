package configs

import "github.com/spf13/viper"

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
}

type RabbitMQRoutingKey struct {
	UserDirectImport  string
	UserSftpImport    string
	UserExport        string
	OrderDirectImport string
	OrderImport       string
	OrderExport       string
}

type RabbitMQQueue struct {
	UserDirectImport  string
	UserSftpImport    string
	UserExport        string
	OrderDirectImport string
	OrderImport       string
	OrderExport       string
}

type RedisClientConfig struct {
	Addr     string
	Password string
	DB       int
}

type SftpClientConfig struct {
	Host     string
	Port     int
	Username string
	Password string
}

type Config struct {
	// RabbitMQ
	RabbitMQConnectURL string
	RabbitMQExchange   string

	RabbitMQRoutingKey RabbitMQRoutingKey

	RabbitMQQueue RabbitMQQueue

	// Database
	Database DatabaseConfig

	// redis client
	RedisClient RedisClientConfig

	// sftp/ssh
	SftpClient SftpClientConfig
}

func InitConfig(v *viper.Viper) *Config {
	cfg := &Config{}

	// --- RabbitMQ connection ---
	cfg.RabbitMQConnectURL = getOrDefaultString(v, "RABBITMQ_CONNECTION_URL", "amqp://guest:guest@localhost:5672/")
	cfg.RabbitMQExchange = getOrDefaultString(v, "MQ_EXCHANGE_GO_APP", "go-app-exchange")

	// --- RabbitMQ Routing Keys ---
	cfg.RabbitMQRoutingKey = RabbitMQRoutingKey{
		UserDirectImport:  getOrDefaultString(v, "MQ_RK_USER_DIRECT_IMPORT", "user.import.direct"),
		UserSftpImport:    getOrDefaultString(v, "MQ_RK_USER_SFTP_IMPORT", "user.import.sftp"),
		UserExport:        getOrDefaultString(v, "MQ_RK_USER_EXPORT", "user.export"),
		OrderDirectImport: getOrDefaultString(v, "MQ_RK_ORDER_DIRECT_IMPORT", "order.import.direct"),
		OrderImport:       getOrDefaultString(v, "MQ_RK_ORDER_IMPORT", "order.import"),
		OrderExport:       getOrDefaultString(v, "MQ_RK_ORDER_EXPORT", "order.export"),
	}

	// --- RabbitMQ Queues ---
	cfg.RabbitMQQueue = RabbitMQQueue{
		UserDirectImport:  getOrDefaultString(v, "MQ_Q_USER_DIRECT_IMPORT", "user.import.direct.q"),
		UserSftpImport:    getOrDefaultString(v, "MQ_Q_USER_SFTP_IMPORT", "user.import.sftp.q"),
		UserExport:        getOrDefaultString(v, "MQ_Q_USER_EXPORT", "user.export.q"),
		OrderDirectImport: getOrDefaultString(v, "MQ_Q_ORDER_DIRECT_IMPORT", "order.import.direct.q"),
		OrderImport:       getOrDefaultString(v, "MQ_Q_ORDER_IMPORT", "order.import.q"),
		OrderExport:       getOrDefaultString(v, "MQ_Q_ORDER_EXPORT", "order.export.q"),
	}

	// --- Database ---
	cfg.Database = DatabaseConfig{
		Host:     getOrDefaultString(v, "DB_HOST", "localhost"),
		Port:     getOrDefaultInt(v, "DB_PORT", 5432),
		User:     getOrDefaultString(v, "DB_USERNAME", "postgres"),
		Password: getOrDefaultString(v, "DB_PASSWORD", "postgres"),
		Name:     getOrDefaultString(v, "DB_NAME", "postgres"),
	}

	// --- Redis Client ---
	cfg.RedisClient = RedisClientConfig{
		Addr:     getOrDefaultString(v, "REDIS_ADDR", "localhost:6379"),
		Password: getOrDefaultString(v, "REDIS_PASSWORD", ""),
		DB:       getOrDefaultInt(v, "REDIS_DB", 0),
	}

	// --- SFTP ---
	cfg.SftpClient = SftpClientConfig{
		Host:     getOrDefaultString(v, "SFTP_HOST", "localhost"),
		Port:     getOrDefaultInt(v, "SFTP_PORT", 2222),
		Username: getOrDefaultString(v, "SFTP_USERNAME", "foo"),
		Password: getOrDefaultString(v, "SFTP_PASSWORD", "pass"),
	}

	return cfg
}

func getOrDefaultString(v *viper.Viper, key, def string) string {
	val := v.GetString(key)
	if val == "" {
		return def
	}
	return val
}

func getOrDefaultInt(v *viper.Viper, key string, def int) int {
	val := v.GetInt(key)
	if val == 0 {
		return def
	}
	return val
}
