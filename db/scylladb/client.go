package scylladb

import (
	"context"

	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx/v3"
	"github.com/scylladb/gocqlx/v3/qb"
)

type Client struct {
	session gocqlx.Session
}

func NewClient(cfg gocql.ClusterConfig) (*Client, error) {
	session, err := gocqlx.WrapSession(gocql.NewSession(cfg))
	if err != nil {
		return nil, err
	}
	return &Client{session}, nil
}

func (c *Client) QueryRow(ctx context.Context, qb qb.Builder, data ...any) (*gocqlx.Iterx, error) {
	q := c.session.Query(qb.ToCql()).WithContext(ctx)
	iter := q.Iter()
	err := q.ExecRelease()
	return iter, err
}

func (c *Client) QueryExec(ctx context.Context, qb qb.Builder, data interface{}) error {
	return c.session.Query(qb.ToCql()).WithContext(ctx).BindStruct(data).ExecRelease()
}
