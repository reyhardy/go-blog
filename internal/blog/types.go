package blog

import (
	"time"

	"github.com/segmentio/ksuid"
)

const (
	Keyspace  string = "blog"
	TablePost string = "post"
)

type Post struct {
	ID      string `db:"id"`
	Title   string `db:"title"`
	Content string `db:"content"`
	Author  string `db:"author"`

	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type Posts []*Post

func MapPost(postParams *PostParams, createdAt *time.Time) *Post {
	post := &Post{
		ID:      postParams.ID,
		Title:   postParams.Title,
		Content: postParams.Content,
		Author:  postParams.Author,
	}

	if createdAt != nil {
		post.CreatedAt = *createdAt
		post.UpdatedAt = *createdAt
	} else {
		post.UpdatedAt = time.Now()
	}

	return post
}

func NewPost(postParams *PostParams) *Post {
	ksuid := ksuid.New()

	createdAt := time.Now()

	postParams.ID = ksuid.String()

	return MapPost(postParams, &createdAt)
}

type PostParams struct {
	ID      string `db:"id"`
	Title   string `db:"title"`
	Content string `db:"content"`
	Author  string `db:"author"`
}
