package handlers

import (
	"SOA/cmd/users/auth"
	pb "SOA/proto/api"
	"context"
	"fmt"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

var CreateOperation = huma.Operation{
	OperationID:   "createPost",
	Method:        http.MethodPost,
	Path:          "/posts/create",
	Summary:       "Create post",
	Description:   "Create post",
	Tags:          []string{"posts"},
	DefaultStatus: http.StatusCreated,
}

type createInput struct {
	Body struct {
		Body     string `json:"string" example:"my first post!11"`
		JWTToken string `json:"token"`
	}
}

type createOutput struct {
	Body *pb.CreatePostResponse
}

type CreateHandler struct {
	Client    pb.PostsServiceClient
	JWTSecret string
}

func (ch *CreateHandler) Handle(ctx context.Context, allInput *createInput) (*createOutput, error) {
	input := allInput.Body

	login, err := auth.UnmarshalToken(input.JWTToken, ch.JWTSecret)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal jwt token: %w", err)
	}

	out, err := ch.Client.CreatePost(context.Background(), &pb.CreatePostRequest{
		Owner: login,
		Body:  input.Body,
	})
	if err != nil {
		return nil, err
	}

	return &createOutput{Body: out}, nil
}
