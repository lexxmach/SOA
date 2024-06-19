package handlers

import (
	"SOA/cmd/users/auth"
	pb "SOA/proto/api"
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

var UpdatePostOperation = huma.Operation{
	OperationID:   "updatePost",
	Method:        http.MethodPost,
	Path:          "/posts/update",
	Summary:       "Update post",
	Description:   "Update post",
	Tags:          []string{"posts"},
	DefaultStatus: http.StatusOK,
}

type updatePostInput struct {
	Body struct {
		ID       uint64 `json:"id"`
		Body     string `json:"body"`
		JWTToken string `json:"token"`
	}
}

type UpdatePostHandler struct {
	Client    pb.PostsServiceClient
	JWTSecret string
}

func (ch *UpdatePostHandler) Handle(ctx context.Context, allInput *updatePostInput) (*struct{}, error) {
	input := allInput.Body

	login, err := auth.UnmarshalToken(input.JWTToken, ch.JWTSecret)
	if err != nil {
		return nil, huma.Error401Unauthorized("Wrong token")
	}

	_, err = ch.Client.UpdatePost(ctx, &pb.UpdatePostRequest{
		Owner: login,
		Id:    input.ID,
		Body:  input.Body,
	})
	if err != nil {
		return nil, err
	}

	return nil, nil
}
