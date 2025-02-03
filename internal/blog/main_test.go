package blog_test

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/reyhardy/go-blog/db/scylladb"
	"github.com/reyhardy/go-blog/internal/blog"
	"github.com/scylladb/gocqlx/v3"
	"github.com/scylladb/gocqlx/v3/gocqlxtest"
)

var (
	svc          blog.Servicer
	session      gocqlx.Session
	testKeyspace string = "blog_test"
)

func TestMain(m *testing.M) {
	cluster := gocqlxtest.CreateCluster()
	dbClient, _ := scylladb.NewClient(*cluster)
	svc = blog.NewService(*dbClient)

	session, _ = gocqlx.WrapSession(cluster.CreateSession())
	defer session.Close()

	session.ExecStmt(fmt.Sprintf("DROP KEYSPACE %v IF EXISTS", testKeyspace))

	err := gocqlxtest.CreateKeyspace(cluster, testKeyspace)
	if err != nil {
		log.Fatalln("error create keyspace: ", err)
	}

	err = session.ExecStmt(fmt.Sprintf("CREATE TABLE IF NOT EXISTS %v.%v (id varchar PRIMARY KEY, title varchar, content varchar, author varchar, created_at timestamp, updated_at timestamp);", testKeyspace, blog.TablePost))
	if err != nil {
		log.Fatalln("error create table: ", err)
	}

	os.Exit(m.Run())
}
