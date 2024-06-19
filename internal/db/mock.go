package db

import (
	"SOA/internal/api"
	"SOA/internal/posts"
	"fmt"
)

type MockDB struct {
	apiStorage   map[api.UserLogin]*api.User
	postsStorage map[uint64]*posts.Post
}

func CreateApiMockDB() ApiDatabase {
	return &MockDB{
		apiStorage: make(map[api.UserLogin]*api.User),
	}
}

func CreatePostsMockDB() PostDatabase {
	return &MockDB{
		postsStorage: make(map[uint64]*posts.Post),
	}
}

func (db *MockDB) CreateUser(user api.User) error {
	if _, ok := db.apiStorage[user.GetLogin()]; ok {
		return fmt.Errorf("user with login %q already exists", user.GetLogin())
	}

	db.apiStorage[user.GetLogin()] = &user
	return nil
}

func (db *MockDB) GetUser(login api.UserLogin) (*api.User, error) {
	if _, ok := db.apiStorage[login]; !ok {
		return nil, fmt.Errorf("user with login %q doesn't exist", login)
	}

	return db.apiStorage[login], nil
}

func (db *MockDB) UpdateUser(user *api.User) error {
	if user == nil {
		return fmt.Errorf("user is nil")
	}

	apiStorageUser, err := db.GetUser(user.GetLogin())
	if err != nil {
		return err
	}

	*apiStorageUser = *user
	return nil
}

func (db *MockDB) GetNewPostId() uint64 {
	return uint64(len(db.postsStorage) + 1)
}

func (db *MockDB) CreatePost(post posts.Post) (uint64, error) {
	post.ID = db.GetNewPostId()
	db.postsStorage[post.ID] = &post

	return post.ID, nil
}

// DeletePost implements PostDatabase.
func (db *MockDB) DeletePost(id uint64) error {
	if _, ok := db.postsStorage[id]; !ok {
		return fmt.Errorf("failed to get post %d", id)
	}
	delete(db.postsStorage, id)

	return nil
}

func (db *MockDB) GetPost(id uint64) (*posts.Post, error) {
	if _, ok := db.postsStorage[id]; !ok {
		return nil, fmt.Errorf("failed to get post %d", id)
	}
	return db.postsStorage[id], nil
}

func (db *MockDB) ListPosts(pageNum uint64, pageSize uint64) ([]posts.Post, error) {
	var boiler []posts.Post
	for id := pageNum * pageSize; id < (pageNum+1)*pageSize; id++ {
		post, err := db.GetPost(id)
		if err != nil {
			return nil, err
		}
		boiler = append(boiler, *post)
	}
	return boiler, nil
}

// UpdatePost implements PostDatabase.
func (db *MockDB) UpdatePost(post *posts.Post) error {
	if _, err := db.GetPost(post.ID); err != nil {
		return err
	}

	db.postsStorage[post.ID] = post
	return nil
}
