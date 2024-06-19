package handlers

import (
	"SOA/internal/posts"
	pb "SOA/proto/api"
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

var ListPostOperation = huma.Operation{
	OperationID:   "listPost",
	Method:        http.MethodPost,
	Path:          "/posts/list",
	Summary:       "List post",
	Description:   "List post",
	Tags:          []string{"posts"},
	DefaultStatus: http.StatusOK,
}

type listPostInput struct {
	Body struct {
		PageNum  uint64 `json:"page_num"`
		PageSize uint64 `json:"page_size"`
	}
}

type listPostOutput struct {
	Body []*posts.Post
}

type ListPostHandler struct {
	Client pb.PostsServiceClient
}

func (ch *ListPostHandler) Handle(ctx context.Context, allInput *listPostInput) (*listPostOutput, error) {
	input := allInput.Body

	out, err := ch.Client.ListPosts(context.Background(), &pb.ListPostsRequest{
		PageNum:  input.PageNum,
		PageSize: input.PageSize,
	})
	if err != nil {
		return nil, err
	}

	boiler := make([]*posts.Post, 0, len(out.Posts))
	for _, post := range out.Posts {
		boiler = append(boiler, &posts.Post{
			ID:    post.Id,
			Owner: post.Owner,
			Body:  post.Body,
		})
	}
	return &listPostOutput{Body: boiler}, nil
}
