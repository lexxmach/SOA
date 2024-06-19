package main

import (
	"SOA/cmd/posts/handlers"
	"SOA/internal/common"
	"SOA/internal/db"
	posts "SOA/internal/posts"
	"SOA/proto/api"
	"fmt"
	"net"

	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"google.golang.org/grpc"
	"gorm.io/driver/postgres"
)

func main() {
	config := setup()

	// CreateDB
	var database db.PostDatabase
	{
		if config.DB.Mock {
			panic("unsupported mock db")
		} else {
			dsn := fmt.Sprintf(
				"host=%s user=%s password=%s dbname=postgres port=%d sslmode=disable",
				config.DB.Host,
				config.DB.User,
				config.DB.Password,
				config.DB.Port,
			)

			database = common.Must(db.CreatePostsGormDB(postgres.Open(dsn)))
		}
	}

	server := grpc.NewServer()
	listen, err := net.Listen("tcp", config.Address)
	if err != nil {
		panic(err)
	}

	api.RegisterPostsServiceServer(server, handlers.PostGRPCHandler{
		DB: database,
	})

	if err := server.Serve(listen); err != nil {
		panic(err)
	}
}

func setup() *posts.Config {
	configPath := pflag.StringP("config", "c", "", "path to config file")

	pflag.Parse()

	if configPath == nil || *configPath == "" {
		panic("config path must not be nil or empty")
	}

	logrus.Info("Starting server...")

	return posts.MustGetConfig(*configPath)
}
