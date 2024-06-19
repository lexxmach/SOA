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
	Path:          "/posts/get/{id}",
	Summary:       "Get post",
	Description:   "Get post",
	Tags:          []string{"posts"},
	DefaultStatus: http.StatusOK,
}

type getPostInput struct {
	ID uint64 `path:"id" json:"id"`
}

type getPostOutput struct {
	Body *posts.Post
}

type GetPostHandler struct {
	Client    pb.PostsServiceClient
	JWTSecret string
}

func (ch *GetPostHandler) Handle(ctx context.Context, input *getPostInput) (*getPostOutput, error) {
	out, err := ch.Client.GetPost(ctx, &pb.GetPostRequest{
		Id: input.ID,
	})
	if err != nil {
		return nil, huma.Error400BadRequest("No such post exist")
	}

	return &getPostOutput{Body: &posts.Post{
		ID:    out.Id,
		Owner: out.Owner,
		Body:  out.Body,
	}}, nil
}
