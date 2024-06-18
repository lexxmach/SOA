package db

import (
	"SOA/internal/api"
	"fmt"

	"gorm.io/gorm"
)

type GormDatabase struct {
	db *gorm.DB
}

func CreateGormDB(dialector gorm.Dialector) (*GormDatabase, error) {
	gormDB, err := gorm.Open(dialector)
	if err != nil {
		return nil, fmt.Errorf("failed to setup gorm: %w", err)
	}

	gormDB.AutoMigrate(&api.User{})

	return &GormDatabase{
		db: gormDB,
	}, nil
}

func (gorm *GormDatabase) CreateUser(user api.User) error {
	tx := gorm.db.Create(&user)
	if tx.Error != nil {
		return fmt.Errorf("failed to create user %q: %w", user, tx.Error)
	}

	return nil
}

func (gorm *GormDatabase) GetUser(login api.UserLogin) (*api.User, error) {
	user := &api.User{}

	tx := gorm.db.First(&user, login)
	if tx.Error != nil {
		return nil, fmt.Errorf("failed to get user %q: %w", login, tx.Error)
	}

	return user, nil
}

func (gorm *GormDatabase) UpdateUser(user *api.User) error {
	tx := gorm.db.Save(user)
	if tx.Error != nil {
		return fmt.Errorf("failed to update user %q: %w", user, tx.Error)
	}

	return nil
}
