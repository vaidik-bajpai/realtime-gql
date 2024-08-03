package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gorilla/websocket"
	"github.com/rs/cors"
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
	srv := handler.New(graph.NewExecutableSchema(graph.Config{Resolvers: resolver}))

	var allowedOrigins = map[string]bool{
		"http://localhost:5173": true,
		"http://localhost:8080": true,
		// Add more origins as needed
	}

	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			origin := r.Header.Get("Origin")
			log.Println(origin)
			// Check if the origin is in the allowedOrigins map
			if allowedOrigins[origin] {
				return true
			}
			log.Printf("WebSocket connection rejected: origin %s not allowed", origin)
			return false
		},
	}

	// Add WebSocket transport to the server
	srv.AddTransport(&transport.Websocket{
		Upgrader:              upgrader,
		KeepAlivePingInterval: 1000 * time.Second, // Add keep-alive pings to maintain the connection
	})
	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})
	srv.AddTransport(transport.MultipartForm{})

	// Set up CORS middleware
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	http.Handle("/", data.Authenticate(playground.Handler("GraphQL playground", "/query")))
	http.Handle("/query", c.Handler(data.Authenticate(srv)))

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
