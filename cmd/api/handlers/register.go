package handlers

import (
	"SOA/cmd/api/auth"
	"SOA/internal/api"
	"SOA/internal/common"
	"SOA/internal/db"
	"context"
	"net/http"
	"net/mail"
	"time"

	"github.com/danielgtaylor/huma/v2"
)

var RegisterOperation = huma.Operation{
	OperationID:   "register",
	Method:        http.MethodPost,
	Path:          "/register",
	Summary:       "Register user",
	Description:   "Register user with response of user token",
	Tags:          []string{"auth"},
	DefaultStatus: http.StatusCreated,
}

const YYYYMMDD = "2006-01-02"

type registerInput struct {
	Body struct {
		Login    string `json:"login" maxLength:"16" example:"Lexmach"`
		Password string `json:"password" example:"mycoolpass"`

		FirstName string `json:"firstname" maxLength:"16" example:"Konstantin"`
		LastName  string `json:"lastname" maxLength:"16" example:"Frolov"`

		BirthDate string `json:"birthdate" example:"2003-09-09"`
		Email     string `json:"email" maxLength:"32" example:"someone@gmail.com"`

		Phone string `json:"phone" maxLength:"16" example:"+78005553535"`
	}
}

func (i *registerInput) Resolve(ctx huma.Context) []error {
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

type registerOutput struct {
	Body api.UserToken
}

type RegisterHandler struct {
	DB        db.Database
	JWTSecret string
}

func (rh *RegisterHandler) Handle(ctx context.Context, allInput *registerInput) (*registerOutput, error) {
	// should never fail cause of validation
	input := allInput.Body
	user := api.User{
		FirstName: input.FirstName,
		LastName:  input.LastName,
		BirthDate: common.Must(time.Parse(YYYYMMDD, input.BirthDate)),
		Email:     *common.Must(mail.ParseAddress(input.Email)),
		Phone:     input.Phone,

		Creds: api.UserCredentials{
			Login: api.UserLogin{
				Login: input.Login,
			},
			Password: auth.MustSaltPassword(input.Password),
		},
	}

	err := rh.DB.CreateUser(user)
	if err != nil {
		return nil, err
	}

	return &registerOutput{auth.MustCreateToken(user.GetLogin(), rh.JWTSecret)}, nil
}
