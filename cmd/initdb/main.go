package main

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
)

func main() {
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, "postgres://go_blog_user:go_blog_password@localhost:5432/go_blog_db?sslmode=disable")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	defer conn.Close(ctx)

	if err := conn.Ping(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "Unable to ping database: %v\n", err)
	}

	fmt.Println("Connected")

	if _, err := conn.Exec(ctx, "LISTEN db_changes;"); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to listen for notifications: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("LISTENING...")

	for {
		notification, err := conn.WaitForNotification(ctx)
		if err != nil {
			fmt.Printf("WaitForNotification error: %v\n", err)
			break
		}
		fmt.Println("notification: \n", notification.Payload)
	}

	// comment out this part if re run this file
	// _, err = conn.Exec(context.Background(), "DROP DATABASE go-blog;")
	// if err != nil {
	// 	fmt.Fprintf(os.Stderr, "drop database failed: %v\n", err)
	// 	os.Exit(1)
	// }

	// _, err = conn.Exec(context.Background(), "CREATE DATABASE go-blog-;")
	// if err != nil {
	// 	fmt.Fprintf(os.Stderr, "create database failed: %v\n", err)
	// 	os.Exit(1)
	// }

	// sqlStmt := fmt.Sprintf(`
	// CREATE TABLE %s (
	// 	id varchar(40) PRIMARY KEY,
	// 	title varchar,
	// 	content varchar,
	// 	author varchar,
	// 	created_at TIMESTAMP,
	// 	updated_at TIMESTAMP
	// );`, blog.TablePost)

	// _, err = conn.Exec(context.Background(), sqlStmt)
	// if err != nil {
	// 	fmt.Fprintf(os.Stderr, "create table failed: %v\n", err)
	// 	os.Exit(1)
	// }

	// fmt.Println("init table success!")
}
