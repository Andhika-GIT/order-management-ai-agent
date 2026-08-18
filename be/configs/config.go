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
	UserImport  string
	UserExport  string
	OrderImport string
	OrderExport string
}

type RabbitMQConfig struct {
	URL          string
	ExchangeName string
	ExchangeType string
}

type RedisClientConfig struct {
	Addr     string
	Password string
	DB       int
}

type S3Config struct {
	Endpoint       string
	PublicEndpoint string
	Region         string
	AccessKey      string
	SecretKey      string
	Bucket         string
	UseSSL         bool
	UsePathStyle   bool
}

type Config struct {
	// RabbitMQ
	RabbitMQRoutingKey RabbitMQRoutingKey

	// Database
	Database DatabaseConfig

	// redis client
	RedisClient RedisClientConfig

	// S3
	S3Config S3Config
}

func InitConfig(v *viper.Viper) *Config {
	cfg := &Config{}

	// --- RabbitMQ Routing Keys ---
	cfg.RabbitMQRoutingKey = RabbitMQRoutingKey{
		UserImport:  getOrDefaultString(v, "MQ_RK_USER_IMPORT", "user.import"),
		UserExport:  getOrDefaultString(v, "MQ_RK_USER_EXPORT", "user.export"),
		OrderImport: getOrDefaultString(v, "MQ_RK_ORDER_IMPORT", "order.import"),
		OrderExport: getOrDefaultString(v, "MQ_RK_ORDER_EXPORT", "order.export"),
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

	// --- S3 ---
	cfg.S3Config = S3Config{
		Endpoint:       getOrDefaultString(v, "S3_ENDPOINT", ""),
		PublicEndpoint: getOrDefaultString(v, "S3_PUBLIC_ENDPOINT", ""),
		Region:         getOrDefaultString(v, "S3_REGION", "us-east-1"),
		AccessKey:      getOrDefaultString(v, "S3_ACCESS_KEY", ""),
		SecretKey:      getOrDefaultString(v, "S3_SECRET_KEY", ""),
		Bucket:         getOrDefaultString(v, "S3_BUCKET", ""),
		UseSSL:         getOrDefaultBool(v, "S3_USE_SSL", true),
		UsePathStyle:   getOrDefaultBool(v, "S3_USE_PATH_STYLE", false),
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

func getOrDefaultBool(v *viper.Viper, key string, def bool) bool {
	if !v.IsSet(key) {
		return def
	}
	return v.GetBool(key)
}
