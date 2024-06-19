package handlers

import (
	"SOA/internal/db"
	"SOA/internal/posts"
	"SOA/proto/api"
	proto "SOA/proto/posts"
	"context"
	"fmt"
)

type PostGRPCHandler struct {
	api.UnimplementedPostsServiceServer
	DB db.PostDatabase
}

func (p PostGRPCHandler) CreatePost(ctx context.Context, req *api.CreatePostRequest) (*api.CreatePostResponse, error) {
	id, err := p.DB.CreatePost(posts.Post{
		Owner: req.Owner,
		Body:  req.Body,
	})
	if err != nil {
		return nil, err
	}
	return &api.CreatePostResponse{
		Id: id,
	}, nil
}

func (p PostGRPCHandler) DeletePost(ctx context.Context, req *api.DeletePostRequest) (*api.DeletePostResponse, error) {
	post, err := p.DB.GetPost(req.Id)
	if err != nil {
		return nil, err
	}
	if post.Owner != req.Owner {
		return nil, fmt.Errorf("permission denied: post owner is not request owner")
	}

	err = p.DB.DeletePost(req.Id)
	if err != nil {
		return nil, err
	}
	return &api.DeletePostResponse{}, nil
}

func (p PostGRPCHandler) GetPost(ctx context.Context, req *api.GetPostRequest) (*proto.Post, error) {
	post, err := p.DB.GetPost(req.GetId())
	if err != nil {
		return nil, err
	}
	return &proto.Post{
		Id:    post.ID,
		Owner: post.Owner,
		Body:  post.Body,
	}, nil
}

func (p PostGRPCHandler) ListPosts(ctx context.Context, req *api.ListPostsRequest) (*api.ListPostsResponse, error) {
	posts, err := p.DB.ListPosts(req.PageNum, req.PageSize)
	if err != nil {
		return nil, err
	}
	resp := &api.ListPostsResponse{}
	for _, post := range posts {
		resp.Posts = append(resp.Posts, &proto.Post{
			Id:    post.ID,
			Owner: post.Owner,
			Body:  post.Body,
		})
	}
	return resp, nil
}

// UpdatePost implements api.PostsServiceServer.
func (p PostGRPCHandler) UpdatePost(ctx context.Context, req *api.UpdatePostRequest) (*api.UpdatePostResponse, error) {
	post, err := p.DB.GetPost(req.Id)
	if err != nil {
		return nil, err
	}
	if post.Owner != req.Owner {
		return nil, fmt.Errorf("permission denied: post owner is not request owner")
	}

	err = p.DB.UpdatePost(&posts.Post{
		ID:    req.Id,
		Owner: req.Owner,
		Body:  req.Body,
	})
	if err != nil {
		return nil, err
	}

	return &api.UpdatePostResponse{}, nil
}
