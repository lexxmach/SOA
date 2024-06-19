package main

import (
	apiHandlers "SOA/cmd/users/handlers/api"
	postsHandlers "SOA/cmd/users/handlers/posts"

	proto "SOA/proto/api"

	"SOA/internal/api"
	"SOA/internal/common"
	"SOA/internal/db"
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"gorm.io/driver/postgres"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
)

func main() {
	config := setup()

	// TOOD: change to golang 1.22 router, add middleware
	router := chi.NewMux()
	api := humachi.New(router, huma.DefaultConfig(config.Title, config.Version))

	// CreateDB
	var database db.ApiDatabase
	{
		if config.DB.Mock {
			database = db.CreateMockDB()
		} else {
			dsn := fmt.Sprintf(
				"host=%s user=%s password=%s dbname=postgres port=%d sslmode=disable",
				config.DB.Host,
				config.DB.User,
				config.DB.Password,
				config.DB.Port,
			)

			database = common.Must(db.CreateApiGormDB(postgres.Open(dsn)))
		}
	}

	// TODO: make generic, this sucks
	// api section
	{
		registerHandler := apiHandlers.RegisterHandler{
			DB:        database,
			JWTSecret: config.JWTSecret,
		}
		huma.Register(api, apiHandlers.RegisterOperation, registerHandler.Handle)

		authHandler := apiHandlers.AuthHandler{
			DB:        database,
			JWTSecret: config.JWTSecret,
		}
		huma.Register(api, apiHandlers.AuthOperation, authHandler.Handle)

		updateHandler := apiHandlers.UpdateHandler{
			DB:        database,
			JWTSecret: config.JWTSecret,
		}
		huma.Register(api, apiHandlers.UpdateOperation, updateHandler.Handle)

		getHandler := apiHandlers.GetHandler{
			DB: database,
		}
		huma.Register(api, apiHandlers.GetOperation, getHandler.Handle)
	}

	// posts section
	{
		client := proto.NewPostsServiceClient(
			common.Must(
				grpc.NewClient(
					config.PostAddress,
					grpc.WithTransportCredentials(insecure.NewCredentials()),
				),
			),
		)

		createHandler := postsHandlers.CreateHandler{
			Client:    client,
			JWTSecret: config.JWTSecret,
		}
		huma.Register(api, postsHandlers.CreateOperation, createHandler.Handle)

		deleteHandler := postsHandlers.DeleteHandler{
			Client:    client,
			JWTSecret: config.JWTSecret,
		}
		huma.Register(api, postsHandlers.DeleteOperation, deleteHandler.Handle)

		getPostHandler := postsHandlers.GetPostHandler{
			Client:    client,
			JWTSecret: config.JWTSecret,
		}
		huma.Register(api, postsHandlers.GetPostOperation, getPostHandler.Handle)

		listPostHandler := postsHandlers.ListPostHandler{
			Client: client,
		}
		huma.Register(api, postsHandlers.ListPostOperation, listPostHandler.Handle)

		updatePostHandler := postsHandlers.UpdatePostHandler{
			Client: client,
		}
		huma.Register(api, postsHandlers.UpdatePostOperation, updatePostHandler.Handle)
	}

	logrus.Info("Starting server at port ", config.Port)

	if err := http.ListenAndServe(fmt.Sprintf(":%d", config.Port), router); err != nil {
		panic(err)
	}
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
