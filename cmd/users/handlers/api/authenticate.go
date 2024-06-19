package handlers

import (
	"SOA/cmd/users/auth"
	"SOA/internal/api"
	"SOA/internal/db"
	"context"
	"fmt"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

var AuthOperation = huma.Operation{
	OperationID:   "authenticateUser",
	Method:        http.MethodPost,
	Path:          "/api/auth",
	Summary:       "Auth user",
	Description:   "Auth user with response of user token",
	Tags:          []string{"api", "auth"},
	DefaultStatus: http.StatusOK,
}

type authInput struct {
	Body struct {
		Login    string `json:"login" maxLength:"16" example:"Lexmach"`
		Password string `json:"password" example:"mycoolpass"`
	}
}

type authOutput struct {
	Body api.UserToken
}

type AuthHandler struct {
	DB        db.ApiDatabase
	JWTSecret string
}

func (ah *AuthHandler) Handle(ctx context.Context, allInput *authInput) (*authOutput, error) {
	input := allInput.Body
	creds := api.UserCredentials{
		Login: api.UserLogin{
			Login: input.Login,
		},
		Password: input.Password,
	}
	user, err := ah.DB.GetUser(creds.Login)
	if err != nil {
		return nil, err
	}
	if auth.CompareSaltedAndOrigin(user.GetPassword(), creds.Password) {
		return nil, fmt.Errorf("wrong password")
	}

	return &authOutput{auth.MustCreateToken(user.GetLogin(), ah.JWTSecret)}, nil
}
