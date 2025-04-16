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
	Get(ctx context.Context, keyspace string, postParams *PostParams) (*Post, error)
	SelectAll(ctx context.Context, keyspace string) (Posts, error)
	Update(ctx context.Context, keyspace string, postParams *PostParams) (*Post, error)
	Delete(ctx context.Context, keyspace string, postParams *PostParams) error
}

func NewService(session scylladb.Client) Servicer {
	return &service{session}
}

var posts Posts

func (s *service) Add(ctx context.Context, keyspace string, postParams *PostParams) (*Post, error) {
	q := qb.Insert(fmt.Sprintf("%s.%s", keyspace, TablePost)).Columns("id", "title", "content", "author", "created_at", "updated_at")

	post := NewPost(postParams)

	_, err := s.db.QueryExec(ctx, q, post)
	if err != nil {
		return nil, err
	}

	posts = append(posts, post)

	return post, nil
}

func (s *service) Get(ctx context.Context, keyspace string, postParams *PostParams) (*Post, error) {
	q := qb.Select(fmt.Sprintf("%s.%s", keyspace, TablePost)).Columns("id", "title", "content", "author", "created_at", "updated_at").Where(qb.Eq("id"))

	getPost := MapPost(postParams, nil)

	stmt, _ := q.ToCql()

	query, err := s.db.QueryExec(ctx, q, getPost)
	if err != nil {
		return nil, err
	}

	fmt.Printf("query: %s;\npost: %+v\nquery struct: %+v", stmt, getPost, query)

	var post *Post

	if err = query.SelectRelease(post); err != nil {
		return nil, err
	}

	return post, nil
}

func (s *service) SelectAll(ctx context.Context, keyspace string) (Posts, error) {
	q := qb.Select(fmt.Sprintf("%s.%s", keyspace, TablePost)).Columns("id", "title", "content", "author", "created_at", "updated_at")
	iter, err := s.db.QueryRow(ctx, q)
	if err != nil {
		return nil, err
	}
	defer iter.Close()

	if err = iter.Select(&posts); err != nil {
		return nil, err
	}

	return posts, nil
}

func (s *service) Update(ctx context.Context, keyspace string, postParams *PostParams) (*Post, error) {
	q := qb.Update(fmt.Sprintf("%s.%s", keyspace, TablePost)).Set("title", "content", "author", "updated_at").Where(qb.Eq("id"))

	updatedPost := MapPost(postParams, nil)

	_, err := s.db.QueryExec(ctx, q, updatedPost)
	if err != nil {
		return nil, err
	}

	for i, post := range posts {
		if post.ID == updatedPost.ID {
			posts[i] = updatedPost
		}
	}
	return updatedPost, nil
}

func (s *service) Delete(ctx context.Context, keyspace string, postParams *PostParams) error {
	q := qb.Delete(fmt.Sprintf("%s.%s", keyspace, TablePost)).Where(qb.Eq("id"))

	deletedPost := MapPost(postParams, nil)

	// var filteredPost Posts

	// for _, post := range posts {
	// 	if post.ID != deletedPost.ID {
	// 		filteredPost = append(filteredPost, post)
	// 	}
	// }

	// posts = filteredPost

	_, err := s.db.QueryExec(ctx, q, deletedPost)
	if err != nil {
		return err
	}

	return nil
}
