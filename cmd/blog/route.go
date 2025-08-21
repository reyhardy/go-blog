package main

import (
	"github.com/labstack/echo/v4"
	"github.com/reyhardy/go-blog/internal/blog"
	"github.com/reyhardy/go-blog/pkg/pgsql"
)

func routes(dbClient pgsql.Client) {
	blogEP := blog.NewAPI(dbClient)

	e := echo.New()

	e.GET("/", blogEP.GetHome)

	g := e.Group("/api")
	// g.GET("/posts/sse", blogEP.GetPostSSE)
	g.GET("/posts/sse", blogEP.GetPost)
	g.POST("/posts", blogEP.AddPost)
	g.DELETE("/post/:id", blogEP.DeletePost)
	g.PUT("/post/:id", blogEP.UpdatePost)

	e.Static("/static", "public")
	e.Logger.Fatal(e.Start(":3000"))
}
