package configs

import (
	"log"

	"github.com/spf13/viper"
)

type AppConfig struct {
	Port string
	Env  string
}

type LLMConfig struct {
	Provider       string
	BaseURL        string
	APIKey         string
	Model          string
	MaxTokens      int
	TimeoutSeconds int
}

type RedisClientConfig struct {
	Addr     string
	Password string
	DB       int
}

type Config struct {
	App AppConfig
	LLM LLMConfig

	// redis client
	RedisClient RedisClientConfig
}

func InitConfig(v *viper.Viper) *Config {
	cfg := &Config{}

	cfg.App = AppConfig{
		Port: getOrDefaultString(v, "APP_PORT", "3011"),
		Env:  getOrDefaultString(v, "APP_ENV", "development"),
	}

	provider := getOrDefaultString(v, "LLM_PROVIDER", "anthropic")
	cfg.LLM = LLMConfig{
		Provider:       provider,
		BaseURL:        getOrDefaultString(v, "LLM_BASE_URL", defaultLLMBaseURL(provider)),
		APIKey:         v.GetString("LLM_API_KEY"), // sengaja tanpa default, wajib diisi manual
		Model:          getOrDefaultString(v, "LLM_MODEL", "claude-sonnet-5"),
		MaxTokens:      getOrDefaultInt(v, "LLM_MAX_TOKENS", 1024),
		TimeoutSeconds: getOrDefaultInt(v, "LLM_TIMEOUT_SECONDS", 60),
	}

	// --- Redis Client ---
	cfg.RedisClient = RedisClientConfig{
		Addr:     getOrDefaultString(v, "REDIS_ADDR", "localhost:6379"),
		Password: getOrDefaultString(v, "REDIS_PASSWORD", ""),
		DB:       getOrDefaultInt(v, "REDIS_DB", 0),
	}

	validate(cfg)

	return cfg
}

func validate(cfg *Config) {
	if cfg.LLM.Provider != "anthropic" && cfg.LLM.Provider != "openai_compatible" {
		log.Fatalf("[config] LLM_PROVIDER tidak dikenal: %q (pakai \"anthropic\" atau \"openai_compatible\")", cfg.LLM.Provider)
	}
	if cfg.LLM.Model == "" {
		log.Fatal("[config] LLM_MODEL wajib diisi")
	}
	if cfg.LLM.APIKey == "" {
		log.Println("[config] warning: LLM_API_KEY kosong — oke kalau target-nya Ollama lokal, tapi wajib untuk provider cloud")
	}
}

func defaultLLMBaseURL(provider string) string {
	switch provider {
	case "anthropic":
		return "https://api.anthropic.com"
	case "openai_compatible":
		return "https://api.openai.com/v1"
	default:
		return ""
	}
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
