package blog

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/reyhardy/go-blog/pkg/pgsql"
	"github.com/starfederation/datastar-go/datastar"
	"maragu.dev/gomponents"
)

const (
	datastarSelector string = "datastar-selector"
	datastarMode     string = "datastar-mode"
)

type endpoint struct {
	svc      servicer
	listener listener
}

type API interface {
	GetHome(c echo.Context) error
	GetPostSSE(c echo.Context) error
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
			if res.IsDeleted {
				sse.PatchElements(
					"",
					datastar.WithModeRemove(),
					datastar.WithSelectorID(fmt.Sprintf("post-%s", res.ID)),
				)
			} else {
				sse.PatchElements(
					gomponents.NodeFunc(PostCard(&res).Render).String(),
					datastar.WithModeAppend(),
					datastar.WithSelectorID("posts"),
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

	return c.NoContent(http.StatusCreated)
}

func (e *endpoint) DeletePost(c echo.Context) error {
	postParams := &PostParams{
		ID: c.Param("id"),
	}

	if err := e.svc.Delete(c.Request().Context(), postParams); err != nil {
		c.Response().WriteHeader(echo.ErrInternalServerError.Code)
		fmt.Fprintf(c.Response(), "Error deleting post: %v", err)
	}

	return c.NoContent(http.StatusNoContent)
}

func (e *endpoint) UpdatePost(c echo.Context) error {
	postParams := &PostParams{
		ID:      c.Param("id"),
		Title:   c.FormValue("title"),
		Content: c.FormValue("content"),
		Author:  c.FormValue("author"),
	}

	res, err := e.svc.Update(c.Request().Context(), postParams)
	if err != nil {
		c.Response().WriteHeader(echo.ErrInternalServerError.Code)
		fmt.Fprintf(c.Response(), "Error adding posts: %v", err)
	}

	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)
	c.Response().Header().Set(datastarSelector, fmt.Sprintf("#post-%s", postParams.ID))

	return PostCard(res).Render(c.Response().Writer)
}
