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
	Method:        http.MethodGet,
	Path:          "/posts/list",
	Summary:       "List post",
	Description:   "List post",
	Tags:          []string{"posts"},
	DefaultStatus: http.StatusOK,
}

type listPostInput struct {
	PageNum  uint64 `query:"pagenum" json:"page_num"`
	PageSize uint64 `query:"pagesize" json:"page_size"`
}

type listPostOutput struct {
	Body []*posts.Post
}

type ListPostHandler struct {
	Client pb.PostsServiceClient
}

func (ch *ListPostHandler) Handle(ctx context.Context, input *listPostInput) (*listPostOutput, error) {
	out, err := ch.Client.ListPosts(ctx, &pb.ListPostsRequest{
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
