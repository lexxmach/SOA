package db

import (
	"SOA/internal/api"
	"SOA/internal/posts"
)

type ApiDatabase interface {
	CreateUser(api.User) error
	UpdateUser(*api.User) error

	GetUser(api.UserLogin) (*api.User, error)
}

type PostDatabase interface {
	CreatePost(posts.Post) (uint64, error)
	UpdatePost(*posts.Post) error
	DeletePost(uint64) error
	GetPost(uint64) (*posts.Post, error)

	ListPosts(uint64, uint64) ([]posts.Post, error)
}
