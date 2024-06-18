package main

import (
	"SOA/cmd/api/handlers"
	"SOA/internal/api"
	"SOA/internal/db"
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
)

func main() {
	config := setup()

	router := chi.NewMux()
	api := humachi.New(router, huma.DefaultConfig(config.Title, config.Version))

	// CreateDB
	var database db.Database
	{
		if config.DB.Mock {
			database = db.CreateMockDB()
		} else {
			panic("UNKNOWN DB")
		}
	}

	// auth section
	{
		registerHandler := handlers.RegisterHandler{
			DB:        database,
			JWTSecret: config.JWTSecret,
		}

		huma.Register(api, handlers.RegisterOperation, registerHandler.Handle)

		authHandler := handlers.AuthHandler{
			DB:        database,
			JWTSecret: config.JWTSecret,
		}

		huma.Register(api, handlers.AuthOperation, authHandler.Handle)
	}

	// operations section
	{
		updateHandler := handlers.UpdateHandler{
			DB:        database,
			JWTSecret: config.JWTSecret,
		}

		huma.Register(api, handlers.UpdateOperation, updateHandler.Handle)
	}

	logrus.Info("Starting server at port ", config.Port)

	http.ListenAndServe(fmt.Sprintf(":%d", config.Port), router)
}

func setup() *api.Config {
	configPath := pflag.StringP("config", "c", "", "path to config file")

	pflag.Parse()

	if configPath == nil || *configPath == "" {
		panic("config path must not be nil or empty")
	}

	logrus.Info("Starting server...")

	return api.MustGetConfig(*configPath)
}
