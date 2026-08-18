package app

import (
	"log"

	"github.com/Andhika-GIT/go-message-broker-monorepo/configs"
	"github.com/go-chi/chi/v5"
)

func InitApp() *chi.Mux {
	r := chi.NewRouter()

	v, err := configs.NewViper()

	if err != nil {
		log.Fatal("error when initiating viper")
	}

	configs.InitConfig(v)

	return r
}
