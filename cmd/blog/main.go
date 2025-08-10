package main

import (
	"context"
	"log"

	"github.com/reyhardy/go-blog/pkg/pgsql"
)

func main() {
	dbClient, err := pgsql.NewClient(context.Background(), "postgres://go_blog_user:go_blog_password@localhost:5432/go_blog_db?sslmode=disable")
	if err != nil {
		log.Fatalln("error init pg client", err)
	}
	routes(*dbClient)
}
