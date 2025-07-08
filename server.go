package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/reginbald/gqlgen-dataloader-subscription/graph"
	"github.com/reginbald/gqlgen-dataloader-subscription/loaders"
	"github.com/reginbald/gqlgen-dataloader-subscription/repository"
	"github.com/vektah/gqlparser/v2/ast"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	connStr := "postgres://postgres:postgres@localhost:5432/data?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	repo := &repository.Repository{DB: db}

	// Change name every second
	go func() {
		for {
			// In our example we'll send the current time every second.
			time.Sleep(1 * time.Second)

			currentTime := time.Now()
			users, err := repo.GetUsers()
			if err != nil {
				return
			}
			for _, user := range users {
				repo.UpdateUser(user.ID, fmt.Sprintf("user %s %d", user.ID, int(currentTime.Unix())))
			}
		}
	}()

	h := handler.New(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{Repo: repo}}))

	h.AddTransport(transport.Websocket{})
	h.AddTransport(transport.Options{})
	h.AddTransport(transport.GET{})
	h.AddTransport(transport.POST{})

	h.SetQueryCache(lru.New[*ast.QueryDocument](1000))

	h.Use(extension.Introspection{})
	h.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	// wrap the query handler with middleware to inject dataloader in requests.
	// pass in your dataloader dependencies, in this case the db connection.
	srv := loaders.Middleware(repo, h)

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
