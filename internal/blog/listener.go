package blog

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/reyhardy/go-blog/pkg/pgsql"
)

type listener struct {
	db pgsql.Client
}

// type dbListener interface {
// 	ListenForNotification(ctx context.Context, channel string)
// }

// func newListener(db pgsql.Client) dbListener {
// 	return &listener{db}
// }

func (l *listener) ListenForNotification(ctx context.Context, channel string, notificationC chan<- Post) {
	conn, err := l.db.Acquire(ctx)
	if err != nil {
		fmt.Println("error acquire:", err)
	}
	defer conn.Release()

	if err := conn.Ping(ctx); err != nil {
		fmt.Println("error ping:", err)
	}

	if _, err := conn.Exec(ctx, fmt.Sprintf("LISTEN %s;", channel)); err != nil {
		fmt.Println("error notify:", err)
	}

	fmt.Printf("listening to %s...\n", channel)

	for {
		var post Post
		notification, err := conn.Conn().WaitForNotification(ctx)
		if err != nil {
			fmt.Println("error notification:", err)
			break
		}
		json.Unmarshal([]byte(notification.Payload), &post)
		fmt.Printf("post struct: \n%+v\n", post)

		notificationC <- post
	}
}
