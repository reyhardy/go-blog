package blog

import (
	"time"

	"github.com/segmentio/ksuid"
)

const (
	TablePost string = "post"
)

type Post struct {
	ID        string `json:"id" db:"id"`
	Title     string `json:"title" db:"title"`
	Content   string `json:"content" db:"content"`
	Author    string `json:"author" db:"author"`
	IsDeleted bool   `json:"deleted" db:"deleted"`

	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
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
	} else {
		post.UpdatedAt = *createdAt
	}

	return post
}

func NewPost(postParams *PostParams) *Post {
	ksuid := ksuid.New()

	createdAt := time.Now()

	postParams.ID = ksuid.String()

	return MapPost(postParams, &createdAt, nil)
}

type CustomTime struct {
	time.Time
}

type PostParams struct {
	ID      string `json:"id" db:"id"`
	Title   string `json:"title" db:"title"`
	Content string `json:"content" db:"content"`
	Author  string `json:"author" db:"author"`
}
