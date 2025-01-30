package blog

import (
	"time"

	"github.com/segmentio/ksuid"
)

const (
	Keyspace  string = "blog"
	TablePost string = "post"
)

// var PostMetadata = table.Metadata{
// 	Name: "post",
// 	Columns: []string{"id", "title", "content", "author", "created_at", "updated_at"},
// }

type Post struct {
	ID      string `db:"id"`
	Title   string `db:"title"`
	Content string `db:"content"`
	Author  string `db:"author"`

	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type Posts []*Post

func MapPost(postParams *PostParams, createdAt *time.Time, updatedAt *time.Time) *Post {
	post := &Post{
		ID:      postParams.ID,
		Title:   postParams.Title,
		Content: postParams.Content,
		Author:  postParams.Author,
	}

	if createdAt != nil {
		post.CreatedAt = *createdAt
	}

	if updatedAt != nil {
		post.UpdatedAt = *updatedAt
	}

	return post
}

func NewPost(postParams *PostParams) *Post {
	ksuid := ksuid.New()

	createdAt := time.Now()

	postParams.ID = ksuid.String()

	return MapPost(postParams, &createdAt, nil)
}

type PostParams struct {
	ID      string `db:"id"`
	Title   string `db:"title"`
	Content string `db:"content"`
	Author  string `db:"author"`
}
