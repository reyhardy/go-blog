package blog

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/reyhardy/go-blog/pkg/pgsql"
)

type listen struct {
	db pgsql.Client
}

type listener interface {
	ListenForNotification(ctx context.Context, channel string, notificationC chan<- Post)
	Listen(ctx context.Context, notification chan<- string)
	Notify(ctx context.Context) error
}

func NewListener(db pgsql.Client) listener {
	return &listen{db}
}

func (l *listen) ListenForNotification(ctx context.Context, channel string, notificationC chan<- Post) {
	conn, err := l.db.Acquire(ctx)
	if err != nil {
		fmt.Println("error acquire:", err)
	}
	defer conn.Release()

	if err := conn.Ping(ctx); err != nil {
		fmt.Println("error ping:", err)
	}

	if _, err := conn.Exec(ctx, "LISTEN db_changes;"); err != nil {
		fmt.Println("error listening:", err)
	}

	fmt.Println("listening to db_change...")

	for {
		var post Post
		n, err := conn.Conn().WaitForNotification(ctx)
		if err != nil {
			fmt.Println("error notification:", err)
			break
		}

		fmt.Printf("notification payload: %v\n", n.Payload)
		json.Unmarshal([]byte(n.Payload), &post)

		notificationC <- post
	}
}

func (l *listen) Listen(ctx context.Context, notification chan<- string) {
	conn, err := l.db.Acquire(ctx)
	if err != nil {
		fmt.Println("error acquire:", err)
	}
	defer conn.Release()

	if err := conn.Ping(ctx); err != nil {
		fmt.Println("error ping:", err)
	}

	if _, err := conn.Exec(ctx, "LISTEN db_changes;"); err != nil {
		fmt.Println("error listening:", err)
	}

	fmt.Println("listening to db_changes...")

	for {
		n, err := conn.Conn().WaitForNotification(ctx)
		if err != nil {
			fmt.Println("error notification:", err)
			break
		}

		fmt.Printf("notification payload: %v\n", n.Payload)
		notification <- n.Payload
	}
}

func (l *listen) Notify(ctx context.Context) error {
	return l.db.Exec(ctx, "NOTIFY db_changes, 'this is payload';")
}
