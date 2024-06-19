package handlers

import (
	"SOA/internal/api"
	"SOA/internal/db"
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

var GetOperation = huma.Operation{
	OperationID:   "getUser",
	Method:        http.MethodGet,
	Path:          "/api/get/{login}",
	Summary:       "Get user",
	Description:   "Get user",
	Tags:          []string{"api", "operations"},
	DefaultStatus: http.StatusOK,
}

type getInput struct {
	Login string `path:"login" json:"login" example:"Lexmach"`
}

type UserOut struct {
	Login string `json:"login" maxLength:"16"`

	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`

	BirthDate string `json:"birthdate"`
	Email     string `json:"email"`

	Phone string `json:"phone"`
}

type getOutput struct {
	Body UserOut
}

type GetHandler struct {
	DB db.ApiDatabase
}

func (uh *GetHandler) Handle(ctx context.Context, input *getInput) (*getOutput, error) {
	user, err := uh.DB.GetUser(api.UserLogin{Login: input.Login})
	if err != nil {
		return nil, huma.Error400BadRequest("User doesnt exist")
	}

	return &getOutput{
		Body: UserOut{
			Login:     user.GetLogin().Login,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			BirthDate: user.BirthDate.Format(YYYYMMDD),
			Email:     user.Email.String(),
			Phone:     user.Phone,
		},
	}, nil
}
