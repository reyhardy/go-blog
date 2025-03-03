package main

import (
	"fmt"
	"net/http"

	"github.com/reyhardy/go-blog/db/scylladb"
	"github.com/reyhardy/go-blog/internal/blog"
)

func routes(dbClient scylladb.Client) {
	blogEP := blog.NewAPI(dbClient)

	http.HandleFunc("GET /home", blogEP.GetHome)
	http.HandleFunc("GET /getpost", blogEP.GetPost)
	http.HandleFunc("POST /addpost", blogEP.AddPost)
	http.HandleFunc("DELETE /deletepost/{id}", blogEP.DeletePost)

	http.Handle("/static/", http.StripPrefix("/static", http.FileServer(http.Dir("public/"))))

	fmt.Println("Listening on localhost:3030")
	http.ListenAndServe(":3030", nil)
}
