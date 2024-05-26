package main

import (
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/vaidik-bajpai/realtime-gql/graph"
	"github.com/vaidik-bajpai/realtime-gql/internal/data"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	resolver := graph.NewResolver()
	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: resolver}))

	srv.AddTransport(&transport.Websocket{})

	http.Handle("/", data.Authenticate(playground.Handler("GraphQL playground", "/query")))
	http.Handle("/query", data.Authenticate(srv))

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
