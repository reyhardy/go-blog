package blog

import (
	"context"
	"fmt"

	"github.com/reyhardy/go-blog/db/scylladb"
	"github.com/scylladb/gocqlx/v3/qb"
)

type service struct {
	db scylladb.Client
}

type Servicer interface {
	Add(ctx context.Context, keyspace string, postParams *PostParams) (*Post, error)
	Get(ctx context.Context, keyspace string) (*Post, error)
	SelectAll(ctx context.Context, keyspace string) (Posts, error)
	Update() (*Post, error)
	Delete(ctx context.Context, keyspace string, postParams *PostParams) error
}

func NewService(session scylladb.Client) Servicer {
	return &service{session}
}

func (s *service) Add(ctx context.Context, keyspace string, postParams *PostParams) (*Post, error) {
	q := qb.Insert(fmt.Sprintf("%s.%s", keyspace, TablePost)).Columns("id", "title", "content", "author", "created_at", "updated_at")
	// q := qb.Insert(TablePost).Columns("id", "title", "content", "author", "created_at", "updated_at")

	post := NewPost(postParams)

	err := s.db.QueryExec(ctx, q, post)
	if err != nil {
		return nil, err
	}

	return post, nil
}

func (s *service) Get(ctx context.Context, keyspace string) (*Post, error) {
	return nil, nil
}

func (s *service) SelectAll(ctx context.Context, keyspace string) (Posts, error) {
	q := qb.Select(fmt.Sprintf("%s.%s", keyspace, TablePost)).Columns("id", "title", "content", "author")
	// q := qb.Select(TablePost).Columns("id", "title", "content", "author")
	iter, err := s.db.QueryRow(ctx, q)
	if err != nil {
		return nil, err
	}
	defer iter.Close()

	var postList Posts

	if err = iter.Select(&postList); err != nil {
		return nil, err
	}

	return postList, nil
}

func (s *service) Update() (*Post, error) {
	return nil, nil
}

func (s *service) Delete(ctx context.Context, keyspace string, postParams *PostParams) error {
	q := qb.Delete(fmt.Sprintf("%s.%s", keyspace, TablePost)).Where(qb.Eq("id"))
	return s.db.QueryExec(ctx, q, postParams)
}
