package handlers

import (
	"SOA/cmd/api/auth"
	"SOA/internal/api"
	"SOA/internal/common"
	"SOA/internal/db"
	"context"
	"fmt"
	"net/http"
	"net/mail"
	"time"

	"github.com/danielgtaylor/huma/v2"
)

var UpdateOperation = huma.Operation{
	OperationID:   "update",
	Method:        http.MethodPost,
	Path:          "/update",
	Summary:       "Update user",
	Description:   "Update user with response of user token",
	Tags:          []string{"operations"},
	DefaultStatus: http.StatusOK,
}

type updateInput struct {
	Body struct {
		JWTToken string `json:"token"`

		FirstName string `json:"firstname" maxLength:"16" example:"Konstantin"`
		LastName  string `json:"lastname" maxLength:"16" example:"Frolov"`

		BirthDate string `json:"birthdate" example:"2003-09-09"`
		Email     string `json:"email" maxLength:"32" example:"someone@gmail.com"`

		Phone string `json:"phone" maxLength:"16" example:"+78005553535"`
	}
}

// TODO(lexmach): remove copy-paste
func (i *updateInput) Resolve(ctx huma.Context) []error {
	var ers []error

	_, err := time.Parse(YYYYMMDD, i.Body.BirthDate)
	if err != nil {
		ers = append(ers, &huma.ErrorDetail{
			Location: "path.birthdate",
			Message:  "BirthDate is in incorrect format",
			Value:    err.Error(),
		})
	}

	_, err = mail.ParseAddress(i.Body.Email)
	if err != nil {
		ers = append(ers, &huma.ErrorDetail{
			Location: "path.email",
			Message:  "Email is in incorrect format",
			Value:    err.Error(),
		})
	}

	return ers
}

type UpdateHandler struct {
	DB        db.Database
	JWTSecret string
}

func (uh *UpdateHandler) Handle(ctx context.Context, allInput *updateInput) (*struct{}, error) {
	input := allInput.Body

	login, err := auth.UnmarshalToken(input.JWTToken, uh.JWTSecret)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal jwt token: %w", err)
	}

	// should never fail cause of validation
	user := &api.User{
		FirstName: input.FirstName,
		LastName:  input.LastName,
		BirthDate: common.Must(time.Parse(YYYYMMDD, input.BirthDate)),
		Email:     *common.Must(mail.ParseAddress(input.Email)),
		Phone:     input.Phone,

		Creds: api.UserCredentials{
			Login: api.UserLogin{
				Login: login,
			},
		},
	}

	err = uh.DB.UpdateUser(user)
	if err != nil {
		return nil, err
	}

	return &struct{}{}, nil
}
