package blog

import (
	"context"
	"fmt"
	"time"

	"github.com/reyhardy/go-blog/pkg/pgsql"
)

type service struct {
	db pgsql.Client
}

type servicer interface {
	Add(ctx context.Context, postParams *PostParams) error
	Get(ctx context.Context) (*Post, error)
	SelectAll(ctx context.Context) (Posts, error)
	Update(ctx context.Context, postParams *PostParams) (*Post, error)
	Delete(ctx context.Context, postParams *PostParams) error
}

func newService(db pgsql.Client) servicer {
	return &service{db}
}

func (s *service) Add(ctx context.Context, postParams *PostParams) error {
	q := fmt.Sprintf(
		"INSERT INTO %s (id, title, content, author, created_at, updated_at) VALUES($1, $2, $3, $4, $5, $6);",
		TablePost,
	)

	post := NewPost(postParams)

	err := s.db.Exec(ctx, q,
		post.ID,
		post.Title,
		post.Content,
		post.Author,
		post.CreatedAt,
		post.UpdatedAt,
	)
	if err != nil {
		return err
	}

	fmt.Printf("new post: %+v\n", post)

	return nil
}

func (s *service) Get(ctx context.Context) (*Post, error) {
	return nil, nil
}

func (s *service) SelectAll(ctx context.Context) (Posts, error) {
	q := fmt.Sprintf("SELECT id, title, content, author, created_at, updated_at FROM %s;", TablePost)

	var posts Posts
	err := s.db.Query(ctx, q, &posts)
	if err != nil {
		return nil, err
	}

	return posts, nil
}

func (s *service) Update(ctx context.Context, postParams *PostParams) (*Post, error) {
	q := fmt.Sprintf("UPDATE %s SET title = $1, content = $2, author = $3, updated_at = $4 WHERE id = $5;", TablePost)

	updatedAt := time.Now()

	post := MapPost(postParams, nil, &updatedAt)

	if err := s.db.Exec(ctx, q,
		post.Title,
		post.Content,
		post.Author,
		post.UpdatedAt,
		post.ID,
	); err != nil {
		return nil, err
	}

	return post, nil
}

func (s *service) Delete(ctx context.Context, postParams *PostParams) error {
	q := fmt.Sprintf("DELETE FROM %s WHERE id = $1;", TablePost)

	return s.db.Exec(ctx, q, postParams.ID)
}
