package handlers

import (
	"SOA/internal/posts"
	pb "SOA/proto/api"
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

var GetPostOperation = huma.Operation{
	OperationID:   "getPost",
	Method:        http.MethodPost,
	Path:          "/posts/get",
	Summary:       "Get post",
	Description:   "Get post",
	Tags:          []string{"posts"},
	DefaultStatus: http.StatusOK,
}

type getPostInput struct {
	Body struct {
		ID uint64 `json:"id"`
	}
}

type getPostOutput struct {
	Body *posts.Post
}

type GetPostHandler struct {
	Client    pb.PostsServiceClient
	JWTSecret string
}

func (ch *GetPostHandler) Handle(ctx context.Context, allInput *getPostInput) (*getPostOutput, error) {
	input := allInput.Body

	out, err := ch.Client.GetPost(context.Background(), &pb.GetPostRequest{
		Id: input.ID,
	})
	if err != nil {
		return nil, err
	}

	return &getPostOutput{Body: &posts.Post{
		ID:    out.Id,
		Owner: out.Owner,
		Body:  out.Body,
	}}, nil
}
