package handlers

import (
	"SOA/cmd/users/auth"
	pb "SOA/proto/api"
	"context"
	"fmt"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

var DeleteOperation = huma.Operation{
	OperationID:   "deletePost",
	Method:        http.MethodDelete,
	Path:          "/posts/delete",
	Summary:       "Delete post",
	Description:   "Delete post",
	Tags:          []string{"posts"},
	DefaultStatus: http.StatusOK,
}

type deleteInput struct {
	Body struct {
		ID       uint64 `json:"id"`
		JWTToken string `json:"token"`
	}
}

type DeleteHandler struct {
	Client    pb.PostsServiceClient
	JWTSecret string
}

func (ch *DeleteHandler) Handle(ctx context.Context, allInput *deleteInput) (*struct{}, error) {
	input := allInput.Body

	login, err := auth.UnmarshalToken(input.JWTToken, ch.JWTSecret)
	if err != nil {
		return &struct{}{}, fmt.Errorf("failed to unmarshal jwt token: %w", err)
	}

	_, err = ch.Client.DeletePost(context.Background(), &pb.DeletePostRequest{
		Owner: login,
		Id:    input.ID,
	})
	if err != nil {
		return nil, err
	}

	return &struct{}{}, nil
}
