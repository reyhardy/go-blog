package main

import (
	"log"

	"github.com/gocql/gocql"
	"github.com/reyhardy/go-blog/db/scylladb"
)

func main() {
	dbClient, err := scylladb.NewClient(*gocql.NewCluster("localhost"))
	if err != nil {
		log.Fatalln("error init client db", err)
	}
	routes(*dbClient)
}
