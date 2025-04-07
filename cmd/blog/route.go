package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/reyhardy/go-blog/db/scylladb"
	"github.com/reyhardy/go-blog/internal/blog"
)

func routes(dbClient scylladb.Client) {
	e := echo.New()
	blogEP := blog.NewAPI(dbClient)

	// view
	e.GET("/", blogEP.GetHomePage)
	e.GET("/add-post", blogEP.GetAddPostFormPage)

	// fragment
	e.GET("/posts", blogEP.GetAllPostsFragment)
	e.GET("/f/add-post", blogEP.GetAddPostFormFragment)
	http.HandleFunc("GET /edit-form/{id}", blogEP.GetEditForm)
	http.HandleFunc("GET /post/{id}", blogEP.GetPost)
	http.HandleFunc("POST /post", blogEP.AddPost)
	http.HandleFunc("PUT /post/{id}", blogEP.EditPost)
	http.HandleFunc("DELETE /post/{id}", blogEP.DeletePost)

	e.Static("/static", "./public")

	e.Logger.Fatal(e.Start(":3030"))
}
