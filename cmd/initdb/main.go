package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gocql/gocql"
	"github.com/reyhardy/go-blog/internal/blog"
)

func main() {
	cluster := gocql.NewCluster("localhost")
	cluster.ProtoVersion = 4
	session, err := cluster.CreateSession()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to session: %v\n", err)
	}

	defer session.Close()

	if err = session.Query(fmt.Sprintf("CREATE KEYSPACE IF NOT EXISTS %v WITH replication = {'class': 'SimpleStrategy', 'replication_factor': 1};", blog.Keyspace)).Exec(); err != nil {
		log.Fatalln("error create keyspace", err)
	} else {
		fmt.Println("Keyspace Created")
	}

	if err = session.Query(fmt.Sprintf("CREATE TABLE IF NOT EXISTS %v.%v (id varchar PRIMARY KEY, title varchar, content varchar, author varchar, created_at timestamp, updated_at timestamp);", blog.Keyspace, blog.TablePost)).Exec(); err != nil {
		log.Fatalln("error create table", err)
	} else {
		fmt.Println("Table Created")
	}
}
