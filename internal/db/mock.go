package db

import (
	"SOA/internal/api"
	"fmt"
)

type MockDB struct {
	storage map[api.UserLogin]*api.User
}

func CreateMockDB() *MockDB {
	return &MockDB{
		storage: make(map[api.UserLogin]*api.User),
	}
}

func (db *MockDB) CreateUser(user api.User) error {
	if _, ok := db.storage[user.GetLogin()]; ok {
		return fmt.Errorf("user with login %q already exists", user.GetLogin())
	}

	db.storage[user.GetLogin()] = &user
	return nil
}

func (db *MockDB) GetUser(login api.UserLogin) (*api.User, error) {
	if _, ok := db.storage[login]; !ok {
		return nil, fmt.Errorf("user with login %q doesn't exist", login)
	}

	return db.storage[login], nil
}

func (db *MockDB) UpdateUser(user *api.User) error {
	if user == nil {
		return fmt.Errorf("user is nil")
	}

	storageUser, err := db.GetUser(user.GetLogin())
	if err != nil {
		return err
	}

	*storageUser = *user
	return nil
}
