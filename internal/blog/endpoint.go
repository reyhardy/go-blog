package blog

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/reyhardy/go-blog/pkg/pgsql"
	"github.com/starfederation/datastar-go/datastar"
	"maragu.dev/gomponents"
)

// const (
// 	datastarSelector string = "datastar-selector"
// 	datastarMode     string = "datastar-mode"
// )

type endpoint struct {
	svc      servicer
	listener listener
}

type API interface {
	GetHome(c echo.Context) error
	GetPostSSE(c echo.Context) error
	GetPost(c echo.Context) error
	AddPost(c echo.Context) error
	DeletePost(c echo.Context) error
	UpdatePost(c echo.Context) error
}

func NewAPI(db pgsql.Client) API {
	return &endpoint{
		svc:      newService(db),
		listener: NewListener(db),
	}
}

func (e *endpoint) GetHome(c echo.Context) error {
	res, err := e.svc.SelectAll(c.Request().Context())
	if err != nil {
		c.Response().WriteHeader(echo.ErrInternalServerError.Code)
		fmt.Fprintf(c.Response().Writer, "Error fetching posts: %v", err)
	}

	return Home(res).Render(c.Response().Writer)
}

func (e *endpoint) GetPost(c echo.Context) error {
	notification := make(chan string)

	sse := datastar.NewSSE(c.Response().Writer, c.Request())

	go e.listener.Listen(c.Request().Context(), notification)

	select {
	case <-c.Request().Context().Done():
		return c.Request().Context().Err()
	case payload := <-notification:
		fmt.Println("payload received: ", payload)

		res, err := e.svc.SelectAll(c.Request().Context())
		if err != nil {
			c.Response().WriteHeader(echo.ErrInternalServerError.Code)
			fmt.Fprintf(c.Response().Writer, "Error fetching posts: %v", err)
		}

		return sse.PatchElements(
			gomponents.NodeFunc(PostList(res).Render).String(),
			datastar.WithModeReplace(),
			datastar.WithSelectorID("posts"),
		)
	}

	// return c.Stream(http.StatusOK, )
}

func (e *endpoint) GetPostSSE(c echo.Context) error {
	notiChan := make(chan Post)

	sse := datastar.NewSSE(c.Response().Writer, c.Request())

	go e.listener.ListenForNotification(c.Request().Context(), "db_changes", notiChan)

	for res := range notiChan {
		select {
		case <-c.Request().Context().Done():
			c.Response().WriteHeader(echo.ErrInternalServerError.Code)
			fmt.Fprintf(c.Response().Writer, "Error streaming posts: %v", c.Request().Context().Err())
		default:
			switch true {
			case res.IsDeleted:
				sse.PatchElements(
					"",
					datastar.WithModeRemove(),
					datastar.WithSelectorID(fmt.Sprintf("post-%s", res.ID)),
				)
			case !res.IsDeleted && res.UpdatedAt.Equal(res.CreatedAt):
				sse.PatchElements(
					gomponents.NodeFunc(PostCard(&res).Render).String(),
					datastar.WithModeAppend(),
					datastar.WithSelectorID("posts"),
				)
			case !res.IsDeleted && !res.UpdatedAt.Equal(res.CreatedAt):
				sse.PatchElements(
					gomponents.NodeFunc(PostCard(&res).Render).String(),
					datastar.WithSelectorID(fmt.Sprintf("post-%s", res.ID)),
				)
			}
		}
	}

	return nil
}

func (e *endpoint) AddPost(c echo.Context) error {
	postParams := &PostParams{
		Title:   c.FormValue("title"),
		Content: c.FormValue("content"),
		Author:  c.FormValue("author"),
	}

	err := e.svc.Add(c.Request().Context(), postParams)
	if err != nil {
		c.Response().WriteHeader(echo.ErrInternalServerError.Code)
		fmt.Fprintf(c.Response().Writer, "Error adding posts: %v", err)
	}

	if err := e.listener.Notify(c.Request().Context()); err != nil {
		fmt.Fprintf(c.Response().Writer, "Error notify posts: %v", err)
	}

	return c.String(http.StatusCreated, "post created")
}

func (e *endpoint) DeletePost(c echo.Context) error {
	postParams := &PostParams{
		ID: c.Param("id"),
	}

	if err := e.svc.Delete(c.Request().Context(), postParams); err != nil {
		c.Response().WriteHeader(echo.ErrInternalServerError.Code)
		fmt.Fprintf(c.Response(), "Error deleting post: %v", err)
	}

	if err := e.listener.Notify(c.Request().Context()); err != nil {
		fmt.Fprintf(c.Response().Writer, "Error notify posts: %v", err)
	}

	return c.String(http.StatusNoContent, "post deleted")
}

func (e *endpoint) UpdatePost(c echo.Context) error {
	postParams := &PostParams{
		ID:      c.Param("id"),
		Title:   c.FormValue("title"),
		Content: c.FormValue("content"),
		Author:  c.FormValue("author"),
	}

	if err := e.svc.Update(c.Request().Context(), postParams); err != nil {
		c.Response().WriteHeader(echo.ErrInternalServerError.Code)
		fmt.Fprintf(c.Response(), "Error updating posts: %v", err)
	}

	if err := e.listener.Notify(c.Request().Context()); err != nil {
		fmt.Fprintf(c.Response().Writer, "Error notify posts: %v", err)
	}

	return c.String(http.StatusNoContent, "post updated")
}
