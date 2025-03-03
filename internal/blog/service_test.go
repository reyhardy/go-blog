package blog_test

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"github.com/reyhardy/go-blog/internal/blog"
	"github.com/scylladb/gocqlx/v3/qb"
)

func createTable(t *testing.T, keyspace, table string) error {
	err := session.ExecStmt(fmt.Sprintf("CREATE TABLE IF NOT EXISTS %v.%v (id varchar PRIMARY KEY, title varchar, content varchar, author varchar, created_at timestamp, updated_at timestamp);", keyspace, table))
	if err != nil {
		t.Fatal("error create table: ", err)
	}

	return nil
}

func dropTable(t *testing.T, keyspace, table string) error {
	err := session.ExecStmt(fmt.Sprintf("DROP TABLE IF EXISTS %v.%v;", keyspace, table))
	if err != nil {
		t.Fatal("error drop table: ", err)
	}

	return nil
}

func TestAdd(t *testing.T) {
	t.Run("Add Post", func(t *testing.T) {
		postParams := blog.PostParams{
			Title:   "test title",
			Content: "test content",
			Author:  "test author",
		}

		createTable(t, testKeyspace, blog.TablePost)

		addedPosts, err := svc.Add(context.Background(), testKeyspace, &postParams)
		if err != nil {
			t.Errorf("error add: %s", err)
		}

		getPostQuery := session.Query(
			qb.Select(fmt.Sprintf("%v.%v", testKeyspace, blog.TablePost)).
				Columns("title", "content", "author").
				Where(qb.Eq("id")).
				ToCql(),
		).WithContext(context.Background()).BindStruct(blog.Post{ID: postParams.ID})

		post := new(blog.Post)

		if err = getPostQuery.Scan(
			// &post.ID,
			&post.Title,
			&post.Content,
			&post.Author,
		); err != nil {
			t.Errorf("scan query err: %v", err)
		}

		for _, addedPost := range addedPosts {
			if addedPost.Title != postParams.Title {
				t.Errorf("added %v got %v", postParams.Title, addedPost.Title)
			}

			if post.Title != addedPost.Title {
				t.Errorf("got %v want %v", post.Title, addedPost.Title)
			}
		}

		dropTable(t, testKeyspace, blog.TablePost)
	})
}

func TestSelectAll(t *testing.T) {
	t.Run("select all posts", func(t *testing.T) {
		postLength := 10
		// var posts blog.Posts

		createTable(t, testKeyspace, blog.TablePost)

		for i := 0; i < postLength; i++ {
			postToAdd := blog.Post{
				ID:    strconv.Itoa(i),
				Title: fmt.Sprintf("title%v", i),
			}

			err := session.Query(
				qb.Insert(fmt.Sprintf("%v.%v", testKeyspace, blog.TablePost)).
					Columns("id", "title", "content", "author", "created_at", "updated_at").
					ToCql(),
			).WithContext(context.Background()).BindStruct(&postToAdd).ExecRelease()
			if err != nil {
				t.Errorf("error add post: %v", err)
			}
		}

		posts, err := svc.SelectAll(context.Background(), testKeyspace)
		if err != nil {
			t.Errorf("select all error: %v", err)
		}

		if len(posts) != postLength {
			t.Errorf("expected length '%v' but got length '%v'", postLength, len(posts))
		}

		dropTable(t, testKeyspace, blog.TablePost)
	})
}

func TestDelete(t *testing.T) {
	t.Run("Delete Post", func(t *testing.T) {
		tests := []struct {
			name   string
			id     string
			errMsg string
		}{
			{"success delete", "12345678", ""},         //id exist, successful deletion
			{"id not exist", "not exist", "not found"}, //id not exist
		}

		for _, tt := range tests {
			t.Log(tt.name)
			addedPost := blog.Post{
				ID: tt.id,
			}

			createTable(t, testKeyspace, blog.TablePost)

			err := session.Query(
				qb.Insert(fmt.Sprintf("%v.%v", testKeyspace, blog.TablePost)).
					Columns("id", "title", "content", "author").
					ToCql(),
			).WithContext(context.Background()).BindStruct(&addedPost).ExecRelease()
			if err != nil {
				t.Errorf("error add post: %v", err)
			}

			postToDelete := blog.PostParams{
				ID: tt.id,
			}

			_, err = svc.Delete(context.Background(), testKeyspace, &postToDelete)
			if err != nil {
				t.Errorf("delete post error: %v", err)
			}

			getPostQuery := session.Query(
				qb.Select(fmt.Sprintf("%v.%v", testKeyspace, blog.TablePost)).
					Columns("title", "content", "author").
					Where(qb.Eq("id")).
					ToCql(),
			).WithContext(context.Background()).BindStruct(blog.Post{ID: postToDelete.ID})

			post := new(blog.Post)

			err = getPostQuery.Scan(
				&post.ID,
				&post.Title,
				&post.Content,
				&post.Author,
			)

			if tt.errMsg != "" {
				if err == nil {
					t.Errorf("expect error '%v' but got nil error", tt.errMsg)
				} else {
					if err.Error() != tt.errMsg {
						t.Errorf("expect '%v' but got '%v'", tt.errMsg, err.Error())
					}
				}
			}

			if post.ID != "" {
				t.Errorf("expect empty ID but got '%v'", post.ID)
			}
		}

		dropTable(t, testKeyspace, blog.TablePost)
	})
}

func TestUpdate(t *testing.T) {
	t.Run("Update post", func(t *testing.T) {
		tests := []struct {
			name     string
			id       string
			title    string
			expected *blog.Post
			errMsg   string
		}{
			{"update no error", "12345678", "updated title", &blog.Post{Title: "updated title"}, ""}, // no error
			{"id not exist", "not exist", "not exist", nil, "not found"},                             // not found
		}

		for _, tt := range tests {
			t.Log(tt.name)

			createTable(t, testKeyspace, blog.TablePost)

			addedPost := blog.Post{
				ID:    "123456789",
				Title: "test title",
			}

			err := session.Query(
				qb.Insert(fmt.Sprintf("%v.%v", testKeyspace, blog.TablePost)).
					Columns("id", "title", "content", "author").
					ToCql(),
			).WithContext(context.Background()).BindStruct(&addedPost).ExecRelease()
			if err != nil {
				t.Errorf("error add post: %v", err)
			}

			updateParams := blog.PostParams{
				ID:    tt.id,
				Title: tt.title,
			}

			_, err = svc.Update(context.Background(), testKeyspace, &updateParams)
			if err != nil {
				t.Errorf("update post error: %v", err)
			}

			getPostQuery := session.Query(
				qb.Select(fmt.Sprintf("%v.%v", testKeyspace, blog.TablePost)).
					Columns("title", "content", "author").
					Where(qb.Eq("id")).
					ToCql(),
			).WithContext(context.Background()).BindStruct(blog.Post{ID: updateParams.ID})

			post := new(blog.Post)

			err = getPostQuery.Scan(
				// &post.ID,
				&post.Title,
				&post.Content,
				&post.Author,
			)

			if tt.errMsg == "" {
				if err != nil {
					t.Errorf("expect nil error got '%v'", err.Error())
				}
			} else {
				if err != nil {
					if tt.errMsg != err.Error() {
						t.Errorf("expected %v got %v", tt.errMsg, err.Error())
					}
				}
			}

			if tt.expected != nil {
				if post.Title != tt.expected.Title {
					t.Errorf("expected %v got %v", tt.expected.Title, post.Title)
				}
			}

			dropTable(t, testKeyspace, blog.TablePost)
		}
	})
}
