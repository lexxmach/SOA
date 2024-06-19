package handlers

import (
	"SOA/cmd/users/auth"
	pb "SOA/proto/api"
	"context"
	"fmt"
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
		return nil, fmt.Errorf("failed to unmarshal jwt token: %w", err)
	}

	_, err = ch.Client.UpdatePost(context.Background(), &pb.UpdatePostRequest{
		Owner: login,
		Id:    input.ID,
		Body:  input.Body,
	})
	if err != nil {
		return nil, err
	}

	return nil, nil
}
