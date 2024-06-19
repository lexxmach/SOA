package api

import (
	"database/sql/driver"
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

type Address mail.Address

type User struct {
	FirstName string
	LastName  string

	BirthDate time.Time
	Email     Address

	// TODO(lexmach): use some fancy package
	Phone string

	Creds UserCredentials `gorm:"embedded"`
}

func ParseAddress(input string) (*Address, error) {
	val, err := mail.ParseAddress(input)
	if err != nil {
		return nil, err
	}
	valueAddr := &Address{}
	*valueAddr = Address(*val)
	return valueAddr, nil
}

func (adr *Address) Scan(value interface{}) error {
	parsed, err := mail.ParseAddress(value.(string))
	if err != nil {
		return err
	}

	*adr = Address(*parsed)
	return nil
}

func (g Address) Value() (driver.Value, error) {
	return g.String(), nil
}

func (g Address) String() string {
	adr := mail.Address(g)
	return adr.String()
}

func (u *User) GetLogin() UserLogin {
	return u.Creds.Login
}

func (u *User) GetPassword() string {
	return u.Creds.Password
}
