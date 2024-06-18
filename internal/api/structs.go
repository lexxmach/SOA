package api

import (
	"net/mail"
	"time"
)

type UserLogin struct {
	Login string `gorm:"primaryKey;<-:create"`
}

type UserCredentials struct {
	Login UserLogin `gorm:"embedded"`

	// Contains already salted password
	Password string `gorm:"<-:create"`
}

type UserToken struct {
	// JWT token
	Token string `json:"oauth"`
}

type User struct {
	FirstName string
	LastName  string

	BirthDate time.Time
	Email     mail.Address

	// TODO(lexmach): use some fancy package
	Phone string

	Creds UserCredentials `gorm:"embedded"`
}

func (u *User) GetLogin() UserLogin {
	return u.Creds.Login
}

func (u *User) GetPassword() string {
	return u.Creds.Password
}
