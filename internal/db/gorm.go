package db

import (
	"SOA/internal/api"
	"SOA/internal/posts"
	"fmt"

	"gorm.io/gorm"
)

type ApiGormDatabase struct {
	db *gorm.DB
}

func CreateApiGormDB(dialector gorm.Dialector) (ApiDatabase, error) {
	gormDB, err := gorm.Open(dialector)
	if err != nil {
		return nil, fmt.Errorf("failed to setup gorm: %w", err)
	}

	gormDB.AutoMigrate(&api.User{})

	return &ApiGormDatabase{
		db: gormDB,
	}, nil
}

func (gorm *ApiGormDatabase) CreateUser(user api.User) error {
	tx := gorm.db.Create(&user)
	if tx.Error != nil {
		return fmt.Errorf("failed to create user %q: %w", user, tx.Error)
	}

	return nil
}

func (gorm *ApiGormDatabase) GetUser(login api.UserLogin) (*api.User, error) {
	user := &api.User{}

	tx := gorm.db.First(&user, login)
	if tx.Error != nil {
		return nil, fmt.Errorf("failed to get user %q: %w", login, tx.Error)
	}

	return user, nil
}

func (gorm *ApiGormDatabase) UpdateUser(user *api.User) error {
	tx := gorm.db.Save(user)
	if tx.Error != nil {
		return fmt.Errorf("failed to update user %q: %w", user, tx.Error)
	}

	return nil
}

type PostsGormDatabase struct {
	db *gorm.DB
}

func (p *PostsGormDatabase) CreatePost(post posts.Post) (uint64, error) {
	tx := p.db.Create(&post)
	if tx.Error != nil {
		return 0, fmt.Errorf("failed to create post: %w", tx.Error)
	}
	return post.ID, nil
}

func (p *PostsGormDatabase) DeletePost(id uint64) error {
	tx := p.db.Delete(&posts.Post{}, &id)
	if tx.Error != nil {
		return fmt.Errorf("failed to delete post: %w", tx.Error)
	}
	return nil
}

func (p *PostsGormDatabase) GetPost(id uint64) (*posts.Post, error) {
	post := posts.Post{}
	tx := p.db.First(&post, id)
	if tx.Error != nil {
		return nil, fmt.Errorf("failed to get post %d: %w", id, tx.Error)
	}
	return &post, nil
}

func (p *PostsGormDatabase) ListPosts(pageNum uint64, pageSize uint64) ([]posts.Post, error) {
	posts := []posts.Post{}
	tx := p.db.Limit(int(pageSize)).Offset(int(pageNum * pageSize)).Find(&posts)
	if tx.Error != nil {
		return nil, fmt.Errorf("failed to get posts: %w", tx.Error)
	}
	return posts, nil
}

func (p *PostsGormDatabase) UpdatePost(post *posts.Post) error {
	tx := p.db.Save(&post)
	if tx.Error != nil {
		return fmt.Errorf("failed to update post: %w", tx.Error)
	}
	return nil
}

func CreatePostsGormDB(dialector gorm.Dialector) (PostDatabase, error) {
	gormDB, err := gorm.Open(dialector)
	if err != nil {
		return nil, fmt.Errorf("failed to setup gorm: %w", err)
	}

	gormDB.AutoMigrate(&posts.Post{})

	return &PostsGormDatabase{
		db: gormDB,
	}, nil
}
