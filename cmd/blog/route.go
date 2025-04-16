package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/reyhardy/go-blog/db/scylladb"
	"github.com/reyhardy/go-blog/internal/blog"
)

func routes(dbClient scylladb.Client) {
	e := echo.New()
	blogEP := blog.NewAPI(dbClient)

	// middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// view
	e.GET("/", blogEP.GetHomePage)
	e.GET("/form/add-post", blogEP.GetAddPostFormPage)

	// fragment
	e.GET("/f/posts", blogEP.GetAllPostsFragment)
	e.GET("/f/form/add-post", blogEP.GetAddPostFormFragment)
	http.HandleFunc("GET /edit-form/{id}", blogEP.GetEditForm)
	http.HandleFunc("GET /post/{id}", blogEP.GetPost)
	http.HandleFunc("PUT /post/{id}", blogEP.EditPost)
	// http.HandleFunc("DELETE /post/{id}", blogEP.DeletePost)

	//endpoint
	e.POST("/post", blogEP.AddPostEcho)
	e.DELETE("/post/:id", blogEP.DeletePostEcho)

	e.Static("/static", "./public")

	e.Logger.Fatal(e.Start(":3030"))
}
