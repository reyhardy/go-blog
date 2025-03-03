package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx/v3"
	"github.com/scylladb/gocqlx/v3/migrate"
)

func main() {
	cluster := gocql.NewCluster("localhost")
	cluster.ProtoVersion = 4
	session, err := gocqlx.WrapSession(cluster.CreateSession())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to session: %v\n", err)
	}

	defer session.Close()

	err = migrate.FromFS(context.Background(), session, os.DirFS("migration"))
	if err != nil {
		log.Fatalln("error migrating from file: ", err)
	}

	// if err = session.ExecStmt(fmt.Sprintf("CREATE KEYSPACE IF NOT EXISTS %v WITH replication = {'class': 'SimpleStrategy', 'replication_factor': 1};", blog.Keyspace)); err != nil {
	// 	log.Fatalln("error create keyspace", err)
	// } else {
	// 	fmt.Println("Keyspace Created")
	// }

	// if err = session.ExecStmt(fmt.Sprintf("CREATE TABLE IF NOT EXISTS %v.%v (id varchar PRIMARY KEY, title varchar, content varchar, author varchar, created_at timestamp, updated_at timestamp);", blog.Keyspace, blog.TablePost)); err != nil {
	// 	log.Fatalln("error create table", err)
	// } else {
	// 	fmt.Println("Table Created")
	// }
}
