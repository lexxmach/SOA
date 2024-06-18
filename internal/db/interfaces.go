package db

import "SOA/internal/api"

type Database interface {
	CreateUser(api.User) error
	UpdateUser(*api.User) error

	GetUser(api.UserLogin) (*api.User, error)
}
