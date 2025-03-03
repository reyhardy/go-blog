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
	Add(ctx context.Context, keyspace string, postParams *PostParams) (Posts, error)
	Get(ctx context.Context, keyspace string) (*Post, error)
	SelectAll(ctx context.Context, keyspace string) (Posts, error)
	Update(ctx context.Context, keyspace string, postParams *PostParams) (Posts, error)
	Delete(ctx context.Context, keyspace string, postParams *PostParams) (Posts, error)
}

func NewService(session scylladb.Client) Servicer {
	return &service{session}
}

var postList Posts

func (s *service) Add(ctx context.Context, keyspace string, postParams *PostParams) (Posts, error) {
	q := qb.Insert(fmt.Sprintf("%s.%s", keyspace, TablePost)).Columns("id", "title", "content", "author", "created_at", "updated_at")

	post := NewPost(postParams)

	err := s.db.QueryExec(ctx, q, post)
	if err != nil {
		return nil, err
	}

	postList = append(postList, post)

	return postList, nil
}

func (s *service) Get(ctx context.Context, keyspace string) (*Post, error) {
	return nil, nil
}

func (s *service) SelectAll(ctx context.Context, keyspace string) (Posts, error) {
	q := qb.Select(fmt.Sprintf("%s.%s", keyspace, TablePost)).Columns("id", "title", "content", "author", "created_at", "updated_at")
	iter, err := s.db.QueryRow(ctx, q)
	if err != nil {
		return nil, err
	}
	defer iter.Close()

	if err = iter.Select(&postList); err != nil {
		return nil, err
	}

	return postList, nil
}

func (s *service) Update(ctx context.Context, keyspace string, postParams *PostParams) (Posts, error) {
	q := qb.Update(fmt.Sprintf("%s.%s", keyspace, TablePost)).Set("title", "content", "author", "updated_at").Where(qb.Eq("id"))

	updatedPost := MapPost(postParams, nil)

	err := s.db.QueryExec(ctx, q, updatedPost)
	if err != nil {
		return nil, err
	}

	for i, post := range postList {
		if post.ID == updatedPost.ID {
			postList[i] = updatedPost
		}
	}

	return postList, nil
}

func (s *service) Delete(ctx context.Context, keyspace string, postParams *PostParams) (Posts, error) {
	q := qb.Delete(fmt.Sprintf("%s.%s", keyspace, TablePost)).Where(qb.Eq("id"))

	deletedPost := MapPost(postParams, nil)

	err := s.db.QueryExec(ctx, q, deletedPost)
	if err != nil {
		return nil, err
	}

	var filteredPost Posts

	for _, post := range postList {
		if post.ID != deletedPost.ID {
			filteredPost = append(filteredPost, post)
		}
	}

	postList = filteredPost

	return postList, nil
}
