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

	if _, err := conn.Exec(ctx, fmt.Sprintf("LISTEN %s;", channel)); err != nil {
		fmt.Println("error listening:", err)
	}

	fmt.Printf("listening to %s...\n", channel)

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
