# Go-Blog

This is my side project of developing a blog site using Go. I also use PicoCSS for semantic CSS styling, Datastar for reactivity, and ScyllaDB for database.
For HTML template, I use Gomponents, a Go library which able us to write HTML in Go style. I also use gocql and gocqlx, which is Go libraries for ScyllaDB driver.

## Initialize ScyllaDB

```
docker run --name go-blog -d scylladb/scylla -p 9042:9042
```

To check if ScyllaDB is already initialized, run `docker exec -it go-blog nodetool status.
If the status is UN, that means your ScyllaDB node is already up and in normal condition.

## Create Keyspace and Table

```
make init-table
```

## Run development server

- [Go](https://go.dev/doc/install) (current version 1.23.5)

Make sure to have Go bin folder in PATH variable.

In terminal at {home} directory, type

```bash
nano .bashrc
```

Then add at the bottom line of the file and save

```bash
export PATH=$PATH:/home/{your username}/go/bin
```

Install go libs dependancies:

- Install **air**

```
go install github.com/air-verse/air@latest
```

Then simply run `make dev`

> Server by default will run on port 3030
