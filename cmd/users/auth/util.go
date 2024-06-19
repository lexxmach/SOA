package auth

import (
	"SOA/internal/api"
	"SOA/internal/common"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func MustSaltPassword(pass string) string {
	return string(common.Must(bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)))
}

// Returns true if equal
func CompareSaltedAndOrigin(salted, origin string) bool {
	return bcrypt.CompareHashAndPassword([]byte(salted), []byte(origin)) != nil
}

func MustCreateToken(login api.UserLogin, jwtSecret string) api.UserToken {
	payload := jwt.MapClaims{
		"login": login.Login,
		"exp":   time.Now().Add(time.Hour * 72).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)

	return api.UserToken{Token: common.Must(token.SignedString([]byte(jwtSecret)))}
}

func UnmarshalToken(jwtToken string, jwtSecret string) (string, error) {
	token, err := jwt.Parse(jwtToken, func(t *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})

	if err != nil {
		return "", err
	}
	mapClaims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", err
	}
	login, ok := mapClaims["login"].(string)
	if !ok {
		return "", err
	}

	return login, nil
}
